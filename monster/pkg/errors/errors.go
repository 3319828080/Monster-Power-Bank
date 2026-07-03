package errors

import (
	"net/http"

	kratos "github.com/go-kratos/kratos/v2/errors"
)

const (
	// Common
	Success            = 0
	UnknownError       = 10001
	BadRequest         = 10002
	Unauthorized       = 10003
	Forbidden          = 10004
	NotFound           = 10005
	MethodNotAllowed   = 10006
	Timeout            = 10007
	TooManyRequests    = 10008
	InternalError      = 10009
	ServiceUnavailable = 10010

	// User
	UserNotFound       = 20001
	UserAlreadyExists  = 20002
	UserDisabled       = 20003
	InvalidCredential  = 20004
	PhoneAlreadyBound  = 20005

	// Wallet
	WalletNotFound      = 30001
	InsufficientBalance = 30002
	WalletFrozen        = 30003
	InvalidAmount       = 30004

	// Order
	OrderNotFound      = 40001
	OrderCannotReturn  = 40002
	OrderExpired       = 40003
	OrderAlreadyPaid   = 40004
	InvalidOrderStatus = 40005

	// Payment
	PaymentFailed       = 50001
	PaymentNotFound     = 50002
	RefundFailed        = 50003
	ChannelNotSupported = 50004

	// Power Bank
	PowerBankNotFound = 60001
	PowerBankBusy     = 60002
	SlotOccupied      = 60003
	CabinetOffline    = 60004
)

var codeToMsg = map[int]string{
	Success:            "success",
	UnknownError:       "unknown error",
	BadRequest:         "bad request",
	Unauthorized:       "unauthorized",
	Forbidden:          "forbidden",
	NotFound:           "not found",
	MethodNotAllowed:   "method not allowed",
	Timeout:            "request timeout",
	TooManyRequests:    "too many requests",
	InternalError:      "internal server error",
	ServiceUnavailable: "service unavailable",

	UserNotFound:       "user not found",
	UserAlreadyExists:  "user already exists",
	UserDisabled:       "user is disabled",
	InvalidCredential:  "invalid credential",
	PhoneAlreadyBound:  "phone already bound",

	WalletNotFound:      "wallet not found",
	InsufficientBalance: "insufficient balance",
	WalletFrozen:        "wallet is frozen",
	InvalidAmount:       "invalid amount",

	OrderNotFound:      "order not found",
	OrderCannotReturn:  "order cannot be returned",
	OrderExpired:       "order expired",
	OrderAlreadyPaid:   "order already paid",
	InvalidOrderStatus: "invalid order status",

	PaymentFailed:       "payment failed",
	PaymentNotFound:     "payment not found",
	RefundFailed:        "refund failed",
	ChannelNotSupported: "payment channel not supported",

	PowerBankNotFound: "power bank not found",
	PowerBankBusy:     "power bank is busy",
	SlotOccupied:      "slot is occupied",
	CabinetOffline:    "cabinet is offline",
}

func Msg(code int) string {
	if msg, ok := codeToMsg[code]; ok {
		return msg
	}
	return "unknown error"
}

func New(code int, msg ...string) *kratos.Error {
	m := Msg(code)
	if len(msg) > 0 {
		m = msg[0]
	}
	return kratos.New(http.StatusOK, Msg(code), m)
}

func Newf(code int, format string, args ...interface{}) *kratos.Error {
	return kratos.New(http.StatusOK, Msg(code), format)
}

func Is(err error, code int) bool {
	if e := kratos.FromError(err); e != nil {
		return e.Reason == Msg(code)
	}
	return false
}

// HTTP status code mapping
func HTTPCode(code int) int {
	switch {
	case code == Success:
		return http.StatusOK
	case code >= 60000:
		return http.StatusServiceUnavailable
	case code >= 50000:
		return http.StatusInternalServerError
	case code >= 40000:
		return http.StatusBadRequest
	case code >= 30000:
		return http.StatusPaymentRequired
	case code >= 20000:
		return http.StatusUnauthorized
	case code >= 10000:
		if code == Unauthorized {
			return http.StatusUnauthorized
		}
		if code == Forbidden {
			return http.StatusForbidden
		}
		if code == NotFound {
			return http.StatusNotFound
		}
		if code == TooManyRequests {
			return http.StatusTooManyRequests
		}
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
