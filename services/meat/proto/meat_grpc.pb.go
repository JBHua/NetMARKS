// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: meat.proto

package proto

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

// MeatClient is the client API for Meat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MeatClient interface {
	Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type meatClient struct {
	cc grpc.ClientConnInterface
}

func NewMeatClient(cc grpc.ClientConnInterface) MeatClient {
	return &meatClient{cc}
}

func (c *meatClient) Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/netmarks_meat.Meat/Produce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MeatServer is the server API for Meat service.
// All implementations must embed UnimplementedMeatServer
// for forward compatibility
type MeatServer interface {
	Produce(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedMeatServer()
}

// UnimplementedMeatServer must be embedded to have forward compatible implementations.
type UnimplementedMeatServer struct {
}

func (UnimplementedMeatServer) Produce(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Produce not implemented")
}
func (UnimplementedMeatServer) mustEmbedUnimplementedMeatServer() {}

// UnsafeMeatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MeatServer will
// result in compilation errors.
type UnsafeMeatServer interface {
	mustEmbedUnimplementedMeatServer()
}

func RegisterMeatServer(s grpc.ServiceRegistrar, srv MeatServer) {
	s.RegisterService(&Meat_ServiceDesc, srv)
}

func _Meat_Produce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeatServer).Produce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/netmarks_meat.Meat/Produce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeatServer).Produce(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Meat_ServiceDesc is the grpc.ServiceDesc for Meat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Meat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "netmarks_meat.Meat",
	HandlerType: (*MeatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Produce",
			Handler:    _Meat_Produce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "meat.proto",
}
