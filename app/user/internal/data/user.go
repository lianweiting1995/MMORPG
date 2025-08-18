package data

import (
	"MMORPG/app/user/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo struct {
	data *Data
	log  *log.Helper
}

// Info implements biz.UserRepo.
func (u UserRepo) Info(ctx context.Context, id int) (*biz.User, error) {
	return &biz.User{
		ID:       1,
		UserName: "lian",
	}, nil
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
