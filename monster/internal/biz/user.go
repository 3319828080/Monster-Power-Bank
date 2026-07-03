package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/google/wire"
)

var UserProviderSet = wire.NewSet(NewUserUsecase, NewWalletUsecase)

// ---------- User ----------

type User struct {
	ID        int64
	OpenID    string
	UnionID   string
	Nickname  string
	Avatar    string
	Phone     string
	Gender    int32
	Status    int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepo interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByOpenID(ctx context.Context, openID string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	Update(ctx context.Context, u *User) (*User, error)
	BindPhone(ctx context.Context, id int64, phone string) error
	VerifySmsCode(ctx context.Context, phone, code string) bool
	SendSmsCode(ctx context.Context, phone string) error
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) LoginOrRegister(ctx context.Context, openID, unionID, nickname, avatar string) (*User, error) {
	user, err := uc.repo.GetByOpenID(ctx, openID)
	if err != nil {
		user = &User{
			OpenID:   openID,
			UnionID:  unionID,
			Nickname: nickname,
			Avatar:   avatar,
			Status:   1,
		}
		return uc.repo.Create(ctx, user)
	}
	return user, nil
}

func (uc *UserUsecase) GetUser(ctx context.Context, id int64) (*User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UserUsecase) LoginByPhone(ctx context.Context, phone, code, nickname, avatar string) (*User, error) {
	if !uc.repo.VerifySmsCode(ctx, phone, code) {
		return nil, errors.New("invalid SMS code")
	}
	user, err := uc.repo.GetByPhone(ctx, phone)
	if err != nil {
		user = &User{
			OpenID:   "phone_" + uuid.New().String(),
			Phone:    phone,
			Nickname: nickname,
			Avatar:   avatar,
			Status:   1,
		}
		return uc.repo.Create(ctx, user)
	}
	return user, nil
}

func (uc *UserUsecase) BindPhone(ctx context.Context, id int64, phone string) error {
	return uc.repo.BindPhone(ctx, id, phone)
}

func (uc *UserUsecase) SendSmsCode(ctx context.Context, phone string) error {
	return uc.repo.SendSmsCode(ctx, phone)
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, u *User) (*User, error) {
	return uc.repo.Update(ctx, u)
}

// ---------- Wallet ----------

type Wallet struct {
	ID            int64
	UserID        int64
	Balance       int64
	Frozen        int64
	TotalRecharge int64
	TotalConsume  int64
	Status        int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type WalletTransaction struct {
	ID            int64
	UserID        int64
	WalletID      int64
	OrderID       string
	TradeType     string
	Amount        int64
	BalanceBefore int64
	BalanceAfter  int64
	Remark        string
	CreatedAt     time.Time
}

type WalletRepo interface {
	GetByUserID(ctx context.Context, userID int64) (*Wallet, error)
	Create(ctx context.Context, w *Wallet) (*Wallet, error)
	UpdateBalance(ctx context.Context, id int64, balance, frozen int64) error
	CreateTransaction(ctx context.Context, t *WalletTransaction) error
	ListTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*WalletTransaction, int, error)
	AddBalance(ctx context.Context, id int64, amount int64) (int64, error)
	AddFrozen(ctx context.Context, id int64, amount int64) (int64, error)
	FreezeBalance(ctx context.Context, id int64, amount int64) error
	UnfreezeBalance(ctx context.Context, id int64, amount int64) error
	GetOrCreate(ctx context.Context, userID int64) (*Wallet, error)
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type WalletUsecase struct {
	repo WalletRepo
	log  *log.Helper
}

func NewWalletUsecase(repo WalletRepo, logger log.Logger) *WalletUsecase {
	return &WalletUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *WalletUsecase) GetWallet(ctx context.Context, userID int64) (*Wallet, error) {
	return uc.repo.GetByUserID(ctx, userID)
}

func (uc *WalletUsecase) ListTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*WalletTransaction, int, error) {
	return uc.repo.ListTransactions(ctx, userID, page, pageSize)
}

func (uc *WalletUsecase) EnsureWallet(ctx context.Context, userID int64) (*Wallet, error) {
	return uc.repo.GetOrCreate(ctx, userID)
}

func (uc *WalletUsecase) PayDeposit(ctx context.Context, userID int64, amount int64) error {
	wallet, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if wallet.Balance < amount {
		return errors.New("insufficient balance")
	}
	if err := uc.repo.FreezeBalance(ctx, wallet.ID, amount); err != nil {
		return err
	}
	_ = uc.repo.CreateTransaction(ctx, &WalletTransaction{
		UserID:        userID,
		WalletID:      wallet.ID,
		TradeType:     "deposit_pay",
		Amount:        -amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance - amount,
		Remark:        "缴纳押金",
	})
	return nil
}

func (uc *WalletUsecase) RechargeDeposit(ctx context.Context, userID int64, amount int64) error {
	wallet, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	if _, err := uc.repo.AddBalance(ctx, wallet.ID, amount); err != nil {
		return err
	}
	if err := uc.repo.FreezeBalance(ctx, wallet.ID, amount); err != nil {
		return err
	}
	_ = uc.repo.CreateTransaction(ctx, &WalletTransaction{
		UserID:        userID,
		WalletID:      wallet.ID,
		TradeType:     "deposit_recharge",
		Amount:        amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance + amount,
		Remark:        "押金快速充值",
	})
	return nil
}

func (uc *WalletUsecase) GetDepositInfo(ctx context.Context, userID int64, depositAmount int64) (bool, int64, int64, bool, string, error) {
	wallet, err := uc.repo.GetOrCreate(ctx, userID)
	if err != nil {
		return false, 0, 0, false, "", err
	}
	hasDeposit := wallet.Frozen > 0
	depositEnough := wallet.Balance+wallet.Frozen >= depositAmount
	rules := "押金在归还充电宝且无欠费后自动退还，预计1-3个工作日到账。如产生费用将优先从余额扣除，不足部分从押金扣除。"
	return hasDeposit, depositAmount, wallet.Frozen, depositEnough, rules, nil
}

func (uc *WalletUsecase) RefundDeposit(ctx context.Context, userID int64) error {
	wallet, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if wallet.Frozen <= 0 {
		return errors.New("no deposit to refund")
	}
	if err := uc.repo.UnfreezeBalance(ctx, wallet.ID, wallet.Frozen); err != nil {
		return err
	}
	_ = uc.repo.CreateTransaction(ctx, &WalletTransaction{
		UserID:        userID,
		WalletID:      wallet.ID,
		TradeType:     "deposit_refund",
		Amount:        wallet.Frozen,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  wallet.Balance + wallet.Frozen,
		Remark:        "退还押金",
	})
	return nil
}
