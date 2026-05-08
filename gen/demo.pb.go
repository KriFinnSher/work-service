package gen

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type GetUserRequest struct {
	Id string
}

type GetUserResponse struct {
	Id    string
	Email string
	Name  string
}

type CreateUserRequest struct {
	Email string
	Name  string
}

type CreateUserResponse struct {
	Id string
}

type DemoServiceServer interface {
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
}

type UnimplementedDemoServiceServer struct{}

func (UnimplementedDemoServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (UnimplementedDemoServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func RegisterDemoServiceServer(s *grpc.Server, srv DemoServiceServer) {
	s.RegisterService(&_DemoService_serviceDesc, srv)
}

var _DemoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "demo.v1.DemoService",
	HandlerType: (*DemoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _DemoService_GetUser_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _DemoService_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/demo.proto",
}

func _DemoService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	var in GetUserRequest
	if err := dec(&in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DemoServiceServer).GetUser(ctx, &in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.v1.DemoService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DemoServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, &in, info, handler)
}

func _DemoService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	var in CreateUserRequest
	if err := dec(&in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DemoServiceServer).CreateUser(ctx, &in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.v1.DemoService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DemoServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, &in, info, handler)
}