package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/golang-jwt/jwt/v4"
)

type ServiceMiddleware func(Greeter) Greeter

type proxyingmw struct {
	next Greeter
	ctx  context.Context
	auth endpoint.Endpoint
}

func (m *proxyingmw) Greet(name string) string {
	res, _ := m.auth(m.ctx, nil)
	return m.next.Greet(res.(string))
}

func ProxyingMiddleware(ctx context.Context, addr string) ServiceMiddleware {
	return func(g Greeter) Greeter {
		return &proxyingmw{g, ctx, getAuthenticatedUser(addr)}
	}
}

// Calls auth service on `addr` and returns its response
func getAuthenticatedUser(addr string) endpoint.Endpoint {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	url, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	url.Path += "/peek"

	return JWTMiddleware()(kithttp.NewClient(
		"GET",
		url,
		kithttp.EncodeJSONRequest,
		getAuthName,
		kithttp.ClientBefore(kitjwt.ContextToHTTP()),
	).Endpoint())
}

func JWTMiddleware() endpoint.Middleware {
	return kitjwt.NewSigner(
		"greeter",
		[]byte("your-256-bit-secret"),
		jwt.SigningMethodHS256,
		jwt.StandardClaims{Subject: "foobar"},
	)
}

// Parses response from auth service
func getAuthName(ctx context.Context, r *http.Response) (any, error) {
	var name string
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		return nil, err
	}
	return name, nil
}
