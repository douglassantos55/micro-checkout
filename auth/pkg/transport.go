package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHttpServer(auth Auth) http.Handler {
	server := http.NewServeMux()

	server.Handle("/verify", makeVerifyHandler())
	server.Handle("/login", makeAuthenticateHandler(auth))

	return server
}

func makeAuthenticateHandler(auth Auth) *kithttp.Server {
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

func makeVerifyHandler() *kithttp.Server {
	return kithttp.NewServer(
		JWTMiddleware()(makeVerifyEndpoint()),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerErrorEncoder(errorEncoder),
	)
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {

	w.Header().Set("content-type", "application/json; charset=utf-8")
	body, err := json.Marshal(map[string]any{"err": err.Error()})

	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
}
