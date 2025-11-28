package client

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example"
)

// 实现业务需要的接口

type Biz3GetUserImpl struct {
	exampleServiceApi example_pb.ExampleServiceApi
}

func NewBiz3GetUserImpl(exampleServiceApi example_pb.ExampleServiceApi) *Biz3GetUserImpl {
	return &Biz3GetUserImpl{
		exampleServiceApi,
	}
}

func (c *Biz3GetUserImpl) GetUser(ctx context.Context, id int64) (*example.User, error) {
	rsp, err := c.exampleServiceApi.GetUser(ctx, &example_pb.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &example.User{
		ID:   rsp.GetData().GetId(),
		Name: rsp.GetData().GetName(),
	}, nil
}
