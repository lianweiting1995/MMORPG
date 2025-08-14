package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Login struct {
}

type LoginRepo interface {
	LoginWithPhone(context.Context) (*Login, error)
}

type LoginUseCase struct {
	repo LoginRepo
	log  *log.Helper
}

func NewLoginUsecase(repo LoginRepo, logger log.Logger) *LoginUseCase {
	return &LoginUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
