package main

import (
	"net/http"

	"example.com/microservices/product/pkg"
)

func main() {
	svc := pkg.NewService(pkg.NewMemoryRepository())
	server := pkg.MakeHTTPHandler(svc)
	http.ListenAndServe(":80", server)
}
