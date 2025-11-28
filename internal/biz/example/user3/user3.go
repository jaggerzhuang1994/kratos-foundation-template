package user3

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/conf"
)

// 其他包都是依赖 biz
// biz依赖其他包的都抽象为接口在对应包中实现

type GetUser interface {
	GetUser(ctx context.Context, id int64) (*example.User, error)
}

type User3Biz struct {
	getUser GetUser
	conf    *conf.Bootstrap
}

func NewUser3Biz(getUser GetUser, conf *conf.Bootstrap) *User3Biz {
	return &User3Biz{
		getUser: getUser,
		conf:    conf,
	}
}

func (biz *User3Biz) GetUser(ctx context.Context, id int64) (*example.User, error) {
	u, err := biz.getUser.GetUser(ctx, id)
	return u, err
}
