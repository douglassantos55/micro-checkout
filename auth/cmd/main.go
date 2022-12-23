package main

import (
	"flag"
	"net/http"
	"os"

	"example.com/microservices/auth/pkg"
	"github.com/go-kit/log"
)

func main() {
	addr := flag.String("addr", ":80", "address to bind to")
	flag.Parse()

	writter := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(writter)

	encoder := pkg.NewJWTEncoder()
	svc := pkg.NewAuth(encoder)
	svc = pkg.LoggingMiddleware(svc, logger)

	server := pkg.NewHttpServer(svc, logger)
	http.ListenAndServe(*addr, server)
}
