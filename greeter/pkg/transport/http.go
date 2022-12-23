package transport

import (
	"net/http"

	"example.com/microservices/greeter/pkg"
	"github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHttpServer(svc pkg.Greeter, authService string) http.Handler {
	server := http.NewServeMux()

	server.Handle("/greet", makeHttpHandler(svc, authService))
	server.Handle("/greetings", makeHttpHandler(svc, authService))

	return server
}

func makeHttpHandler(greeter pkg.Greeter, authService string) *kithttp.Server {
	endpoint := pkg.MakeGreetEndpoint(greeter)

	return kithttp.NewServer(
		// Checks token with auth middleware and
		// forwards the auth service's response to
		// greet endpoint
		pkg.AuthenticationMiddleware(authService)(endpoint),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		// Adds the http auth header into context
		// before processing request
		kithttp.ServerBefore(jwt.HTTPToContext()),
	)
}
