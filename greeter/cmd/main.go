package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"example.com/microservices/greeter/pkg"
	httpt "example.com/microservices/greeter/pkg/transport"
)

func main() {
	addr := flag.String("addr", ":5353", "address to bind")
	auth_addr := flag.String("auth_addr", ":5454", "auth service address")
	flag.Parse()

	svc := pkg.NewGreeter()
	svc = pkg.ProxyingMiddleware(context.Background(), *auth_addr)(svc)

	server := httpt.NewHttpServer(svc)
	log.Print(http.ListenAndServe(*addr, server))
}
