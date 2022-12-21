package pkg

import (
	"context"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v4"
)

func makeAuthEndpoint(svc Auth) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		claims := ctx.Value(kitjwt.JWTClaimsContextKey).(*jwt.StandardClaims)
		return svc.GetAuthenticated(claims.Subject), nil
	}
}
