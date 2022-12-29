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

	broker, err := pkg.NewRabbitMQBroker("guest", "guest", "messaging-service")
	if err != nil {
		panic(err)
	}

	svc := pkg.NewService(
		pkg.NewMemoryRepository(),
		pkg.NewValidator(),
		broker,
	)

	svc = pkg.LoggingMiddleware(logger, svc)
	http.ListenAndServe(":80", pkg.MakeHTTPServer(svc, logger))
}
