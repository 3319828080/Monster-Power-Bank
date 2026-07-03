package data

import (
	"context"
	"monster/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ---------- Payment Model ----------

type paymentModel struct {
	ID                   int64      `gorm:"primaryKey;autoIncrement"`
	PaymentNo            string     `gorm:"column:payment_no;uniqueIndex;size:64"`
	BizType              string     `gorm:"column:biz_type;size:32;index"`
	BizNo                string     `gorm:"column:biz_no;size:64;index"`
	UserID               int64      `gorm:"index"`
	Channel              string     `gorm:"size:32"`
	Amount               int64      `gorm:"default:0"`
	Currency             string     `gorm:"size:8;default:CNY"`
	Status               string     `gorm:"size:32;index;default:pending"`
	ChannelTransactionNo string     `gorm:"column:channel_transaction_no;size:128"`
	PaidAt               *time.Time `gorm:"column:paid_at"`
	CreatedAt            time.Time  `gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime"`
}

func (paymentModel) TableName() string { return "payments" }

// ---------- Refund Model ----------

type refundModel struct {
	ID                   int64      `gorm:"primaryKey;autoIncrement"`
	RefundNo             string     `gorm:"column:refund_no;uniqueIndex;size:64"`
	PaymentNo            string     `gorm:"column:payment_no;index;size:64"`
	OrderNo              string     `gorm:"column:order_no;index;size:64"`
	UserID               int64      `gorm:"index"`
	Amount               int64      `gorm:"default:0"`
	Reason               string     `gorm:"size:255"`
	Status               string     `gorm:"size:32;index;default:pending"`
	ChannelTransactionNo string     `gorm:"column:channel_transaction_no;size:128"`
	RefundedAt           *time.Time `gorm:"column:refunded_at"`
	CreatedAt            time.Time  `gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime"`
}

func (refundModel) TableName() string { return "refunds" }

// ---------- Notification Model ----------

type notificationModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"column:user_id;index"`
	Title     string    `gorm:"size:128"`
	Content   string    `gorm:"size:512"`
	BizType   string    `gorm:"column:biz_type;size:32"`
	BizNo     string    `gorm:"column:biz_no;size:64"`
	IsRead    bool      `gorm:"column:is_read;default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (notificationModel) TableName() string { return "user_notifications" }

// ---------- PaymentRepo ----------

type paymentRepo struct {
	data *Data
	log  *log.Helper
}

func NewPaymentRepo(data *Data, logger log.Logger) biz.PaymentRepo {
	return &paymentRepo{data: data, log: log.NewHelper(logger)}
}

func (r *paymentRepo) Create(ctx context.Context, p *biz.Payment) (*biz.Payment, error) {
	m := &paymentModel{
		PaymentNo: p.PaymentNo,
		BizType:   p.BizType,
		BizNo:     p.BizNo,
		UserID:    p.UserID,
		Channel:   p.Channel,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    p.Status,
	}
	if err := r.data.DB(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(m), nil
}

func (r *paymentRepo) GetByPaymentNo(ctx context.Context, paymentNo string) (*biz.Payment, error) {
	var m paymentModel
	if err := r.data.DB(ctx).Where("payment_no = ?", paymentNo).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(&m), nil
}

func (r *paymentRepo) GetByOrderNo(ctx context.Context, orderNo string) (*biz.Payment, error) {
	var m paymentModel
	if err := r.data.DB(ctx).Where("biz_no = ?", orderNo).Order("id DESC").First(&m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(&m), nil
}

func (r *paymentRepo) UpdateStatus(ctx context.Context, id int64, status, transactionNo string, paidAt *time.Time) error {
	updates := map[string]interface{}{"status": status}
	if transactionNo != "" {
		updates["channel_transaction_no"] = transactionNo
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	return r.data.DB(ctx).Model(&paymentModel{}).Where("id = ?", id).Updates(updates).Error
}

func (r *paymentRepo) CreateRefund(ctx context.Context, rf *biz.Refund) (*biz.Refund, error) {
	m := &refundModel{
		RefundNo:  rf.RefundNo,
		PaymentNo: rf.PaymentNo,
		OrderNo:   rf.OrderNo,
		UserID:    rf.UserID,
		Amount:    rf.Amount,
		Reason:    rf.Reason,
		Status:    rf.Status,
	}
	if err := r.data.DB(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return r.toRefundBiz(m), nil
}

func (r *paymentRepo) GetRefundByRefundNo(ctx context.Context, refundNo string) (*biz.Refund, error) {
	var m refundModel
	if err := r.data.DB(ctx).Where("refund_no = ?", refundNo).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toRefundBiz(&m), nil
}

func (r *paymentRepo) GetRefundByOrderNo(ctx context.Context, orderNo string) (*biz.Refund, error) {
	var m refundModel
	if err := r.data.DB(ctx).Where("order_no = ?", orderNo).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toRefundBiz(&m), nil
}

func (r *paymentRepo) UpdateRefundStatus(ctx context.Context, id int64, status, transactionNo string, refundedAt *time.Time) error {
	updates := map[string]interface{}{"status": status}
	if transactionNo != "" {
		updates["channel_transaction_no"] = transactionNo
	}
	if refundedAt != nil {
		updates["refunded_at"] = refundedAt
	}
	return r.data.DB(ctx).Model(&refundModel{}).Where("id = ?", id).Updates(updates).Error
}

// ---------- Notification ----------

func (r *paymentRepo) CreateNotification(ctx context.Context, n *biz.Notification) error {
	m := &notificationModel{
		UserID:  n.UserID,
		Title:   n.Title,
		Content: n.Content,
		BizType: n.BizType,
		BizNo:   n.BizNo,
	}
	return r.data.DB(ctx).Create(m).Error
}

func (r *paymentRepo) MarkNotificationsRead(ctx context.Context, ids []int64, userID int64) error {
	return r.data.DB(ctx).Model(&notificationModel{}).
		Where("id IN ? AND user_id = ?", ids, userID).
		Update("is_read", true).Error
}

// ---------- Converters ----------

func (r *paymentRepo) toBiz(m *paymentModel) *biz.Payment {
	return &biz.Payment{
		ID:                   m.ID,
		PaymentNo:            m.PaymentNo,
		BizType:              m.BizType,
		BizNo:                m.BizNo,
		UserID:               m.UserID,
		Channel:              m.Channel,
		Amount:               m.Amount,
		Currency:             m.Currency,
		Status:               m.Status,
		ChannelTransactionNo: m.ChannelTransactionNo,
		PaidAt:               m.PaidAt,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func (r *paymentRepo) toRefundBiz(m *refundModel) *biz.Refund {
	return &biz.Refund{
		ID:                   m.ID,
		RefundNo:             m.RefundNo,
		PaymentNo:            m.PaymentNo,
		OrderNo:              m.OrderNo,
		UserID:               m.UserID,
		Amount:               m.Amount,
		Reason:               m.Reason,
		Status:               m.Status,
		ChannelTransactionNo: m.ChannelTransactionNo,
		RefundedAt:           m.RefundedAt,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}
