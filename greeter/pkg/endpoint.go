package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeGreetEndpoint(greeter Greeter) endpoint.Endpoint {
	return func(ctx context.Context, name any) (any, error) {
		authEndpoint := getAuthenticatedUser(":5454")
		user, err := authEndpoint(ctx, nil)

		if err != nil {
			return nil, err
		}

		return greeter.Greet(user.(string)), nil
	}
}
