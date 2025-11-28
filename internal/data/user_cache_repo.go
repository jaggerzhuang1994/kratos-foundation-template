package data

import (
	"context"
	"fmt"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/redis"
)

// 实现业务需要的接口

type UserCacheRepo struct {
	redis *redis.Manager
}

func NewUserCacheRepo(redis *redis.Manager) *UserCacheRepo {
	return &UserCacheRepo{redis: redis}
}

func (r *UserCacheRepo) GetUser(ctx context.Context, id int64) (*example.User, error) {
	if id == 0 {
		return nil, example.ErrUserNotFound
	}
	return &example.User{
		ID:   id,
		Name: fmt.Sprintf("user%d", id),
	}, nil
}
