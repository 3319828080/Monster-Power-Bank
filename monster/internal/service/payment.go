package service

import (
	"context"
	paymentv1 "monster/api/payment/v1"
	"monster/internal/biz"
	"monster/pkg/middleware"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var PaymentProviderSet = wire.NewSet(NewPaymentService)

type PaymentService struct {
	uc  *biz.PaymentUsecase
	log *log.Helper
}

func NewPaymentService(uc *biz.PaymentUsecase, logger log.Logger) *PaymentService {
	return &PaymentService{uc: uc, log: log.NewHelper(logger)}
}

// CreatePayment 统一拉起支付
func (s *PaymentService) CreatePayment(ctx context.Context, req *paymentv1.CreatePaymentRequest) (*paymentv1.CreatePaymentResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	payment, payParams, payURL, err := s.uc.CreatePayment(ctx, userID, req.BizType, req.BizNo, req.Channel, req.Amount, req.Description)
	if err != nil {
		return nil, err
	}
	return &paymentv1.CreatePaymentResponse{
		PaymentNo: payment.PaymentNo,
		PayParams: payParams,
		PayUrl:    payURL,
	}, nil
}

// PaymentCallback 支付异步回调
func (s *PaymentService) PaymentCallback(ctx context.Context, req *paymentv1.PaymentCallbackRequest) (*paymentv1.PaymentCallbackResponse, error) {
	if err := s.uc.HandleNotify(ctx, req.Channel, req.Params); err != nil {
		return nil, err
	}
	return &paymentv1.PaymentCallbackResponse{Message: "success"}, nil
}

// MarkNotificationRead 消息已读
func (s *PaymentService) MarkNotificationRead(ctx context.Context, req *paymentv1.MarkNotificationReadRequest) (*paymentv1.MarkNotificationReadResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if err := s.uc.MarkNotificationsRead(ctx, userID, req.Ids); err != nil {
		return nil, err
	}
	return &paymentv1.MarkNotificationReadResponse{Message: "ok"}, nil
}
