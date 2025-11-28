package service

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user1"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user2"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example/user3"
)

type ExampleService struct {
	example_pb.UnimplementedExampleServiceServer

	biz1 *user1.User1Biz
	biz2 *user2.User2Biz
	biz3 *user3.User3Biz
}

func NewExampleService(
	biz1 *user1.User1Biz,
	biz2 *user2.User2Biz,
	biz3 *user3.User3Biz,
) *ExampleService {
	return &ExampleService{
		biz1: biz1,
		biz2: biz2,
		biz3: biz3,
	}
}

func (s *ExampleService) GetUser(ctx context.Context, in *example_pb.GetUserRequest) (*example_pb.GetUserResponse, error) {
	usr, err := s.biz1.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &example_pb.GetUserResponse{
		Data: bizUserToPbUser(usr),
	}, nil
}

func (s *ExampleService) GetUser2(ctx context.Context, in *example_pb.GetUserRequest) (*example_pb.GetUserResponse, error) {
	usr, err := s.biz2.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &example_pb.GetUserResponse{
		Data: bizUserToPbUser(usr),
	}, nil
}

func (s *ExampleService) GetUser3(ctx context.Context, in *example_pb.GetUserRequest) (*example_pb.GetUserResponse, error) {
	usr, err := s.biz3.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &example_pb.GetUserResponse{
		Data: bizUserToPbUser(usr),
	}, nil
}

func bizUserToPbUser(user *example.User) *example_pb.GetUserResponse_Data {
	return &example_pb.GetUserResponse_Data{
		Id:   user.ID,
		Name: user.Name,
	}
}
