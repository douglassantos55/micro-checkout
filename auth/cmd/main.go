package main

import (
	"flag"
	"net/http"

	"example.com/microservices/auth/pkg"
)

func main() {
	addr := flag.String("addr", ":80", "address to bind to")
	flag.Parse()

	server := pkg.NewHttpServer(pkg.NewAuth())
	http.ListenAndServe(*addr, server)
}
