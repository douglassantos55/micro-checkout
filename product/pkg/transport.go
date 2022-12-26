package pkg

import (
	"context"
	"net/http"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func MakeHTTPHandler(svc Service) http.Handler {
	server := httprouter.New()

	listHandler := kithttp.NewServer(
		makeListProductsEndpoint(svc),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
	server.Handler("GET", "/", listHandler)

	return server
}

func MakeGRPCServer(svc Service) kitgrpc.Handler {
	return kitgrpc.NewServer(
		makeReduceStockEndpoint(svc),
		decodeReduceStockRequest,
		nil,
	)
}

func decodeReduceStockRequest(ctx context.Context, r any) (any, error) {
	return r.(ReduceStockRequest), nil
}
