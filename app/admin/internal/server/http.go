package server

import (
	"context"
	adminV1 "qn-base/api/gen/go/admin/v1"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"qn-base/app/admin/internal/conf"
	"qn-base/app/admin/internal/service/systemuser"
	pkgLogger "qn-base/pkg/logger"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Bootstrap, userService *systemuser.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			newHTTPServerMiddleware(c, logger)...,
		),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != 0 {
		opts = append(opts, http.Timeout(time.Duration(c.Server.Http.Timeout)))
	}
	srv := http.NewServer(opts...)
	adminV1.RegisterUserHTTPServer(srv, userService)
	return srv
}

var options = []jwt.Option{
	jwt.WithClaims(func() jwtV5.Claims {
		return jwtV5.MapClaims{}
	}),
}

func newHTTPServerMiddleware(
	config *conf.Bootstrap,
	logger log.Logger,
) []middleware.Middleware {
	var ms []middleware.Middleware
	ms = append(ms, recovery.Recovery())
	ms = append(ms, tracing.Server())
	ms = append(ms, pkgLogger.SimpleTraceIdProvider())
	ms = append(ms, logging.Server(logger))

	ms = append(ms, selector.Server(
		// 认证
		jwt.Server(func(token *jwtV5.Token) (interface{}, error) {
			return []byte(config.Jwt.System.Secret), nil
		}, options...),
		// 处理ctx参数，将token解析出来的信息放到ctx中

		// 鉴权

	).Match(newHTTPServerWhiteListMatcher()).Build())
	ms = append(ms, validate.Validator())
	return ms
}

// NewWhiteListMatcher 创建jwt白名单
func newHTTPServerWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["demo"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
