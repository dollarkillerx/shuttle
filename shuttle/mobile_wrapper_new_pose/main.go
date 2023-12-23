package main

import (
	"google.dev/google/shuttle/core/app/client"
	"google.dev/google/shuttle/utils/log"

	"net/http"
)

func main() {
	log.SetNoLogger()
	defer log.Sync()

	client.RouterRegister()

	addr := "127.0.0.1:8985"
	log.Printf("GuardLink API Run: %s \n", addr)

	server := &http.Server{Addr: addr, Handler: nil}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
