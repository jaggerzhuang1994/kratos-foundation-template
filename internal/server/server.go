package server

import (
	"net/http"

	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/server"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/websocket"
)

// Setup server 初始化
// 可以在这里注入 server.Middlewares server.HttpServerOptions server.GrpcServerOptions 去修改服务参数
// 但是不可以注入 http/grpc/websocket 实例
func Setup(
	middlewares server.Middlewares,
	httpOpts server.HttpServerOptions,
	grpcOpts server.GrpcServerOptions,
) server.Setup {
	_ = middlewares
	_ = httpOpts
	_ = grpcOpts
	return nil
}

// Boot server 初始化完成，可以在这里注入 http/grpc/websocket 实例 绑定路由
func Boot(
	httpServer server.HttpServer,
	grpcServer server.GrpcServer,
	wss server.WebsocketServer,
	// service
	exampleWsHandler *service.ExampleWsHandler,
	exampleService *service.ExampleService,
) server.Bootstrap {
	// websocket
	wss.Handle("/echo", exampleWsHandler, websocket.Upgrader{
		// 不校验来源,在websocket在线工具下可以调试
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
	// http
	example_pb.RegisterExampleServiceHTTPServer(httpServer, exampleService)
	// grpc
	example_pb.RegisterExampleServiceServer(grpcServer, exampleService)
	return nil
}
