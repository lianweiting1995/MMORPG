package server

import (
	"MMORPG/app/chat/internal/conf"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
)

func NewConsul(c *conf.Registry, logger log.Logger) *consul.Registry {
	config := api.DefaultConfig()
	config.Address = c.Consul.Addrs[0]
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return consul.New(client)
}
