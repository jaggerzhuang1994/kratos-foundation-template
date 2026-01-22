package main

import (
	context2 "context"
	"time"

	log2 "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/websocket"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	job2 "github.com/jaggerzhuang1994/kratos-foundation-template/internal/job"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/context"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/job"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/server"
)

// NewBootstrap
// bootstrap 可以注入需要调用的方法或者实体
// bootstrap 在 app.NewApp 初始化之前调用
func NewBootstrap(
	log log.Log,
	httpServer server.HttpServer,
	grpcServer server.GrpcServer,
	wss server.WebsocketServer,
	register job.Register,
	appHook app.Hook,
	ctxHook context.Hook,
	logHook log.Hook,
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

	appHook.BeforeStart(func(ctx context2.Context) error {
		log.Info("app before start")
		return nil
	})

	appHook.AfterStart(func(ctx context2.Context) error {
		log.Info("app after start")
		return nil
	})

	appHook.BeforeStop(func(ctx context2.Context) error {
		log.Info("app before stop")
		return nil
	})

	appHook.AfterStop(func(ctx context2.Context) error {
		log.Info("app after stop")
		return nil
	})

	ctxHook.WithContext(func(ctx context2.Context) context2.Context {
		return newCtx(ctx, "hello world")
	})

	logHook.With("custom", log2.Valuer(func(ctx context2.Context) any {
		return fromCtx(ctx)
	}), "key2", "value2")

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
