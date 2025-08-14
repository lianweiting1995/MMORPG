package server

import (
	v1 "MMORPG/api/helloworld/v1"
	"MMORPG/app/gateway/internal/conf"
	"MMORPG/app/gateway/internal/service"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

var gateway *Gateway

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, login *service.LoginService, logger log.Logger) *http.Server {
	gateway = InitGateway(logger, login)
	// 定期检测链接,处理异常连接
	go gateway.checkInactiveConnections()
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
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

	router := mux.NewRouter()
	router.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		gateway.HandleConnection(ctx, w, r)
	})

	srv := http.NewServer(opts...)
	srv.HandlePrefix("/", router)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
