package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
)

func MakeHTTPServer(svc Service) http.Handler {
	router := httprouter.New()

	getCustomerHandler := kithttp.NewServer(
		makeGetCustomerEndpoint(svc),
		decodeGetCustomerRequest,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("GET", "/:id", getCustomerHandler)

	listCustomersHandler := kithttp.NewServer(
		makeListCustomersEndpoint(svc),
		decodeListCustomersRequest,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("GET", "/", listCustomersHandler)

	createCustomerHandler := kithttp.NewServer(
		makeCreateCustomerEndpoint(svc),
		decodeCustomerRequest,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("POST", "/", createCustomerHandler)

	return router
}

func decodeGetCustomerRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName("id"), nil
}

func decodeListCustomersRequest(ctx context.Context, r *http.Request) (any, error) {
	return r.URL.Query(), nil
}

func decodeCustomerRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Customer
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}