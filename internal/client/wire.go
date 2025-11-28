package client

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user3"
)

// 与 repo 类似，bind的接口需要在这里约束
var _ user3.GetUser = (*Biz3GetUserImpl)(nil)

var ProviderSet = wire.NewSet(
	example_pb.ExampleServiceApiProvider,
	NewBiz3GetUserImpl,
	wire.Bind(new(user3.GetUser), new(*Biz3GetUserImpl)),
)
