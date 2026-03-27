package cronjob

import (
	"context"
	"time"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/job"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/utils"
	"github.com/jaggerzhuang1994/kratos-foundation/proto/kratos_foundation_pb/config_pb"
)

type Job struct {
	log log.Log
}

func NewJob(log log.Log) *Job {
	return &Job{log: log}
}

func (j *Job) Schedule() string {
	// 任务调度
	// 可选表达式可以参考 https://darjun.github.io/2020/06/25/godailylib/cron/
	return "@every 5s"
}

func (j *Job) Immediately() bool {
	// 是否启动的时候立即执行
	return true
}

func (j *Job) ConcurrentPolicy() job.ConcurrentPolicy {
	// 任务并行策略：当上一个任务还未执行完成，下一个任务就调度的时候，怎么处理
	//  ConcurrentPolicy_OVERLAP: 允许并行执行，默认
	//  ConcurrentPolicy_DELAY: 推迟执行
	//  ConcurrentPolicy_SKIP: 跳过执行
	return config_pb.ConcurrentPolicy_OVERLAP
}

func (j *Job) Run(ctx context.Context) error {
	defer func() {
		j.log.WithContext(ctx).Info("sleep 6s finished")
	}()
	j.log.WithContext(ctx).Info("run cron job, sleep 6s")
	return utils.SleepWithContext(ctx, time.Second*6) // 模拟耗时操作
}
