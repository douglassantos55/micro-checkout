package main

import (
	"net/http"
	"os"

	"example.com/microservices/product/pkg"
	"github.com/go-kit/log"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	svc := pkg.NewService(pkg.NewMemoryRepository())
	svc = pkg.LoggingMiddleware(logger, svc)

	server := pkg.MakeHTTPHandler(svc, logger)
	http.ListenAndServe(":80", server)
}
