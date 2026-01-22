package job

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/utils"
)

var rander *rand.Rand

func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Job struct {
	log   log.Log
	name  string
	delay time.Duration
}

func NewJob(log log.Log, name string, delay time.Duration) *Job {
	return &Job{log, name, delay}
}

func (j *Job) Run(ctx context.Context) (err error) {
	j.log.Info("run job ", j.name)
	if j.delay > 0 {
		err = utils.SleepWithContext(ctx, j.delay)
	}
	j.log.Info("finish job ", j.name)
	if rander.Intn(2) == 0 {
		err = errors.New("job failed")
	}

	if rander.Intn(20) == 0 {
		panic("job panic")
	}
	return
}
