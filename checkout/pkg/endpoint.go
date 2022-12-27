package pkg

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

func makePlaceOrderEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		order, ok := r.(Order)
		if !ok {
			return nil, fmt.Errorf("could not convert data into order: %t", r)
		}
		return svc.PlaceOrder(order)
	}
}

func makeReduceStockEndpoint() endpoint.Endpoint {
	url, err := url.Parse("http://product-service/reduce-stock")
	if err != nil {
		panic(err)
	}

	return kithttp.NewClient(
		"POST",
		url,
		kithttp.EncodeJSONRequest,
		func(ctx context.Context, r *http.Response) (any, error) {
			if r.StatusCode == 500 {
				return nil, fmt.Errorf("could not reduce stock")
			}
			return nil, nil
		},
	).Endpoint()
}
