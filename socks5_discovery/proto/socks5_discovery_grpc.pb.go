// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/socks5_discovery.proto

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

// Socks5DiscoveryClient is the client API for Socks5Discovery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Socks5DiscoveryClient interface {
	Discovery(ctx context.Context, in *DiscoveryRequest, opts ...grpc.CallOption) (*DiscoveryResponse, error)
}

type socks5DiscoveryClient struct {
	cc grpc.ClientConnInterface
}

func NewSocks5DiscoveryClient(cc grpc.ClientConnInterface) Socks5DiscoveryClient {
	return &socks5DiscoveryClient{cc}
}

func (c *socks5DiscoveryClient) Discovery(ctx context.Context, in *DiscoveryRequest, opts ...grpc.CallOption) (*DiscoveryResponse, error) {
	out := new(DiscoveryResponse)
	err := c.cc.Invoke(ctx, "/proto.Socks5Discovery/Discovery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Socks5DiscoveryServer is the server API for Socks5Discovery service.
// All implementations must embed UnimplementedSocks5DiscoveryServer
// for forward compatibility
type Socks5DiscoveryServer interface {
	Discovery(context.Context, *DiscoveryRequest) (*DiscoveryResponse, error)
	mustEmbedUnimplementedSocks5DiscoveryServer()
}

// UnimplementedSocks5DiscoveryServer must be embedded to have forward compatible implementations.
type UnimplementedSocks5DiscoveryServer struct {
}

func (UnimplementedSocks5DiscoveryServer) Discovery(context.Context, *DiscoveryRequest) (*DiscoveryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discovery not implemented")
}
func (UnimplementedSocks5DiscoveryServer) mustEmbedUnimplementedSocks5DiscoveryServer() {}

// UnsafeSocks5DiscoveryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Socks5DiscoveryServer will
// result in compilation errors.
type UnsafeSocks5DiscoveryServer interface {
	mustEmbedUnimplementedSocks5DiscoveryServer()
}

func RegisterSocks5DiscoveryServer(s grpc.ServiceRegistrar, srv Socks5DiscoveryServer) {
	s.RegisterService(&Socks5Discovery_ServiceDesc, srv)
}

func _Socks5Discovery_Discovery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoveryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Socks5DiscoveryServer).Discovery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Socks5Discovery/Discovery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Socks5DiscoveryServer).Discovery(ctx, req.(*DiscoveryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Socks5Discovery_ServiceDesc is the grpc.ServiceDesc for Socks5Discovery service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Socks5Discovery_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Socks5Discovery",
	HandlerType: (*Socks5DiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Discovery",
			Handler:    _Socks5Discovery_Discovery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/socks5_discovery.proto",
}