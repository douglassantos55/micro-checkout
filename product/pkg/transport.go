package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/julienschmidt/httprouter"
)

func MakeHTTPHandler(svc Service, logger log.Logger) http.Handler {
	server := httprouter.New()

	listHandler := kithttp.NewServer(
		makeListProductsEndpoint(svc),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
	server.Handler("GET", "/", listHandler)

	getHandler := kithttp.NewServer(
		makeGetProductEndpoint(svc),
		decodeGetProductRequest,
		kithttp.EncodeJSONResponse,
	)
	server.Handler("GET", "/:id", getHandler)

	reduceStockHandler := kithttp.NewServer(
		loggingMiddleware(logger)(makeReduceStockEndpoint(svc)),
		decodeReduceStockRequest,
		kithttp.EncodeJSONResponse,
	)
	server.Handler("POST", "/reduce-stock", reduceStockHandler)

	return server
}

func decodeGetProductRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName("id"), nil
}

func decodeReduceStockRequest(ctx context.Context, r *http.Request) (any, error) {
	var request ReduceStockRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
