package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"monster/pkg/payment"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/google/uuid"
)

var PaymentProviderSet = wire.NewSet(NewPaymentUsecase)

// ---------- Models ----------

type Payment struct {
	ID                   int64
	PaymentNo            string
	BizType              string // deposit, rental, compensation
	BizNo                string
	UserID               int64
	Channel              string
	Amount               int64
	Currency             string
	Status               string // pending, success, failed
	ChannelTransactionNo string
	PaidAt               *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Refund struct {
	ID                   int64
	RefundNo             string
	PaymentNo            string
	OrderNo              string
	UserID               int64
	Amount               int64
	Reason               string
	Status               string
	ChannelTransactionNo string
	RefundedAt           *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Notification struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	BizType   string
	BizNo     string
	IsRead    bool
	CreatedAt time.Time
}

// ---------- Repo Interfaces ----------

type PaymentRepo interface {
	Create(ctx context.Context, p *Payment) (*Payment, error)
	GetByPaymentNo(ctx context.Context, paymentNo string) (*Payment, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*Payment, error)
	UpdateStatus(ctx context.Context, id int64, status, transactionNo string, paidAt *time.Time) error
	CreateRefund(ctx context.Context, r *Refund) (*Refund, error)
	GetRefundByRefundNo(ctx context.Context, refundNo string) (*Refund, error)
	GetRefundByOrderNo(ctx context.Context, orderNo string) (*Refund, error)
	UpdateRefundStatus(ctx context.Context, id int64, status, transactionNo string, refundedAt *time.Time) error

	// Notifications
	CreateNotification(ctx context.Context, n *Notification) error
	MarkNotificationsRead(ctx context.Context, ids []int64, userID int64) error
}

// ---------- PaymentUsecase ----------

type PaymentUsecase struct {
	repo     PaymentRepo
	pm       *payment.Manager
	order    OrderRepo
	wallet   WalletRepo
	log      *log.Helper
}

func NewPaymentUsecase(repo PaymentRepo, pm *payment.Manager, order OrderRepo, wallet WalletRepo, logger log.Logger) *PaymentUsecase {
	return &PaymentUsecase{repo: repo, pm: pm, order: order, wallet: wallet, log: log.NewHelper(logger)}
}

// CreatePayment 统一拉起支付
func (uc *PaymentUsecase) CreatePayment(ctx context.Context, userID int64, bizType, bizNo, channel string, amount int64, description string) (*Payment, map[string]string, string, error) {
	if amount <= 0 {
		return nil, nil, "", errors.New("无效的支付金额")
	}
	if bizNo == "" {
		return nil, nil, "", errors.New("业务单号不能为空")
	}

	paymentNo := generatePaymentNo()

	// Create payment record
	p, err := uc.repo.Create(ctx, &Payment{
		PaymentNo: paymentNo,
		BizType:   bizType,
		BizNo:     bizNo,
		UserID:    userID,
		Channel:   channel,
		Amount:    amount,
		Currency:  "CNY",
		Status:    "pending",
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("创建支付单失败: %w", err)
	}

	// Call payment channel
	ch, err := uc.pm.GetChannel(channel)
	if err != nil {
		return nil, nil, "", fmt.Errorf("不支持的支付渠道: %s", channel)
	}

	payResp, err := ch.Pay(ctx, &payment.PaymentRequest{
		OrderNo:     paymentNo,
		Amount:      amount,
		Description: description,
		NotifyURL:   "", // channels set their own
	})
	if err != nil {
		uc.repo.UpdateStatus(ctx, p.ID, "failed", "", nil)
		return nil, nil, "", fmt.Errorf("支付请求失败: %w", err)
	}

	return p, payResp.PayInfo, payResp.PayURL, nil
}

// HandleNotify 支付异步回调处理
func (uc *PaymentUsecase) HandleNotify(ctx context.Context, channel string, params map[string]string) error {
	ch, err := uc.pm.GetChannel(channel)
	if err != nil {
		return fmt.Errorf("不支持的支付渠道: %s", channel)
	}

	// Verify signature
	notifyResult, err := ch.VerifyNotify(ctx, params)
	if err != nil {
		return fmt.Errorf("回调验签失败: %w", err)
	}

	paymentNo := notifyResult.PayInfo["order_no"]
	if paymentNo == "" {
		return errors.New("回调中缺少订单号")
	}

	// Lookup payment
	payment, err := uc.repo.GetByPaymentNo(ctx, paymentNo)
	if err != nil {
		return fmt.Errorf("支付单不存在: %s", paymentNo)
	}

	// Status machine - only process pending payments
	if payment.Status != "pending" {
		uc.log.Infof("支付单 %s 已处理，当前状态: %s，跳过重复回调", paymentNo, payment.Status)
		return nil
	}

	tradeStatus := notifyResult.PayInfo["trade_status"]
	now := time.Now()

	switch tradeStatus {
	case "success":
		// Update payment status
		if err := uc.repo.UpdateStatus(ctx, payment.ID, "success", notifyResult.TransactionNo, &now); err != nil {
			return fmt.Errorf("更新支付状态失败: %w", err)
		}

		// Trigger business logic
		if err := uc.handleBizSuccess(ctx, payment.BizType, payment.BizNo, payment.UserID, payment.Amount); err != nil {
			uc.log.Warnf("业务处理失败: bizType=%s bizNo=%s err=%v", payment.BizType, payment.BizNo, err)
		}

		// Create notification
		title := "支付成功"
		content := fmt.Sprintf("您已成功支付 %.2f 元", float64(payment.Amount)/100)
		_ = uc.repo.CreateNotification(ctx, &Notification{
			UserID:  payment.UserID,
			Title:   title,
			Content: content,
			BizType: payment.BizType,
			BizNo:   payment.BizNo,
		})

	default:
		if err := uc.repo.UpdateStatus(ctx, payment.ID, "failed", notifyResult.TransactionNo, &now); err != nil {
			return fmt.Errorf("更新支付失败状态: %w", err)
		}
	}

	return nil
}

// handleBizSuccess triggers post-payment business logic.
func (uc *PaymentUsecase) handleBizSuccess(ctx context.Context, bizType, bizNo string, userID, amount int64) error {
	switch bizType {
	case "deposit":
		wallet, err := uc.wallet.GetOrCreate(ctx, userID)
		if err != nil {
			return err
		}
		_, err = uc.wallet.AddBalance(ctx, wallet.ID, amount)
		return err

	case "rental":
		order, err := uc.order.GetByOrderNo(ctx, bizNo)
		if err != nil {
			return err
		}
		if err := uc.order.UpdateOrderPartial(ctx, order.ID, map[string]interface{}{
			"paid_amount": amount,
			"status":      "completed",
		}); err != nil {
			return err
		}

	case "compensation":
		order, err := uc.order.GetByOrderNo(ctx, bizNo)
		if err != nil {
			return err
		}
		if err := uc.order.UpdateOrderPartial(ctx, order.ID, map[string]interface{}{
			"paid_amount": amount,
			"status":      "compensated",
		}); err != nil {
			return err
		}
	}

	return nil
}

// MarkNotificationsRead 消息已读
func (uc *PaymentUsecase) MarkNotificationsRead(ctx context.Context, userID int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return uc.repo.MarkNotificationsRead(ctx, ids, userID)
}

// ---------- Helpers ----------

func generatePaymentNo() string {
	now := time.Now()
	return fmt.Sprintf("PAY%s%s", now.Format("20060102150405"), uuid.New().String()[:12])
}
