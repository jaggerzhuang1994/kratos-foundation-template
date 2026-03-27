package internal

import (
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/job/globaljob"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/job"
)

// 全局级别的引导，所有 cmd 都会执行这个 boot
// 可以在这边注入一些全局都需要注册的东西

type GlobalBootstrap any

func Boot(
	job job.Register,
	globalJob *globaljob.Job,
) GlobalBootstrap {
	job.Register("global", globalJob)
	return nil
}
