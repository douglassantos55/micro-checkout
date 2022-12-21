package transport

import (
	"context"
	"net/http"

	"example.com/microservices/greeter/pkg"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHttpServer() http.Handler {
	svc := pkg.NewGreeter()

	server := http.NewServeMux()
	server.Handle("/greet", makeHttpHandler(svc))
	server.Handle("/greetings", makeHttpHandler(svc))

	return server
}

func makeHttpHandler(greeter pkg.Greeter) *kithttp.Server {
	return kithttp.NewServer(
		pkg.MakeGreetEndpoint(greeter),
		decodeRequest,
		kithttp.EncodeJSONResponse,
	)
}

func decodeRequest(ctx context.Context, req *http.Request) (any, error) {
	return req.URL.Query().Get("name"), nil
}
