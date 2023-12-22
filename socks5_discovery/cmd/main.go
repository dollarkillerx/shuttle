package main

import (
	"google.dev/google/socks5_discovery/internal/conf"
	"google.dev/google/socks5_discovery/internal/service"
	"google.dev/google/socks5_discovery/proto"
	"google.golang.org/grpc"

	"log"
	"net"
)

func main() {
	cf := new(conf.S5DiscoveryConfig)
	cf.ReadConf()

	ser := service.NewService(*cf)

	listen, err := net.Listen("tcp", cf.Listen)
	if err != nil {
		panic(err)
	}
	g := grpc.NewServer()
	proto.RegisterSocks5DiscoveryServer(g, ser)
	log.Println("GRPC: ", cf.Listen)
	if err := g.Serve(listen); err != nil {
		panic(err)
	}
}
