package data

import (
	"MMORPG/app/gateway/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type LoginRepo struct {
	data *Data
	log  log.Helper
}

// LoginWithPhone implements biz.LoginRepo.
func (l LoginRepo) LoginWithPhone(context.Context) (*biz.Login, error) {
	panic("unimplemented")
}

func NewLoginRepo(data *Data, logger log.Logger) biz.LoginRepo {
	return LoginRepo{
		data: data,
		log:  *log.NewHelper(logger),
	}
}
