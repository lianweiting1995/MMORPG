package server

import (
	v1 "MMORPG/api/helloworld/v1"
	user_v1 "MMORPG/api/user/v1"
	"MMORPG/app/user/internal/conf"
	"MMORPG/app/user/internal/service"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, user *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		opts = append(opts, http.Address(addr))
	} else if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	user_v1.RegisterUserHTTPServer(srv, user)

	return srv
}
