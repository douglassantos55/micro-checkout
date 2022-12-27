package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func MakeHTTPServer(svc Service) http.Handler {
	server := httprouter.New()

	placeOrderHandler := kithttp.NewServer(
		reduceStockMiddleware()(makePlaceOrderEndpoint(svc)),
		decodePlaceOrderRequest,
		kithttp.EncodeJSONResponse,
	)
	server.Handler("POST", "/place-order", placeOrderHandler)

	return server
}

func decodePlaceOrderRequest(ctx context.Context, r *http.Request) (any, error) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("could not decode request data")
	}
	return order, nil
}