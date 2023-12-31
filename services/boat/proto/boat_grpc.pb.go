// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: boat.proto

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

// BoatClient is the client API for Boat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BoatClient interface {
	Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type boatClient struct {
	cc grpc.ClientConnInterface
}

func NewBoatClient(cc grpc.ClientConnInterface) BoatClient {
	return &boatClient{cc}
}

func (c *boatClient) Produce(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/netmarks_boat.Boat/Produce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BoatServer is the server API for Boat service.
// All implementations must embed UnimplementedBoatServer
// for forward compatibility
type BoatServer interface {
	Produce(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedBoatServer()
}

// UnimplementedBoatServer must be embedded to have forward compatible implementations.
type UnimplementedBoatServer struct {
}

func (UnimplementedBoatServer) Produce(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Produce not implemented")
}
func (UnimplementedBoatServer) mustEmbedUnimplementedBoatServer() {}

// UnsafeBoatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BoatServer will
// result in compilation errors.
type UnsafeBoatServer interface {
	mustEmbedUnimplementedBoatServer()
}

func RegisterBoatServer(s grpc.ServiceRegistrar, srv BoatServer) {
	s.RegisterService(&Boat_ServiceDesc, srv)
}

func _Boat_Produce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoatServer).Produce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/netmarks_boat.Boat/Produce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoatServer).Produce(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Boat_ServiceDesc is the grpc.ServiceDesc for Boat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Boat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "netmarks_boat.Boat",
	HandlerType: (*BoatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Produce",
			Handler:    _Boat_Produce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "boat.proto",
}
