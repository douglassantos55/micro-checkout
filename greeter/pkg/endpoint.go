package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

func MakeGreetEndpoint(greeter Greeter) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		name, ok := r.(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("could not convert data to string: %v", r))
		}
		return greeter.Greet(name)
	}
}
