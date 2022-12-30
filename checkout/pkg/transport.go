package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/julienschmidt/httprouter"
)

func MakeHTTPServer(svc Service, logger log.Logger) http.Handler {
	server := httprouter.New()

	placeOrderHandler := kithttp.NewServer(
		makePlaceOrderEndpoint(svc),
		decodePlaceOrderRequest,
		kithttp.EncodeJSONResponse,
	)
	server.Handler("POST", "/orders", placeOrderHandler)

	getOrdersHandler := kithttp.NewServer(
		makeGetOrdersEndpoint(svc),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
	server.Handler("GET", "/orders", getOrdersHandler)

	return server
}

func decodePlaceOrderRequest(ctx context.Context, r *http.Request) (any, error) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("could not decode request data")
	}
	return order, nil
}
