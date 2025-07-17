package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	v1 "github.com/leaf-rain/raindata/admin/api/admin/v1"
	"github.com/leaf-rain/raindata/admin/internal/biz"
	"github.com/leaf-rain/raindata/admin/internal/conf"
	"github.com/leaf-rain/raindata/admin/internal/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	authBiz *biz.AuthBiz,
	authSvc *service.AuthService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		NewMiddleware(authBiz, logger),
		// 跨域
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAuthHTTPServer(srv, authSvc)
	return srv
}

// NewMiddleware 创建中间件
func NewMiddleware(authBiz *biz.AuthBiz, logger log.Logger) http.ServerOption {
	return http.Middleware(
		recovery.Recovery(),
		tracing.Server(),
		logging.Server(logger),
		authBiz.AuthUser,
	)
}
