package pkg

import (
	"context"
	"errors"
	"os"

	"github.com/go-kit/kit/endpoint"
)

func MakeGreetEndpoint(greeter Greeter) endpoint.Endpoint {
	return func(ctx context.Context, name any) (any, error) {
		authServiceAddr, ok := os.LookupEnv("AUTH_SERVICE_ADDR")
		if !ok {
			return nil, errors.New("AUTH_SERVICE_ADDR not found")
		}

		authEndpoint := getAuthenticatedUser(authServiceAddr)
		user, err := authEndpoint(ctx, nil)

		if err != nil {
			return nil, err
		}

		return greeter.Greet(user.(string)), nil
	}
}
