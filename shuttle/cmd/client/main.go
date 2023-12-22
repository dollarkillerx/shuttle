package main

import (
	"net/http"

	proxy_dns "google.dev/google/proxy_dns/cmd/client"
	"google.dev/google/shuttle/core/app/client"
	"google.dev/google/shuttle/utils/log"
)

func main() {
	client.RouterRegister()

	addr := "127.0.0.1:8985"
	log.Infof("GuardLink API Run: %s \n", addr)

	err := proxy_dns.StartLocalDNSService("192.227.234.228:8253", "127.0.0.1:5352")
	if err != nil {
		log.Error(err)
	}

	server := &http.Server{Addr: addr, Handler: nil}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
