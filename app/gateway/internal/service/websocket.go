package service

import (
	v1 "MMORPG/api/websocket/v1"
	"MMORPG/app/gateway/internal/biz"
	"context"
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
