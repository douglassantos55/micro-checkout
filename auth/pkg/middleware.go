package pkg

import (
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware() endpoint.Middleware {
	return kitjwt.NewParser(func(token *jwt.Token) (any, error) {
		return []byte("your-256-bit-secret"), nil
	}, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)
}
