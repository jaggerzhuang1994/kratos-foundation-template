package main

import (
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	job2 "github.com/jaggerzhuang1994/kratos-foundation-template/internal/job"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/job"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/server/websocket"
)

// NewBootstrap
// bootstrap 可以注入需要调用的方法或者实体
// bootstrap 在 app.NewApp 初始化之前调用
func NewBootstrap(
	log *log.Log,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wss *websocket.Server,
	register *job.Register,

	exampleService *service.ExampleService,
	exampleWsHandler *service.ExampleWsHandler,
) bootstrap.Bootstrap {
	wss.Handle("/echo", exampleWsHandler, websocket.Upgrader{
		// 不校验来源,在websocket在线工具下可以调试
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	register.Register("job1", job2.NewJob(log, "job1", time.Second*2))
	register.Register("job2", job2.NewJob(log, "job2", time.Second*2))
	register.Register("job3", job2.NewJob(log, "job3", time.Second*2))
	register.Register("job4", job2.NewJob(log, "job4", 0))
	register.Register("job5", job2.NewJob(log, "job5", 0))

	// 绑定http服务
	example_pb.RegisterExampleServiceHTTPServer(httpServer, exampleService)

	// 绑定grpc服务
	example_pb.RegisterExampleServiceServer(grpcServer, exampleService)

	return nil
}
