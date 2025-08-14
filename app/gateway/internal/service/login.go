package service

import (
	pb "MMORPG/app/gateway/api/websocket/v1"
	"MMORPG/app/gateway/internal/biz"
	"context"
)

type LoginService struct {
	uc *biz.LoginUseCase
}

func NewLoginService(uc *biz.LoginUseCase) *LoginService {
	return &LoginService{
		uc: uc,
	}
}

func (s *LoginService) LoginWithPhone(ctx context.Context, in *pb.LoginInput) (*pb.LoginOutput, error) {
	return nil, nil
}
