//go:build wireinject
// +build wireinject

package main

import (
	"monster/internal/biz"
	"monster/internal/conf"
	"monster/internal/data"
	"monster/internal/server"
	"monster/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, *conf.Sms, *conf.Business, *conf.WeChat, *conf.Alipay, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		conf.NewJWT,
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		wire.Bind(new(biz.Locker), new(*data.RedisLock)),
		pricingConfigProvider,
		newApp,
	))
}
