package bootstrap

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/jaggerzhuang1994/kratos-foundation-template/api/helloworld/v1"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/app"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/log"
)

func NewBootstrap(
	log *log.Log,
	hookMgr *app.HookManager,
	appHook *AppHook,
	httpServer *http.Server,
	grpcServer *grpc.Server,
	service *service.GreeterService,
) bootstrap.Bootstrap {
	log.Info("boot")
	hookMgr.Register(appHook)

	v1.RegisterGreeterHTTPServer(httpServer, service)
	v1.RegisterGreeterServer(grpcServer, service)
	return nil
}
