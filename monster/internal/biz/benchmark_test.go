package biz

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ==================== calculateFee 基准测试 ====================

func BenchmarkCalculateFee_起步时长内(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateFee(30*time.Minute, 200, 60, 100, 2000)
	}
}

func BenchmarkCalculateFee_超出起步时长(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateFee(90*time.Minute, 200, 60, 100, 2000)
	}
}

func BenchmarkCalculateFee_超长时长(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateFee(24*time.Hour, 200, 60, 100, 2000)
	}
}

func BenchmarkCalculateFee_多种时长混合(b *testing.B) {
	timings := []time.Duration{
		5 * time.Minute,
		30 * time.Minute,
		61 * time.Minute,
		90 * time.Minute,
		3 * time.Hour,
		10 * time.Hour,
		24 * time.Hour,
		72 * time.Hour,
	}
	for i := 0; i < b.N; i++ {
		calculateFee(timings[i%len(timings)], 200, 60, 100, 2000)
	}
}

// ==================== generateOrderNo 基准测试 ====================

func BenchmarkGenerateOrderNo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateOrderNo(12345)
	}
}

// ==================== PreCheckBorrow 基准测试 ====================

func BenchmarkPreCheckBorrow(b *testing.B) {
	repo := &mockOrderRepo{
		getCurrentByUserIDFn: func(ctx context.Context, userID int64) (*Order, error) {
			return nil, errors.New("not found")
		},
		getSlotByStationCabinetFn: func(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error) {
			return &Slot{ID: 1, Status: "空闲"}, nil
		},
		getPowerBankFn: func(ctx context.Context, id int64) (*PowerBank, error) {
			return &PowerBank{ID: 1, Status: "空闲"}, nil
		},
	}
	wallet := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{Balance: 100000, Frozen: 0}, nil
		},
	}
	uc := NewOrderUsecase(repo, wallet, &mockLocker{}, &PricingConfig{Deposit: 9900}, log.NewStdLogger(io.Discard))

	for i := 0; i < b.N; i++ {
		uc.PreCheckBorrow(context.Background(), 1, 1, 1, 1, 1)
	}
}

// ==================== GetRealTimeFee 基准测试 (高频读接口压测) ====================

func BenchmarkGetRealTimeFee(b *testing.B) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000}
	borrowTime := time.Now().Add(-90 * time.Minute)

	repo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return &Order{
				ID: 1, UserID: 1, BorrowTime: borrowTime, Status: "租借中",
				StartFee: 200, HourlyFee: 100, DailyCap: 2000,
			}, nil
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, pricing, log.NewStdLogger(io.Discard))

	for i := 0; i < b.N; i++ {
		uc.GetRealTimeFee(context.Background(), 1, "MP20240101")
	}
}

// ==================== 并发压力测试 ====================

// TestCreateOrderConcurrency 模拟并发创建订单，验证分布式锁能否正确防重
func TestCreateOrderConcurrency(t *testing.T) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000, Deposit: 9900}

	var mu sync.Mutex
	orderCount := 0
	createdOrders := make(map[int64]bool) // 按 userID 去重

	repo := &mockOrderRepo{
		getCurrentByUserIDFn: func(ctx context.Context, userID int64) (*Order, error) {
			mu.Lock()
			defer mu.Unlock()
			if createdOrders[userID] {
				return &Order{Status: "租借中"}, nil
			}
			return nil, errors.New("not found")
		},
		getSlotByStationCabinetFn: func(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error) {
			return &Slot{ID: slotID, Status: "空闲", PowerBankID: 1}, nil
		},
		getPowerBankFn: func(ctx context.Context, id int64) (*PowerBank, error) {
			return &PowerBank{ID: 1, Status: "空闲"}, nil
		},
		getStationFn: func(ctx context.Context, id int64) (*Station, error) {
			return &Station{ID: id, Name: "测试站"}, nil
		},
		createFn: func(ctx context.Context, o *Order) (*Order, error) {
			mu.Lock()
			defer mu.Unlock()
			orderCount++
			o.ID = int64(orderCount)
			return o, nil
		},
	}
	wallet := &mockWalletRepo{
		getOrCreateFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 100000, Frozen: 0}, nil
		},
	}

	// 模拟分布式锁：按 key 互斥
	lockMu := sync.Mutex{}
	lockHeld := make(map[string]bool)

	locker := &mockLocker{
		lockFn: func(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
			lockMu.Lock()
			defer lockMu.Unlock()
			if lockHeld[key] {
				return "", false, nil // 锁已被持有
			}
			lockHeld[key] = true
			return "token-" + key, true, nil
		},
		unlockFn: func(ctx context.Context, key, token string) error {
			lockMu.Lock()
			defer lockMu.Unlock()
			delete(lockHeld, key)
			return nil
		},
	}

	uc := NewOrderUsecase(repo, wallet, locker, pricing, log.NewStdLogger(io.Discard))

	// 10个用户并发创建订单
	concurrency := 10
	var wg sync.WaitGroup
	errCh := make(chan error, concurrency)
	successCh := make(chan int64, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(uid int64) {
			defer wg.Done()
			order, err := uc.CreateOrder(context.Background(), uid, 1, 1, 1, 1)
			if err != nil {
				errCh <- err
				return
			}
			successCh <- order.UserID
			mu.Lock()
			createdOrders[uid] = true
			mu.Unlock()
		}(int64(1000 + i))
	}
	wg.Wait()
	close(errCh)
	close(successCh)

	// 统计结果
	successCount := 0
	for range successCh {
		successCount++
	}
	errCount := 0
	for range errCh {
		errCount++
	}

	t.Logf("并发%d个用户创建订单结果：成功 %d，失败 %d", concurrency, successCount, errCount)
	if successCount == 0 {
		t.Error("应有至少一个成功创建订单")
	}
}

// TestReturnPowerBankConcurrency 模拟并发归还，验证事务一致性
func TestReturnPowerBankConcurrency(t *testing.T) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000, Deposit: 9900}
	borrowTime := time.Now().Add(-2 * time.Hour)

	var updateMu sync.Mutex
	updated := false

	repo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return &Order{
				ID: 1, OrderNo: orderNo, UserID: 1, PowerBankID: 1,
				BorrowTime: borrowTime, Status: "租借中",
				StartFee: 200, HourlyFee: 100, DailyCap: 2000, Deposit: 9900,
			}, nil
		},
		findEmptySlotFn: func(ctx context.Context, cabinetID int64) (*Slot, error) {
			return &Slot{ID: 10, StationID: 5, Status: "已借出"}, nil
		},
		updateOrderPartialFn: func(ctx context.Context, id int64, updates map[string]any) error {
			updateMu.Lock()
			defer updateMu.Unlock()
			if updated {
				return errors.New("订单已更新")
			}
			updated = true
			return nil
		},
	}

	wallet := &mockWalletRepo{
		getOrCreateFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 50000, Frozen: 9900}, nil
		},
		unfreezeBalanceFn: func(ctx context.Context, id int64, amount int64) error {
			return nil
		},
	}

	uc := NewOrderUsecase(repo, wallet, &mockLocker{}, pricing, log.NewStdLogger(io.Discard))

	// 5个并发同时请求归还同一个订单
	concurrency := 5
	var wg sync.WaitGroup
	successCount := 0
	var successMu sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := uc.ReturnPowerBank(context.Background(), 1, "MP20240101", 0, 5, 0, 0)
			if err == nil {
				successMu.Lock()
				successCount++
				successMu.Unlock()
			}
		}()
	}
	wg.Wait()

	t.Logf("并发%d次归还结果：成功 %d", concurrency, successCount)
	if successCount != 1 {
		t.Errorf("并发归还应只有1次成功，实际 %d 次", successCount)
	}

}
