package server

import (
	v1 "github.com/asynccnu/classService/api/classService/v1"
	"github.com/asynccnu/classService/internal/conf"
	"github.com/asynccnu/classService/internal/metrics"
	"github.com/asynccnu/classService/internal/pkg/encoder"
	"github.com/asynccnu/classService/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	srv.Handle("/metrics", promhttp.Handler())
	v1.RegisterClassServiceHTTPServer(srv, greeter)
	return srv
}
