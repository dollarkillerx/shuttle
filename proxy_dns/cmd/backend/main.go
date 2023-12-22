package main

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/miekg/dns"
	"google.dev/google/proxy_dns/internal"
	"google.dev/google/proxy_dns/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

const serverPem = `
-----BEGIN CERTIFICATE-----
MIIDRzCCAi+gAwIBAgIJAPnLftTUO93IMA0GCSqGSIb3DQEBCwUAMFIxCzAJBgNV
BAYTAkNOMRAwDgYDVQQIDAdTaUNodWFuMRAwDgYDVQQHDAdDaGVuZ2R1MQ0wCwYD
VQQKDARTdGVwMRAwDgYDVQQDDAd0b25naHVhMB4XDTIzMDMxNTAwMTY0N1oXDTMz
MDMxMjAwMTY0N1owUjELMAkGA1UEBhMCQ04xEDAOBgNVBAgMB1NpQ2h1YW4xEDAO
BgNVBAcMB0NoZW5nZHUxDTALBgNVBAoMBFN0ZXAxEDAOBgNVBAMMB3RvbmdodWEw
ggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC5UvqI8Ip9AQagUQLxqjtm
dABhD4998bPx0t8kd2TS+ah1QMoEyJ58QRJbPY9p+jBxtL0qB/ilawdoycT3694d
T9MQUeN5664bsLaGxSK3yKtuDAcxx+hc45Av+uZ0IHxT3cfr8wkfGj068IPQp6hT
bQpfiQqNSarEmyx6X43nazADicArwbFxqgut3h29Vm71DPh36i/2tbFc2w6uU2Sz
mTh+XDEMPgFzqAS7G1VIAFJ7LzxJh+KxKWr1j+0qKldOhKHR0T9BJIlpPuezgTAe
nap2P8uhsvgm0zuB4BR5npscr71db34Ik/QjUxpk3otcwvWVaI6lBMusLevtvAUN
AgMBAAGjIDAeMBwGA1UdEQQVMBOCC3d3dy5wLXBwLmNuhwTAqB9CMA0GCSqGSIb3
DQEBCwUAA4IBAQA/7jv9E36lkVr0ziv3sDj1W0KTqIU042Gp2WE1LBqcbN+Vs5sK
E3j4q/iamA0TZlH4WbtNtjW64eBehrni4dRImi6GxjJQEkWa0zCe7rPLgZFNIU5G
LeZ9YuLlmt4g9Rwz4h43wDzkO9Hf284260fHRh7BcrrzL7lcDZfbYCS1GXsKuXIb
+CIJ8G7K5VGWK1TRWN3+waMGPvQe64U8nAG68H8seA6Uf9fSYqbQAoJtZNX0DP3v
kKups161Zi9BggPrnqQ0QxPHMHDewuea0FFvwmFQ2l0ow8h08YwL7UCxF5GSPnll
gaq8kjq3u2Wb5jImm7epsco6/Mhu5Q4X+Wsv
-----END CERTIFICATE-----
`

const serverKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAuVL6iPCKfQEGoFEC8ao7ZnQAYQ+PffGz8dLfJHdk0vmodUDK
BMiefEESWz2PafowcbS9Kgf4pWsHaMnE9+veHU/TEFHjeeuuG7C2hsUit8irbgwH
McfoXOOQL/rmdCB8U93H6/MJHxo9OvCD0KeoU20KX4kKjUmqxJssel+N52swA4nA
K8GxcaoLrd4dvVZu9Qz4d+ov9rWxXNsOrlNks5k4flwxDD4Bc6gEuxtVSABSey88
SYfisSlq9Y/tKipXToSh0dE/QSSJaT7ns4EwHp2qdj/LobL4JtM7geAUeZ6bHK+9
XW9+CJP0I1MaZN6LXML1lWiOpQTLrC3r7bwFDQIDAQABAoIBAG5Deb8aZzui7Z9b
NAY3g+ocYNFfIcAMnEToc03OH9YLJvjEmK4p82n4iYRx5y9l5Ybxw48LeRxqxtjJ
HAFqfBgyk2DlaBP1bv0YsjETf+mbYqwySeGLkKwb1YFGpfE4FuELVtUDIE06Hm5A
Bh2Sc2tXuFFJR1bzGsCpltgknFiwgMX8Cki8tMz13Us3zvVLGnR+QG2xWJsGveu6
0R0iL3D7jfiyTjlfLhclguTHPBIrykNXQ8f6x0vS46NBdJpl+VHgUCnMRGgD3qKA
/55W9I586yyTSiBiHhzj4yLuZvy7H2L3wnPxexAgzAeMISgkhth0YlwO6/WqSiCq
wbrrqoECgYEA6zvdKQGFJ6iYTbcZM1FPXMUClVfoSPXPoNj2pc6fWmihahY6AJ7r
tXC1pJMFoCoN1otvkPHmIk6mXJNL4J6HXOyL7c2cr8UwOvT9IfFqBhbLytZHaAO2
fN9vnE1ua4/eysof8wspgbOI8jf8KgVE11qjhFar7naxmVH5B2+IOeECgYEAya8v
2S/ZkCY9SbLnOYbBJVkJPtAVS1M7Sezs4UFSCwU+pwvX9MAxPofldiDw4JOUCHw1
nSB61Kk4AxIn2HHApFAXD4/dnydW084/6Drwvq3BEZH/Dgp0yyv1m3hshru9MzOc
qyvLcdVm13PqhZgLUL94eALOnU0vzfuOQ1tZ6K0CgYEAvAJQtSF950Cc2iBph9aI
88CSXAKyqP2uQQSnvcXzHzNZL40sNqrOAWpgA3VunaB/BubS+KoeIXVzCbLAhnqt
/dshy6L2hJW6AqUkXCizcMJvh2LUF5JAHHYIoohQpK+MhdAe0QYu2ndAETgl1v/3
EZhj8LXFHQbI053sx8CgxIECgYBMbKgTAsDMkND0lmhsMhYKkvyf4rXO/1EeKDty
+A+gwXIGVsSUqCeA7HoVE1JzpziXJooiamZhI2ZoM38J08EOApNagEeYwY1zYVpy
I7OKbckVYV9m8KtlOdkt+qoVPBrrxgj+C/BhyF3aEsCxsvXGuWdrApVMoi0VPtef
yoP9WQKBgQC8xy19uJgYBKEK/qGFpCxcsiIRM0azIe5CXXAEyPus6uQMItjGSfX/
s5I0gbcoDWqtfbz54QQLWjQ3o7YEI2rMj++WCPes74jVpUk9wuUJQE9WVudnUMwK
0xcjx4WcX4PUhjddbWKVv0E2IqfO3nvYsyG64111EKtzu8yqHSqHGw==
-----END RSA PRIVATE KEY-----
`

func main() {
	addr := "0.0.0.0:8253"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := loadTLSCredentials([]byte(serverPem), []byte(serverKey))
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer(
		grpc.Creds(creds),
	)

	dnsCache := internal.InitCache(int64(time.Second * 10))
	dnsProxy := internal.NewDNSProxy(&dnsCache, map[string]interface{}{}, map[string]interface{}{}, "1.1.1.1:53")

	pxs := ProxyDnsServer{dnsProxy: dnsProxy}
	proto.RegisterProxyDnsServer(srv, &pxs)

	log.Printf("server listening at %v\n", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type ProxyDnsServer struct {
	proto.UnimplementedProxyDnsServer
	dnsProxy *internal.DNSProxy
}

func (p *ProxyDnsServer) Dig(ctx context.Context, request *proto.DigRequest) (*proto.DigResponse, error) {
	var r = new(dns.Msg)
	err := r.Unpack(request.Data)
	if err != nil {
		return nil, err
	}

	if len(r.Question) > 0 {
		question := r.Question[0].Name
		log.Printf("lookup: %s \n", question)
	}

	response, err := p.dnsProxy.GetResponse(r)
	if err != nil {
		return nil, err
	}

	msg, err := response.Pack()
	if err != nil {
		return nil, err
	}

	return &proto.DigResponse{Data: msg}, nil
}

func (p *ProxyDnsServer) Hello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &proto.HelloResponse{
		Msg: req.Msg + " hello world" + strconv.Itoa(int(r.Int63())),
	}, nil
}

func loadTLSCredentials(certPEMBlock []byte, keyPEMBlock []byte) (credentials.TransportCredentials, error) {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}}), nil
}
