//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/client"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/conf"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/data"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app_info"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/consul"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*app_info.AppInfo, conf.LocalFilePath) (*kratos.App, func(), error) {
	panic(wire.Build(
		consul.ProviderSet,
		component.ProviderSet,
		wire.NewSet(
			biz.ProviderSet,
			bootstrap.ProviderSet,
			client.ProviderSet,
			conf.ProviderSet,
			data.ProviderSet,
			service.ProviderSet,
		),
	))
}
