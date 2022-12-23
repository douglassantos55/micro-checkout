package pkg

import (
	"context"
	"errors"
	"fmt"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v4"
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

// This endpoint does not require a service call
func makeVerifyEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		claims, ok := ctx.Value(kitjwt.JWTClaimsContextKey).(*jwt.StandardClaims)
		if !ok {
			return nil, errors.New(fmt.Sprintf("could not convert claims: %v", claims))
		}
		return claims.Subject, nil
	}
}
