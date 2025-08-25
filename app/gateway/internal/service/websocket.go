package service

import (
	v1 "MMORPG/api/websocket/v1"
	"MMORPG/app/gateway/internal/biz"
	"context"
	"net/http"
)

type WebsocketService struct {
	v1.UnimplementedWebsocketServer

	uc *biz.WebsocketUsecase
}

func NewWebsocketService(uc *biz.WebsocketUsecase) *WebsocketService {
	return &WebsocketService{
		uc: uc,
	}
}

func (s *WebsocketService) Push(ctx context.Context, req *v1.MsgReq) (*v1.MsgReply, error) {
	return nil, nil
}

func (s *WebsocketService) InitGateway(ctx context.Context) {
	s.uc.InitGateway(ctx)
}

func (s *WebsocketService) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.uc.HandleConnection(ctx, w, r)
}