package data

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user1"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user2"
)

// repo 实现哪些接口在这里约束 wire.Bind 只会在 wire 时提示找不到接口或者接口不满足约束，这里直接写可以实时提示是否满足接口约束
var _ user1.GetUser = (*UserDbRepo)(nil)
var _ user2.GetUser = (*UserCacheRepo)(nil)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewUserDbRepo,
	wire.Bind(new(user1.GetUser), new(*UserDbRepo)),

	NewUserCacheRepo,
	wire.Bind(new(user2.GetUser), new(*UserCacheRepo)),
)
