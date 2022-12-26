package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeGetMethodsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetPaymentMethods()
	}
}
