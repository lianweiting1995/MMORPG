package user

import (
	v1 "MMORPG/api/user/v1"
	"MMORPG/helper"
	"MMORPG/pkg/consul"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
)

// client
type UserClient struct {
	conn *grpc.ClientConn
}

func NewClient() *UserClient {
	cfg_path := fmt.Sprintf("%s/configs", helper.RootPath())
	cfg, err := helper.Config(cfg_path)
	if err != nil {
		log.Errorf("init user grpc client failed ! %+v", err)

		return nil
	}
	conn, err := consul.Client(cfg.Register.Consul, cfg.Services.User)
	if err != nil {
		log.Errorf("init user grpc client failed ! %+v", err)

		return nil
	}

	return &UserClient{
		conn: conn,
	}
}

func (c *UserClient) Server(ctx context.Context) v1.UserClient {
	return v1.NewUserClient(c.conn)
}
