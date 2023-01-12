// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: plugin.proto

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

// CalculateClient is the client API for Calculate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalculateClient interface {
	Calculate(ctx context.Context, in *CalcRequest, opts ...grpc.CallOption) (*ResultResponse, error)
}

type calculateClient struct {
	cc grpc.ClientConnInterface
}

func NewCalculateClient(cc grpc.ClientConnInterface) CalculateClient {
	return &calculateClient{cc}
}

func (c *calculateClient) Calculate(ctx context.Context, in *CalcRequest, opts ...grpc.CallOption) (*ResultResponse, error) {
	out := new(ResultResponse)
	err := c.cc.Invoke(ctx, "/proto.Calculate/Calculate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalculateServer is the server API for Calculate service.
// All implementations must embed UnimplementedCalculateServer
// for forward compatibility
type CalculateServer interface {
	Calculate(context.Context, *CalcRequest) (*ResultResponse, error)
	mustEmbedUnimplementedCalculateServer()
}

// UnimplementedCalculateServer must be embedded to have forward compatible implementations.
type UnimplementedCalculateServer struct {
}

func (UnimplementedCalculateServer) Calculate(context.Context, *CalcRequest) (*ResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Calculate not implemented")
}
func (UnimplementedCalculateServer) mustEmbedUnimplementedCalculateServer() {}

// UnsafeCalculateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalculateServer will
// result in compilation errors.
type UnsafeCalculateServer interface {
	mustEmbedUnimplementedCalculateServer()
}

func RegisterCalculateServer(s grpc.ServiceRegistrar, srv CalculateServer) {
	s.RegisterService(&Calculate_ServiceDesc, srv)
}

func _Calculate_Calculate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalcRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalculateServer).Calculate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Calculate/Calculate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalculateServer).Calculate(ctx, req.(*CalcRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Calculate_ServiceDesc is the grpc.ServiceDesc for Calculate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calculate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Calculate",
	HandlerType: (*CalculateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Calculate",
			Handler:    _Calculate_Calculate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin.proto",
}