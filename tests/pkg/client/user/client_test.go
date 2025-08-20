package user_test

import (
	v1 "MMORPG/api/user/v1"
	"MMORPG/pkg/client/user"
	"context"
	"fmt"
	"testing"
)

func TestNewClient_Server(t *testing.T) {
	ctx := context.Background()
	c := user.NewClient().Server(ctx)
	reply, err := c.Info(ctx, &v1.InfoRequest{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", reply)
}
