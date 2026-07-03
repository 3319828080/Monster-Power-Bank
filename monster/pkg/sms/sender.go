package sms

// Sender 短信发送接口
type Sender interface {
	Send(phone, code string) error
}
