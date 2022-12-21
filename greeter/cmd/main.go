package main

import (
	"flag"
	"log"
	"net/http"

	"example.com/microservices/greeter/pkg"
	httpt "example.com/microservices/greeter/pkg/transport"
)

func main() {
	addr := flag.String("addr", ":5353", "address to bind")
	flag.Parse()

	svc := pkg.NewGreeter()
	server := httpt.NewHttpServer(svc)
	log.Print(http.ListenAndServe(*addr, server))
}
