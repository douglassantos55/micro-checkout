package main

import (
	"net/http"

	"example.com/microservices/customer/pkg"
)

func main() {
	repo := pkg.NewInMemoryRepository()
	svc := pkg.NewService(repo)

	server := pkg.MakeHTTPServer(svc)
	http.ListenAndServe(":80", server)
}
