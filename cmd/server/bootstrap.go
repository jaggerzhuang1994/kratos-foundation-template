package main

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
)

// 应用级别的 bootstrap

func Boot(
	_ internal.Bootstrap, // 引入内部 bootstrap 逻辑
	log log.Log,
	appHook app.Hook,
) app.Bootstrap {
	appHook.BeforeStart(func(ctx context.Context) error {
		log.Info("app before start")
		return nil
	})

	appHook.AfterStart(func(ctx context.Context) error {
		log.Info("app after start")
		return nil
	})

	appHook.BeforeStop(func(ctx context.Context) error {
		log.Info("app before stop")
		return nil
	})

	appHook.AfterStop(func(ctx context.Context) error {
		log.Info("app after stop")
		return nil
	})

	return nil
}
