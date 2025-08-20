package service

import (
	"MMORPG/app/user/internal/biz"
	"context"

	v1 "MMORPG/api/user/v1"
)

type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{
		uc: uc,
	}
}

func (s *UserService) Info(ctx context.Context, req *v1.InfoRequest) (*v1.InfoReply, error) {
	info, err := s.uc.Info(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &v1.InfoReply{
		Id:       info.ID,
		Username: info.UserName,
	}, nil
}
