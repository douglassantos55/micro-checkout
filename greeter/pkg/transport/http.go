package transport

import (
	"net/http"

	"example.com/microservices/greeter/pkg"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHttpServer(svc pkg.Greeter) http.Handler {
	server := http.NewServeMux()

	server.Handle("/greet", makeHttpHandler(svc))
	server.Handle("/greetings", makeHttpHandler(svc))

	return server
}

func makeHttpHandler(greeter pkg.Greeter) *kithttp.Server {
	return kithttp.NewServer(
		pkg.MakeGreetEndpoint(greeter),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
}
