package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"example.com/microservices/greeter/pkg"
	httpt "example.com/microservices/greeter/pkg/transport"
)

func main() {
	addr := flag.String("addr", ":80", "address to bind")
	flag.Parse()

	authService, ok := os.LookupEnv("AUTH_SERVICE_ADDR")
	if !ok {
		log.Fatal("auth_service_addr not found")
	}

	svc := pkg.NewGreeter()
	server := httpt.NewHttpServer(svc, authService)
	log.Print(http.ListenAndServe(*addr, server))
}
