package middleware

import (
	"context"
	"monster/pkg/jwt"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	openIDKey contextKey = "open_id"
)

var (
	ErrMissingToken = errors.Unauthorized("UNAUTHORIZED", "missing auth token")
	ErrInvalidToken = errors.Unauthorized("UNAUTHORIZED", "invalid auth token")
)

func Auth(j *jwt.JWT, skipPaths ...string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				path := tr.Operation()
				if path == "" {
					return handler(ctx, req)
				}
				for _, skip := range skipPaths {
					if strings.HasPrefix(path, skip) {
						return handler(ctx, req)
					}
				}

				token := tr.RequestHeader().Get("Authorization")
				if token == "" {
					return nil, ErrMissingToken
				}
				if len(token) > 7 && token[:7] == "Bearer " {
					token = token[7:]
				}
				claims, err := j.ParseToken(token)
				if err != nil {
					return nil, ErrInvalidToken
				}
				newCtx := context.WithValue(ctx, userIDKey, claims.UserID)
				newCtx = context.WithValue(newCtx, openIDKey, claims.OpenID)

				if r, ok := req.(*http.Request); ok {
					req = r.WithContext(newCtx)
				}
				return handler(newCtx, req)
			}
			return handler(ctx, req)
		}
	}
}

func UserIDFromContext(ctx context.Context) int64 {
	if uid, ok := ctx.Value(userIDKey).(int64); ok {
		return uid
	}
	return 0
}

func OpenIDFromContext(ctx context.Context) string {
	if oid, ok := ctx.Value(openIDKey).(string); ok {
		return oid
	}
	return ""
}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
