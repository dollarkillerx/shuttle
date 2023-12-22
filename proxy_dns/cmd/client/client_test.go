package client

import (
	"context"
	"fmt"
	"google.dev/google/proxy_dns/proto"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestStartLocalDNSService(t *testing.T) {
	err := StartLocalDNSService("192.227.234.228:8253", "0.0.0.0:8822")
	//err := StartLocalDNSService("127.0.0.1:8253", "0.0.0.0:8822")
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second)
	}
}

func TestStartLocalDNSService2(t *testing.T) {
	creds, err := loadTLSCredentials([]byte(clientPem), "www.p-pp.cn")
	if err != nil {
		panic(err)
	}
	//conn, err := rpc.Dial("192.227.234.228:8253", rpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial("192.227.234.228:8253", grpc.WithTransportCredentials(creds))
	//conn, err := rpc.Dial("192.168.31.66:8253", rpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err)
	//}
	proxyDNSClient := proto.NewProxyDnsClient(conn)
	hello, err := proxyDNSClient.Hello(context.TODO(), &proto.HelloRequest{
		Msg: " dollk ",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(hello.Msg)
}
