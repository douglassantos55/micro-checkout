package pkg

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

func authMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		verify := makeAuthEndpoint()
		if _, err := verify(ctx, r); err != nil {
			return nil, err
		}
		return next(ctx, r)
	}
}

type loggingmw struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger, svc Service) Service {
	return &loggingmw{svc, logger}
}

func (m *loggingmw) GetOrders() (orders []*Order, err error) {
	defer func(timestamp time.Time) {
		m.logger.Log(
			"method", "GetOrders",
			"orders", orders,
			"err", err,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return m.next.GetOrders()
}

func (m *loggingmw) PlaceOrder(ctx context.Context, order *Order) (*Order, error) {
	defer func(timestamp time.Time) {
		m.logger.Log(
			"method", "PlaceOrder",
			"order", order,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return m.next.PlaceOrder(ctx, order)
}

func Logging(name string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (res any, err error) {
			defer func(start time.Time) {
				logger.Log(
					"method", name,
					"request", r,
					"res", res,
					"err", err,
					"took", time.Since(start),
				)
			}(time.Now())

			return next(ctx, r)
		}
	}
}
