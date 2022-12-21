package pkg

import (
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHttpServer(auth Auth) http.Handler {
	server := http.NewServeMux()
	server.Handle("/peek", makeHttpHandler(auth))

	return server
}

func makeHttpHandler(auth Auth) *kithttp.Server {
	return kithttp.NewServer(
		JWTMiddleware()(makeAuthEndpoint(auth)),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		kithttp.ServerBefore(jwt.HTTPToContext()),
	)
}
