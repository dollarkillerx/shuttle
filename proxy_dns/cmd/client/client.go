package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/miekg/dns"
	"google.dev/google/proxy_dns/internal"
	"google.dev/google/proxy_dns/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const clientPem = `
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

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

var proxyDNSClient proto.ProxyDnsClient

func StartLocalDNSService(remoteAddress string, localAddress string) error {
	if proxyDNSClient == nil {
		creds, err := loadTLSCredentials([]byte(clientPem), "www.p-pp.cn")
		if err != nil {
			log.Fatalln(err)
		}
		conn, err := grpc.Dial(remoteAddress, grpc.WithTransportCredentials(creds))
		if err != nil {
			return err
		}
		proxyDNSClient = proto.NewProxyDnsClient(conn)
	}

	dnsCache := internal.InitCache(int64(time.Minute))
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		switch r.Opcode {
		case dns.OpcodeQuery:
			if len(r.Question) > 0 {
				question := r.Question[0]
				cacheMsg, found := dnsCache.Get(fmt.Sprintf("%s-%d", question.Name, question.Qtype))
				if found {
					rr := cacheMsg.(*dns.Msg)
					rr.SetReply(r)
					w.WriteMsg(rr)
					return
				}
			}

			msg, err := r.Pack()
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n", r, err.Error())
				w.WriteMsg(r)
				return
			}

			dig, err := proxyDNSClient.Dig(context.TODO(), &proto.DigRequest{
				Data: msg,
			})
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n", r, err.Error())
				w.WriteMsg(r)
				return
			}

			var m = new(dns.Msg)
			err = m.Unpack(dig.Data)
			if err != nil {
				log.Printf("Failed lookup for %s with error: %s\n", r, err.Error())
				m.SetReply(r)
				w.WriteMsg(m)
				return
			}
			if len(m.Answer) > 0 {
				pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
				ipAddress := pattern.FindAllString(m.Answer[0].String(), -1)

				if len(ipAddress) > 0 {
					log.Printf("Lookup for %s with ip %s\n", m.Answer[0].Header().Name, ipAddress[0])
				} else {
					log.Printf("Lookup for %s with response %s\n", m.Answer[0].Header().Name, m.Answer[0])
				}

				if len(r.Question) > 0 {
					question := r.Question[0]
					dnsCache.Set(fmt.Sprintf("%s-%d", question.Name, question.Qtype), m)
				}
			}

			m.SetReply(r)
			w.WriteMsg(m)
		}
	})

	server := &dns.Server{Addr: localAddress, Net: "udp"}
	log.Printf("Starting at %s\n", localAddress)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n ", err.Error())
		}
	}()

	return nil
}

func loadTLSCredentials(caCert []byte, serverNameOverride string) (credentials.TransportCredentials, error) {
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}
