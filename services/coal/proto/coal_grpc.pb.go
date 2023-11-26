// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: coal.proto

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

// CoalClient is the client API for Coal service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoalClient interface {
	Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type coalClient struct {
	cc grpc.ClientConnInterface
}

func NewCoalClient(cc grpc.ClientConnInterface) CoalClient {
	return &coalClient{cc}
}

func (c *coalClient) Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/netmarks_coal.Coal/Produce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CoalServer is the server API for Coal service.
// All implementations must embed UnimplementedCoalServer
// for forward compatibility
type CoalServer interface {
	Produce(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedCoalServer()
}

// UnimplementedCoalServer must be embedded to have forward compatible implementations.
type UnimplementedCoalServer struct {
}

func (UnimplementedCoalServer) Produce(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Produce not implemented")
}
func (UnimplementedCoalServer) mustEmbedUnimplementedCoalServer() {}

// UnsafeCoalServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoalServer will
// result in compilation errors.
type UnsafeCoalServer interface {
	mustEmbedUnimplementedCoalServer()
}

func RegisterCoalServer(s grpc.ServiceRegistrar, srv CoalServer) {
	s.RegisterService(&Coal_ServiceDesc, srv)
}

func _Coal_Produce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoalServer).Produce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/netmarks_coal.Coal/Produce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoalServer).Produce(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Coal_ServiceDesc is the grpc.ServiceDesc for Coal service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Coal_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "netmarks_coal.Coal",
	HandlerType: (*CoalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Produce",
			Handler:    _Coal_Produce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coal.proto",
}