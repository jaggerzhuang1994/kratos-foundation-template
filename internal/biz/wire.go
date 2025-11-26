package biz

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/biz1"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/biz2"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	biz1.NewBiz1,
	biz2.NewBiz1,
)
