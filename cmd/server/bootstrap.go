package main

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
)

// 应用级别的 bootstrap

func Boot(
	_ internal.GlobalBootstrap, // 引入全局 bootstrap
	log log.Log,
	appHook app.Hook,
) app.Bootstrap {
	log.Info("cmd/server bootstrap")

	appHook.BeforeStart(func(ctx context.Context) error {
		log.WithContext(ctx).Info("app before start")
		return nil
	})

	return nil
}
