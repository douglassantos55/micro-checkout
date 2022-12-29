package main

import (
	"net/http"
	"os"

	"example.com/microservices/customer/pkg"
	"github.com/go-kit/log"
)

func main() {
	repo := pkg.NewInMemoryRepository()
	svc := pkg.NewService(repo, pkg.NewValidator())

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	go pkg.NewMessageBroker(logger, "guest", "guest", "messaging-service")

	server := pkg.MakeHTTPServer(svc)
	http.ListenAndServe(":80", server)
}
