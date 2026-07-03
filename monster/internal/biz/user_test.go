package biz

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

// ==================== LoginOrRegister 测试 ====================

func TestLoginOrRegister_用户已存在直接返回(t *testing.T) {
	repo := &mockUserRepo{
		getByOpenIDFn: func(ctx context.Context, openID string) (*User, error) {
			return &User{
				ID: 100, OpenID: openID, Nickname: "老用户", Status: 1,
			}, nil
		},
	}
	uc := NewUserUsecase(repo, log.NewStdLogger(io.Discard))

	user, err := uc.LoginOrRegister(context.Background(), "openid_001", "", "新用户", "")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if user.ID != 100 {
		t.Errorf("应返回已有用户(100)，实际 %d", user.ID)
	}
	if user.Nickname != "老用户" {
		t.Errorf("昵称不应被覆盖，应为'老用户'，实际 %s", user.Nickname)
	}
}

func TestLoginOrRegister_用户不存在创建新用户(t *testing.T) {
	repo := &mockUserRepo{
		getByOpenIDFn: func(ctx context.Context, openID string) (*User, error) {
			return nil, errors.New("not found")
		},
		createFn: func(ctx context.Context, u *User) (*User, error) {
			u.ID = 200
			return u, nil
		},
	}
	uc := NewUserUsecase(repo, log.NewStdLogger(io.Discard))

	user, err := uc.LoginOrRegister(context.Background(), "openid_new", "union_new", "新用户", "http://avatar.jpg")
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}
	if user.ID != 200 {
		t.Errorf("新用户ID应为200，实际 %d", user.ID)
	}
	if user.Nickname != "新用户" {
		t.Errorf("昵称应为'新用户'，实际 %s", user.Nickname)
	}
	if user.OpenID != "openid_new" {
		t.Errorf("OpenID不对")
	}
	if user.Status != 1 {
		t.Error("新用户状态应为1(正常)")
	}
}

// ==================== PayDeposit 缴纳押金测试 ====================

func TestPayDeposit_余额不足(t *testing.T) {
	repo := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 50, Frozen: 0}, nil
		},
	}
	uc := NewWalletUsecase(repo, log.NewStdLogger(io.Discard))

	err := uc.PayDeposit(context.Background(), 1, 9900)
	if err == nil || err.Error() != "insufficient balance" {
		t.Errorf("余额不足应返回错误，实际: %v", err)
	}
}

func TestPayDeposit_钱包不存在(t *testing.T) {
	repo := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewWalletUsecase(repo, log.NewStdLogger(io.Discard))

	err := uc.PayDeposit(context.Background(), 1, 9900)
	if err == nil {
		t.Error("钱包不存在应返回错误")
	}
}

func TestPayDeposit_成功缴纳(t *testing.T) {
	var frozen bool
	repo := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 100000, Frozen: 0}, nil
		},
		freezeBalanceFn: func(ctx context.Context, id int64, amount int64) error {
			frozen = true
			return nil
		},
	}
	uc := NewWalletUsecase(repo, log.NewStdLogger(io.Discard))

	err := uc.PayDeposit(context.Background(), 1, 9900)
	if err != nil {
		t.Fatalf("缴纳押金失败: %v", err)
	}
	if !frozen {
		t.Error("应调用冻结余额")
	}
}

// ==================== RechargeDeposit 押金充值测试 ====================

func TestRechargeDeposit_金额无效(t *testing.T) {
	repo := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 100000, Frozen: 0}, nil
		},
	}
	uc := NewWalletUsecase(repo, log.NewStdLogger(io.Discard))

	err := uc.RechargeDeposit(context.Background(), 1, 0)
	if err == nil || err.Error() != "invalid amount" {
		t.Errorf("无效金额应返回错误，实际: %v", err)
	}
}

func TestRechargeDeposit_成功充值(t *testing.T) {
	var balanceAdded int64
	var frozen bool
	repo := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 100000, Frozen: 0}, nil
		},
		addBalanceFn: func(ctx context.Context, id int64, amount int64) (int64, error) {
			balanceAdded = amount
			return 100000 + amount, nil
		},
		freezeBalanceFn: func(ctx context.Context, id int64, amount int64) error {
			frozen = true
			return nil
		},
	}
	uc := NewWalletUsecase(repo, log.NewStdLogger(io.Discard))

	err := uc.RechargeDeposit(context.Background(), 1, 5000)
	if err != nil {
		t.Fatalf("押金充值失败: %v", err)
	}
	if balanceAdded != 5000 {
		t.Errorf("充值金额应为5000，实际 %d", balanceAdded)
	}
	if !frozen {
		t.Error("充值后应冻结押金")
	}
}

// ==================== RefundDeposit 退还押金测试 ====================

func TestRefundDeposit_无押金可退(t *testing.T) {
	repo := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 10000, Frozen: 0}, nil
		},
	}
	uc := NewWalletUsecase(repo, log.NewStdLogger(io.Discard))

	err := uc.RefundDeposit(context.Background(), 1)
	if err == nil || err.Error() != "no deposit to refund" {
		t.Errorf("无押金应返回错误，实际: %v", err)
	}
}

func TestRefundDeposit_成功退还(t *testing.T) {
	var unfrozen int64
	repo := &mockWalletRepo{
		getByUserIDFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 10000, Frozen: 9900}, nil
		},
		unfreezeBalanceFn: func(ctx context.Context, id int64, amount int64) error {
			unfrozen = amount
			return nil
		},
	}
	uc := NewWalletUsecase(repo, log.NewStdLogger(io.Discard))

	err := uc.RefundDeposit(context.Background(), 1)
	if err != nil {
		t.Fatalf("退还押金失败: %v", err)
	}
	if unfrozen != 9900 {
		t.Errorf("退还押金金额应为9900，实际 %d", unfrozen)
	}
}
