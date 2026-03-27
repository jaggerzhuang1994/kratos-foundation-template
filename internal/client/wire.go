package client

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user3"
)

// 与 repo 类似，bind的接口需要在这里约束
var _ user3.GetUser = (*Biz3GetUserImpl)(nil)

var ProviderSet = wire.NewSet(
	// 外部服务pb的provider
	example_pb.ExampleServiceApiProvider,

	// 自定义业务biz接口的实现
	NewBiz3GetUserImpl,
	wire.Bind(new(user3.GetUser), new(*Biz3GetUserImpl)),
)
