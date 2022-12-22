package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHttpServer(auth Auth) http.Handler {
	server := http.NewServeMux()
	server.Handle("/peek", makeHttpHandler(auth))

	return server
}

func makeHttpHandler(auth Auth) *kithttp.Server {
	return kithttp.NewServer(
		makeAuthEndpoint(auth),
		decodeCredentialsRequest,
		kithttp.EncodeJSONResponse,
	)
}

func decodeCredentialsRequest(ctx context.Context, r *http.Request) (any, error) {
	if r.Method != "POST" {
		return nil, errors.New("invalid request method")
	}
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return nil, err
	}
	return credentials, nil
}
