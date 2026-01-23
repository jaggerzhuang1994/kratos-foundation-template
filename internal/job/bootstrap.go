package job

import (
	"time"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/job"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
)

type Bootstrap any

func NewBootstrap(
	log log.Log,
	register job.Register,
) Bootstrap {
	register.Register("job1", NewJob(log, "job1", time.Second*2))
	register.Register("job2", NewJob(log, "job2", time.Second*2))
	register.Register("job3", NewJob(log, "job3", time.Second*2))
	register.Register("job4", NewJob(log, "job4", 0))
	register.Register("job5", NewJob(log, "job5", 0))

	return nil
}
