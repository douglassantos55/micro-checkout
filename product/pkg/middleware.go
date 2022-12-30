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

func (l *loggingmw) GetProduct(id string) (prod *Product, err error) {
	defer func(timestamp time.Time) {
		l.logger.Log(
			"method", "GetProduct",
			"id", id,
			"product", prod,
			"error", err,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return l.next.GetProduct(id)
}

func (l *loggingmw) ListProducts() (products []*Product, err error) {
	defer func(timestamp time.Time) {
		l.logger.Log(
			"method", "ListProducts",
			"products", products,
			"error", err,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return l.next.ListProducts()
}

func (l *loggingmw) ReduceStock(id string, qty int) (err error) {
	defer func(timestamp time.Time) {
		l.logger.Log(
			"method", "ReduceStock",
			"product_id", id,
			"qty", qty,
			"error", err,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return l.next.ReduceStock(id, qty)
}

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (res any, err error) {
			defer logger.Log("msg", "endpoint called", "req", r, "res", res, "err", err)
			return next(ctx, r)
		}
	}
}
