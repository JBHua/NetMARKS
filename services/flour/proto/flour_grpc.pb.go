// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: flour.proto

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

// FlourClient is the client API for Flour service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FlourClient interface {
	Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type flourClient struct {
	cc grpc.ClientConnInterface
}

func NewFlourClient(cc grpc.ClientConnInterface) FlourClient {
	return &flourClient{cc}
}

func (c *flourClient) Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/netmarks_flour.Flour/Produce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FlourServer is the server API for Flour service.
// All implementations must embed UnimplementedFlourServer
// for forward compatibility
type FlourServer interface {
	Produce(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedFlourServer()
}

// UnimplementedFlourServer must be embedded to have forward compatible implementations.
type UnimplementedFlourServer struct {
}

func (UnimplementedFlourServer) Produce(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Produce not implemented")
}
func (UnimplementedFlourServer) mustEmbedUnimplementedFlourServer() {}

// UnsafeFlourServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FlourServer will
// result in compilation errors.
type UnsafeFlourServer interface {
	mustEmbedUnimplementedFlourServer()
}

func RegisterFlourServer(s grpc.ServiceRegistrar, srv FlourServer) {
	s.RegisterService(&Flour_ServiceDesc, srv)
}

func _Flour_Produce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlourServer).Produce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/netmarks_flour.Flour/Produce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlourServer).Produce(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Flour_ServiceDesc is the grpc.ServiceDesc for Flour service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Flour_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "netmarks_flour.Flour",
	HandlerType: (*FlourServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Produce",
			Handler:    _Flour_Produce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "flour.proto",
}
