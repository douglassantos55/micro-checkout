package main

import (
	"net/http"

	"example.com/microservices/customer/pkg"
)

func main() {
	repo := pkg.NewInMemoryRepository()
	svc := pkg.NewService(repo, pkg.NewValidator())

	server := pkg.MakeHTTPServer(svc)
	http.ListenAndServe(":80", server)
}
