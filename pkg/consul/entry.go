package consul

import (
	"MMORPG/internal/conf"
	"context"
	"fmt"
	"sync"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	g "google.golang.org/grpc"
)

type ConsulInterface interface {
	GetAddrs() []string
}

var (
	conn *g.ClientConn
	mu   sync.Mutex
)

func Client(cs *conf.Consul, sr *conf.Service) (*g.ClientConn, error) {
	if conn == nil {
		mu.Lock()
		defer mu.Unlock()

		if conn == nil {
			cfg := api.DefaultConfig()
			cfg.Address = cs.Addrs[0]
			client, err := api.NewClient(cfg)
			if err != nil {
				return nil, err
			}
			consule_registry := consul.New(client)
			filter := filter.Version(sr.Version)
			selector.SetGlobalSelector(wrr.NewBuilder())
			endpoint := fmt.Sprintf("discovery:///%s", sr.Name)
			c, err := grpc.DialInsecure(
				context.Background(),
				grpc.WithEndpoint(endpoint),
				grpc.WithDiscovery(consule_registry),
				grpc.WithNodeFilter(filter),
				grpc.WithMiddleware(
					circuitbreaker.Client(),
				),
			)
			if err != nil {
				return nil, err
			}
			conn = c
		}
	}

	return conn, nil
}
