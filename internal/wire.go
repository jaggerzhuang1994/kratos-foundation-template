package internal

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/client"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/conf"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/data"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/server"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
)

var ProviderSet = wire.NewSet(
	Boot,
	biz.ProviderSet,
	client.ProviderSet,
	conf.ProviderSet,
	data.ProviderSet,
	server.ProviderSet,
	service.ProviderSet,
)
