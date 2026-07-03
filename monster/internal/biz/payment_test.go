package biz

import (
	"context"
	"errors"
	"io"
	"monster/pkg/payment"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ==================== CreatePayment 测试 ====================

func TestCreatePayment_金额无效(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{name: "wechat"})
	uc := NewPaymentUsecase(&mockPaymentRepo{}, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	_, _, _, err := uc.CreatePayment(context.Background(), 1, "deposit", "BIZ001", "wechat", 0, "test")
	if err == nil || err.Error() != "无效的支付金额" {
		t.Errorf("金额为0应返回错误，实际: %v", err)
	}
}

func TestCreatePayment_不支持的支付渠道(t *testing.T) {
	mgr := payment.NewManager()
	uc := NewPaymentUsecase(&mockPaymentRepo{}, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	_, _, _, err := uc.CreatePayment(context.Background(), 1, "deposit", "BIZ001", "unknown_channel", 100, "test")
	if err == nil || err.Error() != "不支持的支付渠道: unknown_channel" {
		t.Errorf("预期不支持渠道错误，实际: %v", err)
	}
}

func TestCreatePayment_支付渠道返回失败(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "alipay",
		payFn: func(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
			return nil, errors.New("channel timeout")
		},
	})
	uc := NewPaymentUsecase(&mockPaymentRepo{}, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	_, _, _, err := uc.CreatePayment(context.Background(), 1, "deposit", "BIZ001", "alipay", 100, "test")
	if err == nil {
		t.Error("支付渠道失败应返回错误")
	}
}

func TestCreatePayment_成功(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "wechat",
		payFn: func(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
			return &payment.PaymentResponse{
				TransactionNo: "TXN001",
				PayURL:        "https://pay.weixin.qq.com/xxx",
				PayInfo:       map[string]string{"prepay_id": "wx12345"},
			}, nil
		},
	})

	var createdPayment *Payment
	paymentRepo := &mockPaymentRepo{
		createFn: func(ctx context.Context, p *Payment) (*Payment, error) {
			createdPayment = p
			p.ID = 1
			return p, nil
		},
	}

	uc := NewPaymentUsecase(paymentRepo, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	p, payInfo, payURL, err := uc.CreatePayment(context.Background(), 1001, "deposit", "BIZ001", "wechat", 9900, "支付押金")
	if err != nil {
		t.Fatalf("创建支付失败: %v", err)
	}
	if p == nil || p.PaymentNo == "" {
		t.Error("支付单号不应为空")
	}
	if createdPayment.BizType != "deposit" || createdPayment.BizNo != "BIZ001" {
		t.Errorf("业务信息不对: %+v", createdPayment)
	}
	if payURL != "https://pay.weixin.qq.com/xxx" {
		t.Errorf("支付URL不对: %s", payURL)
	}
	if payInfo == nil || payInfo["prepay_id"] != "wx12345" {
		t.Errorf("支付参数不对: %v", payInfo)
	}
}

// ==================== HandleNotify 回调测试 ====================

func TestHandleNotify_不支持的支付渠道(t *testing.T) {
	mgr := payment.NewManager()
	uc := NewPaymentUsecase(&mockPaymentRepo{}, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	err := uc.HandleNotify(context.Background(), "unknown_channel", map[string]string{})
	if err == nil {
		t.Error("不支持渠道应返回错误")
	}
}

func TestHandleNotify_验签失败(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "wechat",
		verifyNotifyFn: func(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error) {
			return nil, errors.New("signature mismatch")
		},
	})
	uc := NewPaymentUsecase(&mockPaymentRepo{}, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	err := uc.HandleNotify(context.Background(), "wechat", map[string]string{"order_no": "PAY001"})
	if err == nil || err.Error() != "回调验签失败: signature mismatch" {
		t.Errorf("验签失败应返回错误，实际: %v", err)
	}
}

func TestHandleNotify_支付单不存在(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "wechat",
		verifyNotifyFn: func(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error) {
			return &payment.PaymentResponse{
				PayInfo: map[string]string{"order_no": "PAY_NOT_EXIST"},
			}, nil
		},
	})
	paymentRepo := &mockPaymentRepo{
		getByPaymentNoFn: func(ctx context.Context, paymentNo string) (*Payment, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewPaymentUsecase(paymentRepo, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	err := uc.HandleNotify(context.Background(), "wechat", map[string]string{})
	if err == nil || err.Error() != "支付单不存在: PAY_NOT_EXIST" {
		t.Errorf("预期'支付单不存在'，实际: %v", err)
	}
}

func TestHandleNotify_重复回调跳过(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "wechat",
		verifyNotifyFn: func(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error) {
			return &payment.PaymentResponse{
				PayInfo: map[string]string{"order_no": "PAY001"},
			}, nil
		},
	})
	paymentRepo := &mockPaymentRepo{
		getByPaymentNoFn: func(ctx context.Context, paymentNo string) (*Payment, error) {
			return &Payment{ID: 1, PaymentNo: "PAY001", Status: "success"}, nil
		},
	}
	uc := NewPaymentUsecase(paymentRepo, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))

	// 已成功的支付单重复回调不应报错
	err := uc.HandleNotify(context.Background(), "wechat", map[string]string{})
	if err != nil {
		t.Fatalf("重复回调不应报错: %v", err)
	}
}

func TestHandleNotify_支付成功_deposit充值(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "wechat",
		verifyNotifyFn: func(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error) {
			return &payment.PaymentResponse{
				TransactionNo: "TXN001",
				PayInfo:       map[string]string{"order_no": "PAY001", "trade_status": "success"},
			}, nil
		},
	})

	paymentRepo := &mockPaymentRepo{
		getByPaymentNoFn: func(ctx context.Context, paymentNo string) (*Payment, error) {
			return &Payment{ID: 1, PaymentNo: "PAY001", Status: "pending", BizType: "deposit", BizNo: "BIZ001", UserID: 1001, Amount: 10000}, nil
		},
	}
	walletAdded := int64(0)
	wallet := &mockWalletRepo{
		getOrCreateFn: func(ctx context.Context, userID int64) (*Wallet, error) {
			return &Wallet{ID: 1, UserID: userID, Balance: 50000, Frozen: 0}, nil
		},
		addBalanceFn: func(ctx context.Context, id int64, amount int64) (int64, error) {
			walletAdded = amount
			return 50000 + amount, nil
		},
	}

	uc := NewPaymentUsecase(paymentRepo, mgr, &mockOrderRepo{}, wallet, log.NewStdLogger(io.Discard))
	err := uc.HandleNotify(context.Background(), "wechat", map[string]string{})
	if err != nil {
		t.Fatalf("回调处理失败: %v", err)
	}
	if walletAdded != 10000 {
		t.Errorf("充值金额应为10000，实际 %d", walletAdded)
	}
}

func TestHandleNotify_支付成功_rental更新订单(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "wechat",
		verifyNotifyFn: func(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error) {
			return &payment.PaymentResponse{
				TransactionNo: "TXN002",
				PayInfo:       map[string]string{"order_no": "PAY002", "trade_status": "success"},
			}, nil
		},
	})

	paymentRepo := &mockPaymentRepo{
		getByPaymentNoFn: func(ctx context.Context, paymentNo string) (*Payment, error) {
			return &Payment{ID: 2, PaymentNo: "PAY002", Status: "pending", BizType: "rental", BizNo: "ORD001", UserID: 1001, Amount: 500}, nil
		},
	}

	var updatedStatus string
	orderRepo := &mockOrderRepo{
		getByOrderNoFn: func(ctx context.Context, orderNo string) (*Order, error) {
			return &Order{ID: 1, OrderNo: "ORD001"}, nil
		},
		updateOrderPartialFn: func(ctx context.Context, id int64, updates map[string]any) error {
			updatedStatus = updates["status"].(string)
			return nil
		},
	}

	uc := NewPaymentUsecase(paymentRepo, mgr, orderRepo, &mockWalletRepo{}, log.NewStdLogger(io.Discard))
	err := uc.HandleNotify(context.Background(), "wechat", map[string]string{})
	if err != nil {
		t.Fatalf("回调处理失败: %v", err)
	}
	if updatedStatus != "completed" {
		t.Errorf("订单状态应更新为completed，实际 %s", updatedStatus)
	}
}

func TestHandleNotify_支付失败(t *testing.T) {
	mgr := payment.NewManager()
	mgr.Register(&mockPaymentChannel{
		name: "wechat",
		verifyNotifyFn: func(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error) {
			return &payment.PaymentResponse{
				TransactionNo: "TXN003",
				PayInfo:       map[string]string{"order_no": "PAY003", "trade_status": "failed"},
			}, nil
		},
	})

	var updatedStatus string
	paymentRepo := &mockPaymentRepo{
		getByPaymentNoFn: func(ctx context.Context, paymentNo string) (*Payment, error) {
			return &Payment{ID: 3, PaymentNo: "PAY003", Status: "pending"}, nil
		},
		updateStatusFn: func(ctx context.Context, id int64, status, transactionNo string, paidAt *time.Time) error {
			updatedStatus = status
			return nil
		},
	}

	uc := NewPaymentUsecase(paymentRepo, mgr, &mockOrderRepo{}, &mockWalletRepo{}, log.NewStdLogger(io.Discard))
	err := uc.HandleNotify(context.Background(), "wechat", map[string]string{})
	if err != nil {
		t.Fatalf("回调处理失败: %v", err)
	}
	if updatedStatus != "failed" {
		t.Errorf("支付失败状态应更新为failed，实际 %s", updatedStatus)
	}
}
