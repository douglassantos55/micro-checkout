package pkg

import (
	"context"
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
		decodeRequest,
		encodeResponse,
	)
}

func decodeRequest(ctx context.Context, r *http.Request) (any, error) {
	return kithttp.NopRequestDecoder(ctx, r)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, r any) error {
	return kithttp.EncodeJSONResponse(ctx, w, r)
}
