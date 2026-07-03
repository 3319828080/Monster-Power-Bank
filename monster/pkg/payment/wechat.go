package payment

import (
	"context"
)

type WechatPayChannel struct {
	appID     string
	mchID     string
	apiKey    string
	notifyURL string
}

func NewWechatPayChannel(appID, mchID, apiKey, notifyURL string) *WechatPayChannel {
	return &WechatPayChannel{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		notifyURL: notifyURL,
	}
}

func (c *WechatPayChannel) Name() string { return "wechat" }

func (c *WechatPayChannel) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	// TODO: integrate with WeChat Pay v3 API
	return &PaymentResponse{
		TransactionNo: "",
		PayURL:        "",
		PayInfo:       map[string]string{"package": "prepay_id=..."},
		Channel:       "wechat",
	}, nil
}

func (c *WechatPayChannel) Refund(ctx context.Context, req *RefundRequest) (*RefundResponse, error) {
	// TODO: integrate with WeChat Pay v3 refund API
	return &RefundResponse{
		RefundNo:        "",
		ChannelRefundNo: "",
		Status:          "success",
	}, nil
}

func (c *WechatPayChannel) VerifyNotify(ctx context.Context, params map[string]string) (*PaymentResponse, error) {
	return &PaymentResponse{Channel: "wechat"}, nil
}

func (c *WechatPayChannel) VerifyRefundNotify(ctx context.Context, params map[string]string) (*RefundResponse, error) {
	return &RefundResponse{}, nil
}
