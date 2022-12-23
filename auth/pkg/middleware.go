package pkg

import (
	"time"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware() endpoint.Middleware {
	return kitjwt.NewParser(func(token *jwt.Token) (any, error) {
		return []byte("your-256-bit-secret"), nil
	}, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)
}

type loggingAuth struct {
	svc    Auth
	logger log.Logger
}

func (s *loggingAuth) Authenticate(credentials Credentials) (user string, err error) {
	defer func(timestamp time.Time) {
		s.logger.Log(
			"timestamp", timestamp,
			"method", "Authenticate",
			"credentials", credentials,
			"user", user,
			"error", err,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return s.svc.Authenticate(credentials)
}

func LoggingMiddleware(svc Auth, logger log.Logger) Auth {
	return &loggingAuth{svc, logger}
}
