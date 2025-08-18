package service

import (
	"MMORPG/app/user/internal/biz"
	"context"

	pb "MMORPG/app/user/api/user/v1"
)

type UserService struct {
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{
		uc: uc,
	}
}

func (s *UserService) Info(ctx context.Context, id int) (*pb.InfoReply, error) {
	info, err := s.uc.Info(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.InfoReply{
		Id:       info.ID,
		Username: info.UserName,
	}, nil
}
