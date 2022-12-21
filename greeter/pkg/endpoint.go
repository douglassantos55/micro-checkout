package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeGreetEndpoint(greeter Greeter) endpoint.Endpoint {
	return func(ctx context.Context, name any) (any, error) {
		return greeter.Greet(name.(string)), nil
	}
}

type GreetParams struct {
	Name string `json:"name"`
}
