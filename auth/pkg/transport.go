package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func NewHttpServer(auth Auth, logger log.Logger) http.Handler {
	server := http.NewServeMux()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	server.Handle("/verify", makeVerifyHandler(opts))
	server.Handle("/login", makeAuthenticateHandler(auth, opts))

	return server
}

func makeAuthenticateHandler(auth Auth, opts []kithttp.ServerOption) *kithttp.Server {
	return kithttp.NewServer(
		makeAuthEndpoint(auth),
		decodeCredentialsRequest,
		kithttp.EncodeJSONResponse,
		opts...,
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

func makeVerifyHandler(opts []kithttp.ServerOption) *kithttp.Server {
	return kithttp.NewServer(
		JWTMiddleware()(makeVerifyEndpoint()),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		append(opts, kithttp.ServerBefore(kitjwt.HTTPToContext()))...,
	)
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json; charset=utf-8")

	switch err {
	case ErrInvalidCredentials:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]any{"err": err.Error()})
}
