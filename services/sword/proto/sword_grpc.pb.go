// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: sword.proto

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

// SwordClient is the client API for Sword service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SwordClient interface {
	Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type swordClient struct {
	cc grpc.ClientConnInterface
}

func NewSwordClient(cc grpc.ClientConnInterface) SwordClient {
	return &swordClient{cc}
}

func (c *swordClient) Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/netmarks_sword.Sword/Produce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SwordServer is the server API for Sword service.
// All implementations must embed UnimplementedSwordServer
// for forward compatibility
type SwordServer interface {
	Produce(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedSwordServer()
}

// UnimplementedSwordServer must be embedded to have forward compatible implementations.
type UnimplementedSwordServer struct {
}

func (UnimplementedSwordServer) Produce(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Produce not implemented")
}
func (UnimplementedSwordServer) mustEmbedUnimplementedSwordServer() {}

// UnsafeSwordServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SwordServer will
// result in compilation errors.
type UnsafeSwordServer interface {
	mustEmbedUnimplementedSwordServer()
}

func RegisterSwordServer(s grpc.ServiceRegistrar, srv SwordServer) {
	s.RegisterService(&Sword_ServiceDesc, srv)
}

func _Sword_Produce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwordServer).Produce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/netmarks_sword.Sword/Produce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwordServer).Produce(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Sword_ServiceDesc is the grpc.ServiceDesc for Sword service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sword_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "netmarks_sword.Sword",
	HandlerType: (*SwordServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Produce",
			Handler:    _Sword_Produce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sword.proto",
}
