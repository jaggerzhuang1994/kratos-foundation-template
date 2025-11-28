package user2

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example"
)

// 其他包都是依赖 biz
// biz依赖其他包的都抽象为接口在对应包中实现

type GetUser interface {
	GetUser(ctx context.Context, id int64) (*example.User, error)
}

type User2Biz struct {
	repo GetUser
}

func NewUser2Biz(repo GetUser) *User2Biz {
	return &User2Biz{
		repo,
	}
}

func (biz *User2Biz) GetUser(ctx context.Context, id int64) (*example.User, error) {
	u, err := biz.repo.GetUser(ctx, id)
	return u, err
}
