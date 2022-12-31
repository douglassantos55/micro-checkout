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

	svc := pkg.NewService(pkg.NewMemoryRepository(), pkg.NewValidator())
	svc = pkg.LoggingMiddleware(logger, pkg.MakeSubtotalProxy(svc, pkg.MakeGetProductEndpoint()))
	svc = pkg.LoggingMiddleware(logger, pkg.MakeStockProxy(svc, pkg.MakeReduceStockEndpoint()))
	svc = pkg.LoggingMiddleware(logger, pkg.MakePaymentProxy(svc, pkg.MakeProcessPaymentEndpoint()))

	http.ListenAndServe(":80", pkg.MakeHTTPServer(svc, logger))
}
