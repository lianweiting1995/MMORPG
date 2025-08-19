package consul

import (
	"MMORPG/contract"
	"context"
	"errors"
	"fmt"
	"sync"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	g "google.golang.org/grpc"
)

type ConsulInterface interface {
	GetAddrs() []string
}

var (
	conn *g.ClientConn
	mu   sync.Mutex
)

func Client(ic any) (*g.ClientConn, error) {
	if ic == nil {
		return nil, errors.New("config is empty")
	}
	if config, ok := ic.(contract.IConfig); ok {
		if conn == nil {
			mu.Lock()
			defer mu.Unlock()

			if conn == nil {
				cfg := api.DefaultConfig()
				cfg.Address = config.GetRegistry().GetConsul().GetAddrs()[0]
				client, err := api.NewClient(cfg)
				if err != nil {
					return nil, err
				}
				consule_registry := consul.New(client)
				filter := filter.Version(config.GetServer().GetVersion())
				selector.SetGlobalSelector(wrr.NewBuilder())
				endpoint := fmt.Sprintf("discovery:///%s", config.GetServer().GetName())
				c, err := grpc.DialInsecure(
					context.Background(),
					grpc.WithEndpoint(endpoint),
					grpc.WithDiscovery(consule_registry),
					grpc.WithNodeFilter(filter),
				)
				if err != nil {
					return nil, err
				}
				conn = c
			}
		}
	}

	return conn, nil
}
