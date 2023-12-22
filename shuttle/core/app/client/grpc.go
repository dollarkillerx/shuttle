package client

import (
	"bytes"
	"context"
	"log"
	"net"

	"google.dev/google/shuttle/proto/grpc_client_plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: 实现...

func (c *Client) wrapGRPC(conn net.Conn) net.Conn {
	return newGRPCWrapper(conn, c)
}

func newGRPCWrapper(conn net.Conn, client *Client) *grpcWrapper {
	rpcConn, err := grpc.Dial(client.Config.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return nil
	}

	pluginClient := grpc_client_plugin.NewGuardLinkGrpcPluginClient(rpcConn)

	stream, err := pluginClient.Stream(context.TODO())

	return &grpcWrapper{
		buf:        bytes.NewBuffer(make([]byte, 0, 1024)),
		Conn:       conn,
		client:     client,
		grpcClient: pluginClient,
		stream:     stream,
	}
}

type grpcWrapper struct {
	net.Conn
	grpcClient grpc_client_plugin.GuardLinkGrpcPluginClient
	buf        *bytes.Buffer
	client     *Client
	stream     grpc_client_plugin.GuardLinkGrpcPlugin_StreamClient
}

func (g *grpcWrapper) Read(b []byte) (n int, err error) {

	if g.buf.Len() > 0 {
		return g.buf.Read(b)
	}

	recv, err := g.stream.Recv()
	if err != nil {
		return 0, err
	}
	n = copy(b, recv.Data)
	g.buf.Write(recv.Data)

	return
}

func (g *grpcWrapper) Write(b []byte) (n int, err error) {

	err = g.stream.Send(&grpc_client_plugin.StreamRequest{
		Data: b,
	})
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (g *grpcWrapper) Close() error {
	err := g.Conn.Close()
	if err != nil {
		return nil
	}

	g.stream.CloseSend()

	return nil
}
