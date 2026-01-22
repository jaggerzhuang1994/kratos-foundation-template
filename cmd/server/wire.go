//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/conf"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app_info"
)

// wireApp init kratos application.
func wireApp(app_info.Version, conf.FileConfigSource) (*kratos.App, func(), error) {
	panic(wire.Build(
		kratos_foundation.ProviderSet,
		internal.ProviderSet,
		NewBootstrap,
	))
}
