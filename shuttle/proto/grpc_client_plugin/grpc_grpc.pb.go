// Code generated by protoc-gen-go-rpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-rpc v1.2.0
// - protoc             v3.21.12
// source: proto/grpc_client_plugin/rpc.proto

package grpc_client_plugin

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the rpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GuardLinkGrpcPluginClient is the client API for GuardLinkGrpcPlugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GuardLinkGrpcPluginClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (GuardLinkGrpcPlugin_StreamClient, error)
}

type guardLinkGrpcPluginClient struct {
	cc grpc.ClientConnInterface
}

func NewGuardLinkGrpcPluginClient(cc grpc.ClientConnInterface) GuardLinkGrpcPluginClient {
	return &guardLinkGrpcPluginClient{cc}
}

func (c *guardLinkGrpcPluginClient) Stream(ctx context.Context, opts ...grpc.CallOption) (GuardLinkGrpcPlugin_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &GuardLinkGrpcPlugin_ServiceDesc.Streams[0], "/proto.GuardLinkGrpcPlugin/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &guardLinkGrpcPluginStreamClient{stream}
	return x, nil
}

type GuardLinkGrpcPlugin_StreamClient interface {
	Send(*StreamRequest) error
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type guardLinkGrpcPluginStreamClient struct {
	grpc.ClientStream
}

func (x *guardLinkGrpcPluginStreamClient) Send(m *StreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *guardLinkGrpcPluginStreamClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GuardLinkGrpcPluginServer is the server API for GuardLinkGrpcPlugin service.
// All implementations must embed UnimplementedGuardLinkGrpcPluginServer
// for forward compatibility
type GuardLinkGrpcPluginServer interface {
	Stream(GuardLinkGrpcPlugin_StreamServer) error
	mustEmbedUnimplementedGuardLinkGrpcPluginServer()
}

// UnimplementedGuardLinkGrpcPluginServer must be embedded to have forward compatible implementations.
type UnimplementedGuardLinkGrpcPluginServer struct {
}

func (UnimplementedGuardLinkGrpcPluginServer) Stream(GuardLinkGrpcPlugin_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedGuardLinkGrpcPluginServer) mustEmbedUnimplementedGuardLinkGrpcPluginServer() {}

// UnsafeGuardLinkGrpcPluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GuardLinkGrpcPluginServer will
// result in compilation errors.
type UnsafeGuardLinkGrpcPluginServer interface {
	mustEmbedUnimplementedGuardLinkGrpcPluginServer()
}

func RegisterGuardLinkGrpcPluginServer(s grpc.ServiceRegistrar, srv GuardLinkGrpcPluginServer) {
	s.RegisterService(&GuardLinkGrpcPlugin_ServiceDesc, srv)
}

func _GuardLinkGrpcPlugin_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GuardLinkGrpcPluginServer).Stream(&guardLinkGrpcPluginStreamServer{stream})
}

type GuardLinkGrpcPlugin_StreamServer interface {
	Send(*StreamResponse) error
	Recv() (*StreamRequest, error)
	grpc.ServerStream
}

type guardLinkGrpcPluginStreamServer struct {
	grpc.ServerStream
}

func (x *guardLinkGrpcPluginStreamServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *guardLinkGrpcPluginStreamServer) Recv() (*StreamRequest, error) {
	m := new(StreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GuardLinkGrpcPlugin_ServiceDesc is the grpc.ServiceDesc for GuardLinkGrpcPlugin service.
// It's only intended for direct use with rpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GuardLinkGrpcPlugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GuardLinkGrpcPlugin",
	HandlerType: (*GuardLinkGrpcPluginServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _GuardLinkGrpcPlugin_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/grpc_client_plugin/rpc.proto",
}
