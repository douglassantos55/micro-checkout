package main

import (
	"net/http"
	"os"
	"time"

	"example.com/microservices/payment/pkg"
	"github.com/go-kit/kit/log"
)

func main() {
	repo := pkg.NewMemoryRepository()

	writter := log.NewSyncWriter(os.Stderr)
	logger := log.NewJSONLogger(writter)
	logger = log.With(logger, "ts", time.Now())

	svc := pkg.NewService(repo)
	svc = pkg.LoggingMiddleware(logger, svc)

    go pkg.MakeAMQPSubscriber(svc)

	server := pkg.MakeHTTPServer(svc, logger)
	http.ListenAndServe(":80", server)
}
