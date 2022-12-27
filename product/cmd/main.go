package main

import (
	"net/http"
	"os"
	"time"

	"example.com/microservices/product/pkg"
	"github.com/go-kit/log"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", time.Now())

	svc := pkg.NewService(pkg.NewMemoryRepository())
	svc = pkg.LoggingMiddleware(logger, svc)

	server := pkg.MakeHTTPHandler(svc, logger)
	http.ListenAndServe(":80", server)
}
