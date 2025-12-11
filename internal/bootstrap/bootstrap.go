package bootstrap

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/server/websocket"
)

// NewBootstrap
// bootstrap 可以注入需要调用的方法或者实体
// bootstrap 在 app.NewApp 初始化之前调用
func NewBootstrap(
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wss *websocket.Server,
	exampleService *service.ExampleService,
	exampleWsHandler *service.ExampleWsHandler,
) bootstrap.Bootstrap {
	wss.Handle("/echo", exampleWsHandler, websocket.Upgrader{
		// 不校验来源,在websocket在线工具下可以调试
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	// 绑定http服务
	example_pb.RegisterExampleServiceHTTPServer(httpServer, exampleService)

	// 绑定grpc服务
	example_pb.RegisterExampleServiceServer(grpcServer, exampleService)

	return nil
}
