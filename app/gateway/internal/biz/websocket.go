package biz

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
)

type WebsocketRepo interface {
}

type WebsocketUsecase struct {
	repo   WebsocketRepo
	logger *log.Helper
}

var gateway *Gateway

func NewWebsocketUsecase(repo WebsocketRepo, logger log.Logger) *WebsocketUsecase {
	return &WebsocketUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// InitGateway 初始化路由
func (uc *WebsocketUsecase) InitGateway(ctx context.Context) {
	gateway = InitGateway(uc.logger.Logger())

	// 定时检测链接，处理异常连接
	go gateway.CheckInactiveConnections()
}

// HandleConnection 处理连接
func (uc *WebsocketUsecase) HandleConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	gateway.HandleConnection(ctx, w, r)
}

