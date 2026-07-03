package payment

import (
	"context"
	"errors"
)

var (
	ErrPaymentFailed   = errors.New("payment failed")
	ErrRefundFailed    = errors.New("refund failed")
	ErrChannelNotFound = errors.New("payment channel not found")
	ErrInvalidAmount   = errors.New("invalid payment amount")
)

type PaymentRequest struct {
	OrderNo     string
	Amount      int64
	Description string
	OpenID      string
	ClientIP    string
	NotifyURL   string
}

type PaymentResponse struct {
	TransactionNo string
	PayURL        string
	PayInfo       map[string]string
	Channel       string
}

type RefundRequest struct {
	OrderNo   string
	PaymentNo string
	Amount    int64
	Reason    string
	NotifyURL string
}

type RefundResponse struct {
	RefundNo        string
	ChannelRefundNo string
	Status          string
}

type Channel interface {
	Name() string
	Pay(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
	Refund(ctx context.Context, req *RefundRequest) (*RefundResponse, error)
	VerifyNotify(ctx context.Context, params map[string]string) (*PaymentResponse, error)
	VerifyRefundNotify(ctx context.Context, params map[string]string) (*RefundResponse, error)
}

type Manager struct {
	channels map[string]Channel
}

func NewManager() *Manager {
	return &Manager{channels: make(map[string]Channel)}
}

func (m *Manager) Register(channel Channel) {
	m.channels[channel.Name()] = channel
}

func (m *Manager) GetChannel(name string) (Channel, error) {
	ch, ok := m.channels[name]
	if !ok {
		return nil, ErrChannelNotFound
	}
	return ch, nil
}

func (m *Manager) Pay(ctx context.Context, channelName string, req *PaymentRequest) (*PaymentResponse, error) {
	ch, err := m.GetChannel(channelName)
	if err != nil {
		return nil, err
	}
	return ch.Pay(ctx, req)
}

func (m *Manager) Refund(ctx context.Context, channelName string, req *RefundRequest) (*RefundResponse, error) {
	ch, err := m.GetChannel(channelName)
	if err != nil {
		return nil, err
	}
	return ch.Refund(ctx, req)
}
