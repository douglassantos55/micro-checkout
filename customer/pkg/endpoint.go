package pkg

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

var (
	ErrInvalidCustomerID = makeError(400, "invalid customer ID")
	ErrInvalidFilters    = makeError(400, "invalid filters")
	ErrInvalidData       = makeError(400, "invalid data")
)

func makeGetCustomerEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		customerId, ok := r.(string)
		if !ok {
			return nil, ErrInvalidCustomerID
		}
		return svc.GetCustomer(customerId)
	}
}

func makeListCustomersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		filters, ok := r.(Filters)
		if !ok {
			return nil, ErrInvalidFilters
		}
		return svc.ListCustomers(filters)
	}
}

func makeCreateCustomerEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		data, ok := r.(Customer)
		if !ok {
			return nil, ErrInvalidData
		}
		return svc.CreateCustomer(data)
	}
}

func makeUpdateCustomerEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		data, ok := r.(UpdateRequest)
		if !ok {
			return nil, ErrInvalidData
		}
		if data.ID == "" {
			return nil, ErrInvalidCustomerID
		}
		return svc.UpdateCustomer(data.ID, data.Customer)
	}
}

func makeDeleteCustomerEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		id, ok := r.(string)
		if !ok {
			return nil, ErrInvalidCustomerID
		}
		return nil, svc.DeleteCustomer(id)
	}
}

type UpdateRequest struct {
	ID       string
	Customer Customer
}
