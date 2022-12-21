package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeAuthEndpoint(auth Auth) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return auth.GetAuthenticated(), nil
	}
}
