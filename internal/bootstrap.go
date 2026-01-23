package internal

import (
	context2 "context"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/context"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
)

// 全局级别的引导

type Bootstrap any

func Boot(
	ctxHook context.Hook,
	logHook log.Hook,
) Bootstrap {
	ctxHook.WithContext(func(ctx context2.Context) context2.Context {
		return newCtx(ctx, "hello world")
	})
	logHook.With("custom", func(ctx context2.Context) any {
		return fromCtx(ctx)
	}, "key2", "value2")
	return nil
}

type myKey struct {
}

func newCtx(ctx context2.Context, val any) context2.Context {
	return context2.WithValue(ctx, myKey{}, val)
}

func fromCtx(ctx context2.Context) any {
	return ctx.Value(myKey{})
}
