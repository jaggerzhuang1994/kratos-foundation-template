package globaljob

import (
	"context"
	"time"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
)

type Job struct {
	log log.Log
}

func NewJob(log log.Log) *Job {
	return &Job{log}
}

func (j *Job) Run(ctx context.Context) error {
	timer := time.NewTicker(time.Second * 5)
	defer timer.Stop()

	log := j.log.WithContext(ctx)

	for {
		log.Info("run global job")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
