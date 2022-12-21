package transport

import (
	"context"
	"encoding/json"
	"fmt"
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
		decodeGreetRequest,
		encodeGreetResponse,
	)
}

func decodeGreetRequest(ctx context.Context, req *http.Request) (any, error) {
	return req.URL.Query().Get("name"), nil
}

func encodeGreetResponse(ctx context.Context, w http.ResponseWriter, res any) error {
	fmt.Printf("res: %v\n", res)
	return json.NewEncoder(w).Encode(res)
}
