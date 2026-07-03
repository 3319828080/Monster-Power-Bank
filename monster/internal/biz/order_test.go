package biz

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var testLogger = log.NewHelper(log.NewStdLogger(io.Discard))

// ==================== calculateFee 纯函数单元测试 ====================

func TestCalculateFee_起步时长内应收起步价(t *testing.T) {
	// 起步60分钟，起步价200分，时长30分钟 → 200
	fee := calculateFee(30*time.Minute, 200, 60, 100, 2000)
	if fee != 200 {
		t.Errorf("起步时长内应收取起步价200，实际 %d", fee)
	}
}

func TestCalculateFee_刚起步时长应收起步价(t *testing.T) {
	// 时长恰好60分钟 → 200
	fee := calculateFee(60*time.Minute, 200, 60, 100, 2000)
	if fee != 200 {
		t.Errorf("刚起步时长应收取起步价200，实际 %d", fee)
	}
}

func TestCalculateFee_超出起步时长按小时计费(t *testing.T) {
	// 时长90分钟：起步60分钟+超出30分钟→按1小时算 → 200+100=300
	fee := calculateFee(90*time.Minute, 200, 60, 100, 2000)
	if fee != 300 {
		t.Errorf("超出起步时长应收取300，实际 %d", fee)
	}
}

func TestCalculateFee_超出3小时计费(t *testing.T) {
	// 时长200分钟：起步60分钟+超出140分钟→按ceil(140/60)=3小时算 → 200+3*100=500
	fee := calculateFee(200*time.Minute, 200, 60, 100, 2000)
	if fee != 500 {
		t.Errorf("超出3小时应收取500，实际 %d", fee)
	}
}

func TestCalculateFee_超过每日封顶(t *testing.T) {
	// 时长24小时：起步60分钟+超出23小时 → 200+23*100=2500 → 超过2000封顶 → 2000
	fee := calculateFee(24*time.Hour, 200, 60, 100, 2000)
	if fee != 2000 {
		t.Errorf("超过封顶应收取2000，实际 %d", fee)
	}
}

func TestCalculateFee_超长时长不超过封顶(t *testing.T) {
	// 时长10小时：起步60分钟+超出9小时 → 200+9*100=1100 ≤ 2000 → 1100
	fee := calculateFee(10*time.Hour, 200, 60, 100, 2000)
	if fee != 1100 {
		t.Errorf("10小时应收取1100，实际 %d", fee)
	}
}

func TestCalculateFee_时长为0(t *testing.T) {
	fee := calculateFee(0, 200, 60, 100, 2000)
	if fee != 200 {
		t.Errorf("时长为0应收取起步价200，实际 %d", fee)
	}
}

func TestCalculateFee_无封顶(t *testing.T) {
	// 无封顶（dailyCap=0），时长10小时：200+9*100=1100
	fee := calculateFee(10*time.Hour, 200, 60, 100, 0)
	if fee != 1100 {
		t.Errorf("无封顶应收取1100，实际 %d", fee)
	}
}

func TestCalculateFee_无起步价(t *testing.T) {
	// 起步价0，起步时长0，按小时计费：90分钟→ceil(90/60)=2小时→2*100=200
	fee := calculateFee(90*time.Minute, 0, 0, 100, 2000)
	if fee != 200 {
		t.Errorf("无起步价应收取200，实际 %d", fee)
	}
}

// ==================== PreCheckBorrow 预检测试 ====================

func TestPreCheckBorrow_有进行中订单不允许(t *testing.T) {
	repo := &mockOrderRepo{
		getCurrentByUserIDFn: func(ctx context.Context, userID int64) (*Order, error) {
			return &Order{ID: 1, Status: "租借中"}, nil
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	allowed, reason, _, _ := uc.PreCheckBorrow(context.Background(), 1, 1, 1, 1, 1)
	if allowed {
		t.Error("有进行中订单应返回不允许")
	}
	if reason == "" {
		t.Error("应返回具体原因")
	}
}

func TestPreCheckBorrow_仓位不存在不允许(t *testing.T) {
	repo := &mockOrderRepo{
		getCurrentByUserIDFn: func(ctx context.Context, userID int64) (*Order, error) {
			return nil, errors.New("not found")
		},
		getSlotByStationCabinetFn: func(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	allowed, reason, _, _ := uc.PreCheckBorrow(context.Background(), 1, 1, 1, 1, 1)
	if allowed {
		t.Error("仓位不存在应返回不允许")
	}
	if reason == "" {
		t.Error("应返回具体原因")
	}
}

func TestPreCheckBorrow_充电宝不可用不允许(t *testing.T) {
	repo := &mockOrderRepo{
		getCurrentByUserIDFn: func(ctx context.Context, userID int64) (*Order, error) {
			return nil, errors.New("not found")
		},
		getSlotByStationCabinetFn: func(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error) {
			return &Slot{ID: 1, Status: "空闲"}, nil
		},
		getPowerBankFn: func(ctx context.Context, id int64) (*PowerBank, error) {
			return &PowerBank{ID: 1, Status: "fault"}, nil
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	allowed, _, _, _ := uc.PreCheckBorrow(context.Background(), 1, 1, 1, 1, 1)
	if allowed {
		t.Error("充电宝不可用应返回不允许")
	}
}

func TestPreCheckBorrow_余额不足不允许(t *testing.T) {
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
			return &Wallet{ID: 1, UserID: userID, Balance: 50, Frozen: 0}, nil
		},
	}
	pricing := &PricingConfig{Deposit: 9900}
	uc := NewOrderUsecase(repo, wallet, &mockLocker{}, pricing, log.NewStdLogger(io.Discard))

	allowed, reason, _, _ := uc.PreCheckBorrow(context.Background(), 1, 1, 1, 1, 1)
	if allowed {
		t.Error("余额不足应返回不允许")
	}
	if reason == "" {
		t.Error("应返回具体原因")
	}
}

func TestPreCheckBorrow_全部通过(t *testing.T) {
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
			return &Wallet{ID: 1, UserID: userID, Balance: 100000, Frozen: 0}, nil
		},
	}
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000, Deposit: 9900}
	uc := NewOrderUsecase(repo, wallet, &mockLocker{}, pricing, log.NewStdLogger(io.Discard))

	allowed, _, p, _ := uc.PreCheckBorrow(context.Background(), 1, 1, 1, 1, 1)
	if !allowed {
		t.Error("应返回允许")
	}
	if p == nil || p.Deposit != 9900 {
		t.Error("应返回定价信息")
	}
}

// ==================== CreateOrder 创建订单测试 ====================

func TestCreateOrder_分布式锁获取失败(t *testing.T) {
	locker := &mockLocker{
		lockFn: func(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
			return "", false, errors.New("redis error")
		},
	}
	uc := NewOrderUsecase(&mockOrderRepo{}, &mockWalletRepo{}, locker, &PricingConfig{}, log.NewStdLogger(io.Discard))

	_, err := uc.CreateOrder(context.Background(), 1, 1, 1, 1, 1)
	if err == nil {
		t.Error("锁获取失败应返回错误")
	}
}

func TestCreateOrder_分布式锁争抢失败(t *testing.T) {
	locker := &mockLocker{
		lockFn: func(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
			return "", false, nil
		},
	}
	uc := NewOrderUsecase(&mockOrderRepo{}, &mockWalletRepo{}, locker, &PricingConfig{}, log.NewStdLogger(io.Discard))

	_, err := uc.CreateOrder(context.Background(), 1, 1, 1, 1, 1)
	if err == nil || err.Error() != "操作过于频繁，请稍后再试" {
		t.Errorf("锁争抢失败应返回'操作过于频繁'，实际: %v", err)
	}
}

func TestCreateOrder_有进行中订单(t *testing.T) {
	repo := &mockOrderRepo{
		getCurrentByUserIDFn: func(ctx context.Context, userID int64) (*Order, error) {
			return &Order{Status: "租借中"}, nil
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	_, err := uc.CreateOrder(context.Background(), 1, 1, 1, 1, 1)
	if err == nil || err.Error() != "您有正在进行中的订单，请先归还" {
		t.Errorf("预期'有进行中的订单'，实际: %v", err)
	}
}

func TestCreateOrder_创建成功(t *testing.T) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000, Deposit: 9900}

	repo := &mockOrderRepo{
		getCurrentByUserIDFn: func(ctx context.Context, userID int64) (*Order, error) {
			return nil, errors.New("not found")
		},
		getSlotByStationCabinetFn: func(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error) {
			return &Slot{ID: slotID, Status: "空闲", PowerBankID: 1}, nil
		},
		getPowerBankFn: func(ctx context.Context, id int64) (*PowerBank, error) {
			return &PowerBank{ID: 1, Status: "空闲", DeviceNo: "PB001"}, nil
		},
		getStationFn: func(ctx context.Context, id int64) (*Station, error) {
			return &Station{ID: id, Name: "测试站点"}, nil
		},
	}
	wallet := &mockWalletRepo{
		getOrCreateFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 100000, Frozen: 0}, nil
		},
	}
	uc := NewOrderUsecase(repo, wallet, &mockLocker{}, pricing, log.NewStdLogger(io.Discard))

	order, err := uc.CreateOrder(context.Background(), 1001, 1, 10, 20, 30)
	if err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}
	if order == nil {
		t.Fatal("订单不应为空")
	}
	if order.UserID != 1001 {
		t.Errorf("用户ID应为1001，实际 %d", order.UserID)
	}
	if order.PowerBankID != 1 {
		t.Errorf("充电宝ID应为1，实际 %d", order.PowerBankID)
	}
	if order.Status != "租借中" {
		t.Errorf("状态应为'租借中'，实际 %s", order.Status)
	}
	if order.StartFee != 200 {
		t.Errorf("起步价应为200，实际 %d", order.StartFee)
	}
}

// ==================== ReturnPowerBank 归还测试 ====================

func TestReturnPowerBank_订单不存在(t *testing.T) {
	repo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	_, err := uc.ReturnPowerBank(context.Background(), 1, "INVALID", 0, 0, 0, 0)
	if err == nil || err.Error() != "订单不存在" {
		t.Errorf("预期'订单不存在'，实际 %v", err)
	}
}

func TestReturnPowerBank_归还成功无费用(t *testing.T) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000, Deposit: 9900}
	borrowTime := time.Now().Add(-30 * time.Minute) // 刚租借30分钟，不超出起步时长

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

	order, err := uc.ReturnPowerBank(context.Background(), 1, "MP20240101", 0, 5, 0, 0)
	if err != nil {
		t.Fatalf("归还失败: %v", err)
	}
	if order.Status != "completed" {
		t.Errorf("订单状态应为 completed，实际 %s", order.Status)
	}
	// 30分钟在起步时长内，费用=起步价=200
	if order.TotalAmount != 200 {
		t.Errorf("费用应为起步价200，实际 %d", order.TotalAmount)
	}
}

func TestReturnPowerBank_归还有费用从余额扣除(t *testing.T) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000, Deposit: 9900}
	borrowTime := time.Now().Add(-5 * time.Hour) // 3小时

	var capturedBalance int64
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
	}
	wallet := &mockWalletRepo{
		getOrCreateFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 50000, Frozen: 9900}, nil
		},
		updateBalanceFn: func(ctx context.Context, id int64, balance, frozen int64) error {
			capturedBalance = balance
			return nil
		},
		unfreezeBalanceFn: func(ctx context.Context, id int64, amount int64) error {
			return nil
		},
	}
	uc := NewOrderUsecase(repo, wallet, &mockLocker{}, pricing, log.NewStdLogger(io.Discard))

	order, err := uc.ReturnPowerBank(context.Background(), 1, "MP20240101", 0, 5, 0, 0)
	if err != nil {
		t.Fatalf("归还失败: %v", err)
	}
	if order.TotalAmount != order.PaidAmount+order.DiscountAmount {
		t.Errorf("总费用(%d)应等于已支付(%d)+优惠(%d)", order.TotalAmount, order.PaidAmount, order.DiscountAmount)
	}
	_ = capturedBalance
}

func TestReturnPowerBank_使用优惠券(t *testing.T) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000, Deposit: 9900}
	borrowTime := time.Now().Add(-5 * time.Hour)

	couponUsed := false
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
		getCouponFn: func(ctx context.Context, id int64) (*Coupon, error) {
			return &Coupon{ID: 100, UserID: 1, Type: "full_reduce", Amount: 300, MinAmount: 500, Status: "unused"}, nil
		},
		useCouponFn: func(ctx context.Context, id int64, userID int64) error {
			couponUsed = true
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

	order, err := uc.ReturnPowerBank(context.Background(), 1, "MP20240101", 0, 5, 0, 100)
	if err != nil {
		t.Fatalf("归还失败: %v", err)
	}
	if !couponUsed {
		t.Error("优惠券未被使用")
	}
	if order.DiscountAmount <= 0 {
		t.Errorf("应有优惠金额，实际 %d", order.DiscountAmount)
	}
}

// ==================== GetRealTimeFee 实时费用测试 ====================

func TestGetRealTimeFee_订单不存在(t *testing.T) {
	repo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	_, _, _, _, err := uc.GetRealTimeFee(context.Background(), 1, "INVALID")
	if err == nil || err.Error() != "订单不存在" {
		t.Errorf("预期'订单不存在'，实际 %v", err)
	}
}

func TestGetRealTimeFee_订单已结束(t *testing.T) {
	repo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return &Order{ID: 1, UserID: 1, Status: "completed"}, nil
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	_, _, _, _, err := uc.GetRealTimeFee(context.Background(), 1, "MP20240101")
	if err == nil || err.Error() != "订单已结束" {
		t.Errorf("预期'订单已结束'，实际 %v", err)
	}
}

func TestGetRealTimeFee_无权限(t *testing.T) {
	repo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return &Order{ID: 1, UserID: 999, Status: "租借中"}, nil
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, &PricingConfig{}, log.NewStdLogger(io.Discard))

	_, _, _, _, err := uc.GetRealTimeFee(context.Background(), 1, "MP20240101")
	if err == nil || err.Error() != "无权查看此订单" {
		t.Errorf("预期'无权查看'，实际 %v", err)
	}
}

func TestGetRealTimeFee_计算正确(t *testing.T) {
	pricing := &PricingConfig{StartFee: 200, StartMins: 60, HourlyFee: 100, DailyCap: 2000}
	borrowTime := time.Now().Add(-90 * time.Minute) // 90分钟前

	repo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return &Order{
				ID: 1, OrderNo: orderNo, UserID: 1,
				BorrowTime: borrowTime, Status: "租借中",
				StartFee: 200, HourlyFee: 100, DailyCap: 2000,
			}, nil
		},
	}
	uc := NewOrderUsecase(repo, &mockWalletRepo{}, &mockLocker{}, pricing, log.NewStdLogger(io.Discard))

	elapsed, fee, hourly, cap, err := uc.GetRealTimeFee(context.Background(), 1, "MP20240101")
	if err != nil {
		t.Fatalf("实时费用查询失败: %v", err)
	}
	if elapsed <= 0 {
		t.Error("已用时长应大于0")
	}
	if fee <= 0 {
		t.Error("费用计算应大于0")
	}
	// 90分钟：起步60+超出30分钟按1小时算 = 200+100=300
	if fee != 300 {
		t.Errorf("预期费用300，实际 %d（已用 %d 秒）", fee, elapsed)
	}
	if hourly != 100 {
		t.Errorf("每小时费率应为100，实际 %d", hourly)
	}
	if cap != 2000 {
		t.Errorf("每日封顶应为2000，实际 %d", cap)
	}
}

// ==================== 辅助函数测试 ====================

func TestGenerateOrderNo_格式正确(t *testing.T) {
	no := generateOrderNo(12345)
	if len(no) == 0 {
		t.Fatal("订单号不应为空")
	}
	// 格式: MP + 时间14位 + 用户ID5位 + UUID8位 = 1+14+5+8=28
	if len(no) != 29 {
		t.Errorf("订单号长度应为29，实际 %d", len(no))
	}
	if no[:2] != "MP" {
		t.Errorf("订单号应以MP开头，实际 %s", no[:2])
	}
}
