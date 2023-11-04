// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: water.proto

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

// WaterClient is the client API for Water service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WaterClient interface {
	ProduceWater(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Single, error)
}

type waterClient struct {
	cc grpc.ClientConnInterface
}

func NewWaterClient(cc grpc.ClientConnInterface) WaterClient {
	return &waterClient{cc}
}

func (c *waterClient) ProduceWater(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Single, error) {
	out := new(Single)
	err := c.cc.Invoke(ctx, "/netmarks_grain.Water/ProduceWater", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WaterServer is the server API for Water service.
// All implementations must embed UnimplementedWaterServer
// for forward compatibility
type WaterServer interface {
	ProduceWater(context.Context, *Request) (*Single, error)
	mustEmbedUnimplementedWaterServer()
}

// UnimplementedWaterServer must be embedded to have forward compatible implementations.
type UnimplementedWaterServer struct {
}

func (UnimplementedWaterServer) ProduceWater(context.Context, *Request) (*Single, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceWater not implemented")
}
func (UnimplementedWaterServer) mustEmbedUnimplementedWaterServer() {}

// UnsafeWaterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WaterServer will
// result in compilation errors.
type UnsafeWaterServer interface {
	mustEmbedUnimplementedWaterServer()
}

func RegisterWaterServer(s grpc.ServiceRegistrar, srv WaterServer) {
	s.RegisterService(&Water_ServiceDesc, srv)
}

func _Water_ProduceWater_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WaterServer).ProduceWater(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/netmarks_grain.Water/ProduceWater",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WaterServer).ProduceWater(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Water_ServiceDesc is the grpc.ServiceDesc for Water service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Water_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "netmarks_grain.Water",
	HandlerType: (*WaterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProduceWater",
			Handler:    _Water_ProduceWater_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "water.proto",
}
