package conf

import (
	"time"

	"monster/pkg/jwt"
)

type Sms struct {
	APIID  string `yaml:"APIID" json:"APIID"`
	APIKEY string `yaml:"APIKEY" json:"APIKEY"`
}

type Auth struct {
	JWTSecret string `yaml:"jwt_secret" json:"jwt_secret"`
	JWTExpire int    `yaml:"jwt_expire" json:"jwt_expire"`
	Issuer    string `yaml:"issuer" json:"issuer"`
}

type Alipay struct {
	AppId     string `yaml:"AppId" json:"AppId"`
	Private   string `yaml:"Private" json:"Private"`
	ReturnUrl string `yaml:"ReturnUrl" json:"ReturnUrl"`
	NotifyUrl string `yaml:"NotifyUrl" json:"NotifyUrl"`
}

type Config struct {
	Server  *Server  `yaml:"server" json:"server"`
	Data    *Data    `yaml:"data" json:"data"`
	Auth    *Auth    `yaml:"auth" json:"auth"`
	Seed    *Sms     `yaml:"Seed" json:"Seed"`
	Alipay  *Alipay  `yaml:"Alipay" json:"Alipay"`
	WeChat  *WeChat  `yaml:"wechat" json:"wechat"`
	Business *Business `yaml:"business" json:"business"`
}

func NewJWT(a *Auth) *jwt.JWT {
	return jwt.NewJWT(a.JWTSecret, time.Duration(a.JWTExpire)*time.Second, a.Issuer)
}

type WeChat struct {
	AppID        string `yaml:"app_id" json:"app_id"`
	AppSecret    string `yaml:"app_secret" json:"app_secret"`
	MchID        string `yaml:"mch_id" json:"mch_id"`
	APIKey       string `yaml:"api_key" json:"api_key"`
	NotifyURL    string `yaml:"notify_url" json:"notify_url"`
	RefundNotify string `yaml:"refund_notify" json:"refund_notify"`
}

type Business struct {
	DefaultStartFee   int64 `yaml:"default_start_fee" json:"default_start_fee"`
	DefaultStartMins  int   `yaml:"default_start_mins" json:"default_start_mins"`
	DefaultHourlyFee  int64 `yaml:"default_hourly_fee" json:"default_hourly_fee"`
	DefaultDailyCap   int64 `yaml:"default_daily_cap" json:"default_daily_cap"`
	DefaultDeposit    int64 `yaml:"default_deposit" json:"default_deposit"`
}
