package bootstrap

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/app"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/server"
)

func NewBootstrap(
	// bootstrap 可以注入需要调用的方法或者实体
	// bootstrap 在 *kratos.App/*http.Server/*grpc.Server 初始化之前调用
	// 因此不能构造函数注入这些实例，否则会循环引用
	// 可以使用对应包中提供的 hook 方法
	log *log.Log,
	appHook *app.HookManager,
	serverHook *server.HookManager,

	// 注入自己的service绑定到 httpServer/grpcServer
	exampleService *service.ExampleService,
) bootstrap.Bootstrap {
	log.Info("boot")

	appHook.BeforeStart(func(ctx context.Context) error {
		log.Info("before start")
		return nil
	})

	// 绑定http服务
	serverHook.HookHttpServer(func(s *http.Server) {
		example_pb.RegisterExampleServiceHTTPServer(s, exampleService)
	})

	// 绑定grpc服务
	serverHook.HookGrpcServer(func(s *grpc.Server) {
		example_pb.RegisterExampleServiceServer(s, exampleService)
	})

	return nil
}
