package user

import (
	v1 "MMORPG/app/user/api/user/v1"
	"MMORPG/app/user/helper"
	"MMORPG/pkg/consul"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
)

// client
type UserClient struct {
	client *grpc.ClientConn
}

func NewClient() *UserClient {
	cfg, err := helper.Config("../../../app/user/configs")
	if err != nil {
		log.Errorf("init user grpc client failed ! %+v", err)
		return nil
	}
	c, err := consul.Client(cfg)
	if err != nil {
		log.Errorf("init user grpc client failed ! %+v", err)
		return nil
	}

	return &UserClient{
		client: c,
	}
}

func (c *UserClient) Server(ctx context.Context) v1.UserClient {
	return v1.NewUserClient(c.client)
}
