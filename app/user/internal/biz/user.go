package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID       int64
	UserName string
}

type UserRepo interface {
	Info(ctx context.Context, id int) (*User, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *UserUsecase) Info(ctx context.Context, id int) (*User, error) {
	return uc.repo.Info(ctx, id)
}
