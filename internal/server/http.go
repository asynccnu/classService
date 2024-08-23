package server

import (
	v1 "classService/api/classService/v1"
	"classService/internal/conf"
	"classService/internal/metrics"
	"classService/internal/pkg/encoder"
	"classService/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/validate"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.ClassServiceService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metrics.QPSMiddleware(),
			metrics.DelayMiddleware(),
		),
		http.ResponseEncoder(encoder.RespEncoder), // Notice: 将响应格式化
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
	v1.RegisterClassServiceHTTPServer(srv, greeter)
	return srv
}