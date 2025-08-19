package client

import (
	"github.com/hashicorp/consul/api"
)

type UserClient struct {
	client *api.Client
}

func NewUserClient() (*UserClient, error) {
	client, err := NewConsuleRegister()
	if err != nil {
		return nil, err
	}
	return &UserClient{
		client: client,
	}, nil
}
