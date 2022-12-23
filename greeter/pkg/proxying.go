package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Calls auth service on `addr` to verify token
func verifyAuthTokenProxy(addr string) endpoint.Endpoint {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	url, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	url.Path += "/verify"

	return kithttp.NewClient(
		"GET",
		url,
		kithttp.EncodeJSONRequest,
		decodeVerifyResponse,
		// Add token from context into http header
		// before sending the request
		kithttp.ClientBefore(jwt.ContextToHTTP()),
	).Endpoint()
}

// Get username from response
func decodeVerifyResponse(ctx context.Context, r *http.Response) (any, error) {
	if r.StatusCode == http.StatusBadRequest {
		return nil, errors.New("could not verify request")
	}
	var name string
	err := json.NewDecoder(r.Body).Decode(&name)
	return name, err
}
