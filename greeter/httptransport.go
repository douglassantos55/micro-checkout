package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	transport "github.com/go-kit/kit/transport/http"
)

func MakeHttpHandler(greeter Greeter) *transport.Server {
	return transport.NewServer(
		MakeGreetEndpoint(greeter),
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
