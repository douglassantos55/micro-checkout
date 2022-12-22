package pkg

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

func makeAuthEndpoint(svc Auth) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		credentials, ok := r.(Credentials)
		if !ok {
			return nil, errors.New("invalid parameters (credentials)")
		}
		return svc.Authenticate(credentials)
	}
}
