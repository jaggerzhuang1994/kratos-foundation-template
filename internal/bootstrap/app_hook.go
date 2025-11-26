package bootstrap

import (
	context "context"

	"github.com/go-kratos/kratos/v2"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/log"
)

type AppHook struct {
	*log.Log
}

func NewAppHook(log *log.Log) *AppHook {
	return &AppHook{log}
}

func (h *AppHook) OnInitContext(ctx context.Context) context.Context {
	h.Info("OnInitContext")
	return ctx
}

func (h *AppHook) OnInitOption(options []kratos.Option) []kratos.Option {
	h.Info("OnInitOption")
	return options
}

func (h *AppHook) OnBeforeStart(ctx context.Context) error {
	h.Info("OnBeforeStart")
	return nil
}

func (h *AppHook) OnAfterStart(ctx context.Context) error {
	h.Info("OnAfterStart")
	return nil
}

func (h *AppHook) OnBeforeStop(ctx context.Context) error {
	h.Info("OnBeforeStop")
	return nil
}

func (h *AppHook) OnAfterStop(ctx context.Context) error {
	h.Info("OnAfterStop")
	return nil
}
