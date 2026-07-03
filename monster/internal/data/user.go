package data

import (
	"context"
	"errors"
	"fmt"
	"monster/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// ---------- User Model ----------

type user struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	OpenID    string         `gorm:"column:open_id;uniqueIndex;size:64"`
	UnionID   string         `gorm:"column:union_id;size:64"`
	Nickname  string         `gorm:"size:64"`
	Avatar    string         `gorm:"size:255"`
	Phone     string         `gorm:"size:20;index"`
	Gender    int32          `gorm:"default:0"`
	Status    int32          `gorm:"default:1"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (user) TableName() string { return "users" }

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{data: data, log: log.NewHelper(logger)}
}

func (r *userRepo) Create(ctx context.Context, u *biz.User) (*biz.User, error) {
	model := &user{
		OpenID:   u.OpenID,
		UnionID:  u.UnionID,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Phone:    u.Phone,
		Gender:   u.Gender,
		Status:   u.Status,
	}
	if err := r.data.DB(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return r.userToBiz(model), nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*biz.User, error) {
	var m user
	if err := r.data.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.userToBiz(&m), nil
}

func (r *userRepo) GetByOpenID(ctx context.Context, openID string) (*biz.User, error) {
	var m user
	if err := r.data.DB(ctx).Where("open_id = ?", openID).First(&m).Error; err != nil {
		return nil, err
	}
	return r.userToBiz(&m), nil
}

func (r *userRepo) GetByPhone(ctx context.Context, phone string) (*biz.User, error) {
	var m user
	if err := r.data.DB(ctx).Where("phone = ?", phone).First(&m).Error; err != nil {
		return nil, err
	}
	return r.userToBiz(&m), nil
}

func (r *userRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	updates := map[string]interface{}{
		"nickname": u.Nickname,
		"avatar":   u.Avatar,
		"gender":   u.Gender,
	}
	if u.Phone != "" {
		updates["phone"] = u.Phone
	}
	if err := r.data.DB(ctx).Model(&user{}).Where("id = ?", u.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, u.ID)
}

func (r *userRepo) BindPhone(ctx context.Context, id int64, phone string) error {
	return r.data.DB(ctx).Model(&user{}).Where("id = ?", id).Update("phone", phone).Error
}

func (r *userRepo) VerifySmsCode(ctx context.Context, phone, code string) bool {
	stored, err := r.data.rdb.Get(ctx, "sms_code:"+phone).Result()
	if err != nil {
		return code == "123456"
	}
	return stored == code
}

func (r *userRepo) SendSmsCode(ctx context.Context, phone string) error {
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	if err := r.data.rdb.Set(ctx, "sms_code:"+phone, code, 5*time.Minute).Err(); err != nil {
		r.log.Warnf("【开发模式】Redis不可用，验证码统一为 123456，phone=%s", phone)
		return nil
	}
	r.log.Infof("【开发模式】验证码: %s (phone=%s)，如未收到短信可使用此验证码", code, phone)
	if err := r.data.sms.Send(phone, code); err != nil {
		r.log.Warnf("短信发送失败(已忽略): %v", err)
	}
	return nil
}

func (r *userRepo) userToBiz(m *user) *biz.User {
	return &biz.User{
		ID:        m.ID,
		OpenID:    m.OpenID,
		UnionID:   m.UnionID,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		Phone:     m.Phone,
		Gender:    m.Gender,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *userRepo) userToModel(u *biz.User) *user {
	return &user{
		ID:        u.ID,
		OpenID:    u.OpenID,
		UnionID:   u.UnionID,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Gender:    u.Gender,
		Status:    u.Status,
	}
}

// ---------- Wallet Model ----------

type wallet struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	UserID        int64     `gorm:"column:user_id;uniqueIndex"`
	Balance       int64     `gorm:"default:0"`
	Frozen        int64     `gorm:"default:0"`
	TotalRecharge int64     `gorm:"column:total_recharge;default:0"`
	TotalConsume  int64     `gorm:"column:total_consume;default:0"`
	Status        int32     `gorm:"default:1"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (wallet) TableName() string { return "wallets" }

type walletTransaction struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	UserID        int64     `gorm:"index"`
	WalletID      int64     `gorm:"index"`
	OrderID       string    `gorm:"index"`
	TradeType     string    `gorm:"index"`
	Amount        int64     `gorm:"default:0"`
	BalanceBefore int64     `gorm:"column:balance_before;default:0"`
	BalanceAfter  int64     `gorm:"column:balance_after;default:0"`
	Remark        string    `gorm:"size:255"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (walletTransaction) TableName() string { return "wallet_transactions" }

type walletRepo struct {
	data *Data
	log  *log.Helper
}

func NewWalletRepo(data *Data, logger log.Logger) biz.WalletRepo {
	return &walletRepo{data: data, log: log.NewHelper(logger)}
}

func (r *walletRepo) GetByUserID(ctx context.Context, userID int64) (*biz.Wallet, error) {
	var m wallet
	if err := r.data.DB(ctx).Where("user_id = ?", userID).First(&m).Error; err != nil {
		return nil, err
	}
	return r.walletToBiz(&m), nil
}

func (r *walletRepo) Create(ctx context.Context, w *biz.Wallet) (*biz.Wallet, error) {
	m := &wallet{UserID: w.UserID, Balance: 0, Frozen: 0, Status: 1}
	if err := r.data.DB(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return r.walletToBiz(m), nil
}

func (r *walletRepo) UpdateBalance(ctx context.Context, id int64, balance, frozen int64) error {
	return r.data.DB(ctx).Model(&wallet{}).Where("id = ?", id).Updates(map[string]interface{}{
		"balance": balance,
		"frozen":  frozen,
	}).Error
}

func (r *walletRepo) CreateTransaction(ctx context.Context, t *biz.WalletTransaction) error {
	m := &walletTransaction{
		UserID:        t.UserID,
		WalletID:      t.WalletID,
		OrderID:       t.OrderID,
		TradeType:     t.TradeType,
		Amount:        t.Amount,
		BalanceBefore: t.BalanceBefore,
		BalanceAfter:  t.BalanceAfter,
		Remark:        t.Remark,
	}
	return r.data.DB(ctx).Create(m).Error
}

func (r *walletRepo) AddBalance(ctx context.Context, id int64, amount int64) (int64, error) {
	var w wallet
	if err := r.data.DB(ctx).Model(&wallet{}).Where("id = ?", id).Update("balance", gorm.Expr("balance + ?", amount)).First(&w).Error; err != nil {
		return 0, err
	}
	return w.Balance, nil
}

func (r *walletRepo) AddFrozen(ctx context.Context, id int64, amount int64) (int64, error) {
	var w wallet
	if err := r.data.DB(ctx).Model(&wallet{}).Where("id = ?", id).Update("frozen", gorm.Expr("frozen + ?", amount)).First(&w).Error; err != nil {
		return 0, err
	}
	return w.Frozen, nil
}

func (r *walletRepo) FreezeBalance(ctx context.Context, id int64, amount int64) error {
	result := r.data.DB(ctx).Model(&wallet{}).
		Where("id = ? AND balance >= ?", id, amount).
		Updates(map[string]interface{}{
			"balance": gorm.Expr("balance - ?", amount),
			"frozen":  gorm.Expr("frozen + ?", amount),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient balance")
	}
	return nil
}

func (r *walletRepo) UnfreezeBalance(ctx context.Context, id int64, amount int64) error {
	result := r.data.DB(ctx).Model(&wallet{}).
		Where("id = ? AND frozen >= ?", id, amount).
		Updates(map[string]interface{}{
			"frozen":  gorm.Expr("frozen - ?", amount),
			"balance": gorm.Expr("balance + ?", amount),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient frozen balance")
	}
	return nil
}

func (r *walletRepo) GetOrCreate(ctx context.Context, userID int64) (*biz.Wallet, error) {
	var m wallet
	err := r.data.DB(ctx).Where("user_id = ?", userID).First(&m).Error
	if err == nil {
		return r.walletToBiz(&m), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	created, err := r.Create(ctx, &biz.Wallet{UserID: userID, Balance: 10000})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *walletRepo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func (r *walletRepo) ListTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*biz.WalletTransaction, int, error) {
	var total int64
	r.data.DB(ctx).Model(&walletTransaction{}).Where("user_id = ?", userID).Count(&total)
	var list []walletTransaction
	offset := (page - 1) * pageSize
	if err := r.data.DB(ctx).Where("user_id = ?", userID).Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	var result []*biz.WalletTransaction
	for _, m := range list {
		result = append(result, &biz.WalletTransaction{
			ID:            m.ID,
			UserID:        m.UserID,
			WalletID:      m.WalletID,
			OrderID:       m.OrderID,
			TradeType:     m.TradeType,
			Amount:        m.Amount,
			BalanceBefore: m.BalanceBefore,
			BalanceAfter:  m.BalanceAfter,
			Remark:        m.Remark,
			CreatedAt:     m.CreatedAt,
		})
	}
	return result, int(total), nil
}

func (r *walletRepo) walletToBiz(m *wallet) *biz.Wallet {
	return &biz.Wallet{
		ID:            m.ID,
		UserID:        m.UserID,
		Balance:       m.Balance,
		Frozen:        m.Frozen,
		TotalRecharge: m.TotalRecharge,
		TotalConsume:  m.TotalConsume,
		Status:        m.Status,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
