package main

import (
	"net/http"
	"os"
	"time"

	"example.com/microservices/checkout/pkg"
	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", time.Now())

	svc := pkg.NewService(
		pkg.NewMemoryRepository(),
		pkg.NewValidator(),
	)

	svc = pkg.LoggingMiddleware(logger, svc)
	http.ListenAndServe(":80", pkg.MakeHTTPServer(svc, logger))
}
