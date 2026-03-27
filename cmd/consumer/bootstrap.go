package main

import (
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/job/kafka_consumer"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/job"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/server"
)

// 应用级别的 bootstrap

func Boot(
	_ internal.GlobalBootstrap, // 引入全局 bootstrap
	_ server.DisableGrpc, // 脚本可以关闭grpc服务，也可以选择在配置里关闭
	log log.Log,
	job job.Register,

	consumer *kafka_consumer.Consumer,
) app.Bootstrap {
	log.Info("cmd/job bootstrap")

	job.Register("kafka_consumer", consumer)
	return nil
}
