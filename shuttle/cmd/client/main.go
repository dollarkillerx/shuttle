package main

import (
	"log"
	"net/http"

	proxy_dns "google.dev/google/proxy_dns/cmd/client"
	"google.dev/google/shuttle/core/app/client"
)

func main() {
	client.RouterRegister()

	addr := "127.0.0.1:8985"
	log.Println("GuardLink API Run: %s \n", addr)

	err := proxy_dns.StartLocalDNSService("192.227.234.228:8253", "127.0.0.1:5352")
	if err != nil {
		log.Println(err)
	}

	server := &http.Server{Addr: addr, Handler: nil}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
