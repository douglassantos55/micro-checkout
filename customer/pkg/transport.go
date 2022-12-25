package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

	updateCustomerHandler := kithttp.NewServer(
		makeUpdateCustomerEndpoint(svc),
		decodeUpdateRequest,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("PUT", "/:id", updateCustomerHandler)

	deleteCustomerHandler := kithttp.NewServer(
		makeDeleteCustomerEndpoint(svc),
		decodeGetCustomerRequest,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("DELETE", "/:id", deleteCustomerHandler)

	return router
}

func decodeGetCustomerRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName("id"), nil
}

func decodeListCustomersRequest(ctx context.Context, r *http.Request) (any, error) {
	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 0)
	if err != nil {
		page = 1
	}
	filters := Filters{Page: int(page)}
	return filters, nil
}

func decodeCustomerRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Customer
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (any, error) {
	var data Customer
	params := httprouter.ParamsFromContext(r.Context())
	customerId := params.ByName("id")

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return UpdateRequest{customerId, data}, nil
}
