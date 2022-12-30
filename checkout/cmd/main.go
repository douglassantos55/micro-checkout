package main

import (
	"net/http"
	"os"

	"example.com/microservices/checkout/pkg"
	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	svc := pkg.NewService(
		pkg.NewMemoryRepository(),
		pkg.NewValidator(),
	)

	svc = pkg.NewProxyingMiddleware(
		pkg.LoggingMiddleware(logger, svc),
		pkg.MakeGetProductEndpoint(),
		pkg.MakeReduceStockEndpoint(),
		pkg.MakeProcessPaymentEndpoint(),
	)

	http.ListenAndServe(":80", pkg.MakeHTTPServer(svc, logger))
}
