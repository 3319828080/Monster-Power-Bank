package payment

import (
	"context"
	"fmt"
	"net/url"

	"github.com/smartwalle/alipay/v3"
)

type AlipayChannel struct {
	client    *alipay.Client
	appId     string
	returnUrl string
	notifyUrl string
}

func NewAlipayChannel(appId, privateKey, returnUrl, notifyUrl string) (*AlipayChannel, error) {
	client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		return nil, fmt.Errorf("create alipay client failed: %w", err)
	}

	return &AlipayChannel{
		client:    client,
		appId:     appId,
		returnUrl: returnUrl,
		notifyUrl: notifyUrl,
	}, nil
}

func (c *AlipayChannel) Name() string { return "alipay" }

func (c *AlipayChannel) Pay(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	p := alipay.TradePagePay{}
	p.NotifyURL = c.notifyUrl
	p.ReturnURL = c.returnUrl
	p.Subject = req.Description
	p.OutTradeNo = req.OrderNo
	p.TotalAmount = fmt.Sprintf("%.2f", float64(req.Amount)/100)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	urlObj, err := c.client.TradePagePay(p)
	if err != nil {
		return nil, fmt.Errorf("alipay page pay failed: %w", err)
	}

	return &PaymentResponse{
		PayURL:  urlObj.String(),
		Channel: "alipay",
	}, nil
}

func (c *AlipayChannel) Refund(ctx context.Context, req *RefundRequest) (*RefundResponse, error) {
	r := alipay.TradeRefund{
		OutTradeNo:   req.OrderNo,
		RefundAmount: fmt.Sprintf("%.2f", float64(req.Amount)/100),
		RefundReason: req.Reason,
		OutRequestNo: req.OrderNo,
	}

	result, err := c.client.TradeRefund(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("alipay refund failed: %w", err)
	}
	if result.IsFailure() {
		return nil, fmt.Errorf("alipay refund error: code=%s msg=%s subCode=%s subMsg=%s",
			result.Code, result.Msg, result.SubCode, result.SubMsg)
	}

	return &RefundResponse{
		RefundNo:        req.OrderNo,
		ChannelRefundNo: result.TradeNo,
		Status:          "success",
	}, nil
}

func (c *AlipayChannel) VerifyNotify(ctx context.Context, params map[string]string) (*PaymentResponse, error) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	if err := c.client.VerifySign(ctx, values); err != nil {
		return nil, fmt.Errorf("alipay notify sign verification failed: %w", err)
	}

	tradeNo := params["trade_no"]
	orderNo := params["out_trade_no"]
	status := params["trade_status"]

	var payStatus string
	switch status {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		payStatus = "success"
	default:
		payStatus = "pending"
	}

	return &PaymentResponse{
		TransactionNo: tradeNo,
		Channel:       "alipay",
		PayInfo:       map[string]string{"order_no": orderNo, "trade_status": payStatus},
	}, nil
}

func (c *AlipayChannel) VerifyRefundNotify(ctx context.Context, params map[string]string) (*RefundResponse, error) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	if err := c.client.VerifySign(ctx, values); err != nil {
		return nil, fmt.Errorf("alipay refund notify sign verification failed: %w", err)
	}

	return &RefundResponse{
		RefundNo:        params["out_trade_no"],
		ChannelRefundNo: params["trade_no"],
		Status:          "success",
	}, nil
}
