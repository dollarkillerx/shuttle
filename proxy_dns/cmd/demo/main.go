package main

import (
	"log"
	"regexp"
	"time"

	"github.com/miekg/dns"
	"google.dev/google/proxy_dns/internal"
)

func main() {
	dnsCache := internal.InitCache(int64(time.Second * 10))
	dnsProxy := internal.NewDNSProxy(&dnsCache, map[string]interface{}{}, map[string]interface{}{}, "8.8.8.8:53")

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		switch r.Opcode {
		case dns.OpcodeQuery:
			m, err := dnsProxy.GetResponse(r)
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
			}
			m.SetReply(r)
			w.WriteMsg(m)
		}
	})

	server := &dns.Server{Addr: "127.0.0.1:53", Net: "udp"}
	log.Printf("Starting at %s\n", "127.0.0.1:53")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
