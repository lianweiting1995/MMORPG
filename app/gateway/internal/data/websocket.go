package data

import (
	"MMORPG/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type WebsocketRepo struct {
	data   *Data
	logger *log.Helper
}

func NewWebsocketRepo(data *Data, logger log.Logger) biz.WebsocketRepo {
	return &WebsocketRepo{
		data:   data,
		logger: log.NewHelper(logger),
	}
}
