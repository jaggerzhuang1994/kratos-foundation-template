package biz

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user1"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user2"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user3"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	user1.NewUser1Biz,
	user2.NewUser2Biz,
	user3.NewUser3Biz,
)
