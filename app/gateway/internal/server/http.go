package server

import (
	v1 "MMORPG/api/helloworld/v1"
	"MMORPG/app/gateway/internal/conf"
	"MMORPG/app/gateway/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
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

		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
		}
		defer c.Close()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Errorf("connect ws failed %s", err)
				break
			}
			log.Infof("id: %d recv: %s", mt, message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Errorf("id: %d, write ws failed %s", mt, err)
				break
			}
		}
	})

	srv := http.NewServer(opts...)
	srv.HandlePrefix("/", router)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
