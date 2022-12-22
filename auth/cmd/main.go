package main

import (
	"flag"
	"net/http"

	"example.com/microservices/auth/pkg"
)

func main() {
	addr := flag.String("addr", ":80", "address to bind to")
	flag.Parse()

	encoder := pkg.NewJWTEncoder()
	svc := pkg.NewAuth(encoder)
	server := pkg.NewHttpServer(svc)

	http.ListenAndServe(*addr, server)
}
