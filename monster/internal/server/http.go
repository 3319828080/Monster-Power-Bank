package server

import (
	paymentv1 "monster/api/payment/v1"
	orderv1 "monster/api/order/v1"
	userv1 "monster/api/user/v1"
	"monster/internal/conf"
	"monster/internal/service"
	"monster/pkg/jwt"
	"monster/pkg/middleware"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, j *jwt.JWT, user *service.UserService, order *service.OrderService, payment *service.PaymentService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			middleware.Auth(j,
				"/user.v1.User/Login",
				"/user.v1.User/LoginPhone",
				"/user.v1.User/SendSmsCode",
				"/payment.v1.Payment/PaymentCallback",
			),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout > 0 {
		opts = append(opts, http.Timeout(time.Duration(c.Http.Timeout)*time.Millisecond))
	}
	srv := http.NewServer(opts...)

	userv1.RegisterUserHTTPServer(srv, user)
	orderv1.RegisterOrderHTTPServer(srv, order)
	paymentv1.RegisterPaymentHTTPServer(srv, payment)

	return srv
}
