package main

import (
	"flag"
	"log"
	"net/http"

	httpt "example.com/microservices/greeter/pkg/transport"
)

func main() {
	addr := flag.String("addr", ":5353", "address to server")
	flag.Parse()

	server := httpt.NewHttpServer()
	log.Print(http.ListenAndServe(*addr, server))
}
