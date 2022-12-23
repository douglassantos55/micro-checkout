package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Verifies auth token before proceeding with request
func AuthenticationMiddleware(authService string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			// Makes a request to the auth service
			// to check token validity
			endpoint := verifyAuthTokenProxy(authService)
			name, err := endpoint(ctx, r)
			if err != nil {
				return nil, err
			}
			// Proceeds with request
			return next(ctx, name)
		}
	}
}
