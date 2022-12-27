package pkg

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type loggingmw struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger, svc Service) Service {
	return &loggingmw{svc, logger}
}

func (m *loggingmw) PlaceOrder(order Order) (*Order, error) {
	defer func(timestamp time.Time) {
		m.logger.Log(
			"method", "PlaceOrder",
			"order", order,
			"order", order,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return m.next.PlaceOrder(order)
}

func reduceStockMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			order := res.(*Order)
			reduceStock := makeReduceStockEndpoint()
			for _, item := range order.Items {
				req := struct {
					ProductID string `json:"product_id"`
					Qty       int    `json:"qty"`
				}{item.ProductID, item.Qty}

				if _, err := reduceStock(ctx, req); err != nil {
					return nil, err
				}
			}

			return order, nil
		}
	}
}
