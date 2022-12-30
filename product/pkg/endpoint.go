package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeListProductsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.ListProducts()
	}
}

func makeGetProductEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetProduct(r.(string))
	}
}

func makeReduceStockEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		in := r.(ReduceStockRequest)
		return nil, svc.ReduceStock(in.ProductID, in.Qty)
	}
}

type ReduceStockRequest struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}
