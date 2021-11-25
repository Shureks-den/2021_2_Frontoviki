// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package category

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CategoryClient is the client API for Category service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CategoryClient interface {
	GetCategories(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Categories, error)
}

type categoryClient struct {
	cc grpc.ClientConnInterface
}

func NewCategoryClient(cc grpc.ClientConnInterface) CategoryClient {
	return &categoryClient{cc}
}

func (c *categoryClient) GetCategories(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Categories, error) {
	out := new(Categories)
	err := c.cc.Invoke(ctx, "/category.Category/GetCategories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CategoryServer is the server API for Category service.
// All implementations should embed UnimplementedCategoryServer
// for forward compatibility
type CategoryServer interface {
	GetCategories(context.Context, *Nothing) (*Categories, error)
}

// UnimplementedCategoryServer should be embedded to have forward compatible implementations.
type UnimplementedCategoryServer struct {
}

func (UnimplementedCategoryServer) GetCategories(context.Context, *Nothing) (*Categories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategories not implemented")
}

// UnsafeCategoryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CategoryServer will
// result in compilation errors.
type UnsafeCategoryServer interface {
	mustEmbedUnimplementedCategoryServer()
}

func RegisterCategoryServer(s grpc.ServiceRegistrar, srv CategoryServer) {
	s.RegisterService(&Category_ServiceDesc, srv)
}

func _Category_GetCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServer).GetCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/category.Category/GetCategories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServer).GetCategories(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

// Category_ServiceDesc is the grpc.ServiceDesc for Category service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Category_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "category.Category",
	HandlerType: (*CategoryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCategories",
			Handler:    _Category_GetCategories_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "category.proto",
}
