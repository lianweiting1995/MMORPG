package biz

import "github.com/go-kratos/kratos/v2/log"



type WebsocketRepo interface {

}

type WebsocketUsecase struct {
	logger *log.Helper
}

func NewWebsocketUsecase(logger log.Logger) *WebsocketUsecase {
	return &WebsocketUsecase{
		logger: log.NewHelper(logger),
	}
}
