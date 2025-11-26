package data

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/biz2"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewBiz1Repo,
	NewBiz2Repo,
	wire.Bind(new(biz2.Repo), new(*Biz2Repo)),
	wire.Bind(new(biz2.Repo2), new(*Biz2Repo)),
)
