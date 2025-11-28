package user1

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example"
)

// 其他包都是依赖 biz
// biz依赖其他包的都抽象为接口在对应包中实现

type GetUser interface {
	GetUser(ctx context.Context, id int64) (*example.User, error)
}

type User1Biz struct {
	getUserRepo GetUser
}

func NewUser1Biz(getUserRepo GetUser) *User1Biz {
	return &User1Biz{getUserRepo}
}

func (biz *User1Biz) GetUser(ctx context.Context, id int64) (*example.User, error) {
	u, err := biz.getUserRepo.GetUser(ctx, id)
	return u, err
}
