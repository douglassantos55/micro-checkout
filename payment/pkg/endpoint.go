package pkg

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

func makeGetMethodsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetPaymentMethods()
	}
}

func makeGetInvoicesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetInvoices()
	}
}

func makeProcessPaymentEndpoint(svc Service) endpoint.Endpoint {
	return func(cxt context.Context, r any) (any, error) {
		order, ok := r.(Order)
		if !ok {
			return nil, fmt.Errorf("could not parse request data")
		}
		return svc.ProcessPayment(order)
	}
}
