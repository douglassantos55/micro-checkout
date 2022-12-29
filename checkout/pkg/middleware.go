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

func (m *loggingmw) PlaceOrder(ctx context.Context, order Order) (*Order, error) {
	defer func(timestamp time.Time) {
		m.logger.Log(
			"method", "PlaceOrder",
			"order", order,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return m.next.PlaceOrder(ctx, order)
}

func reduceStockMiddleware(logger log.Logger) endpoint.Middleware {
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

				logger.Log("reducing stock", req)

				if _, err := reduceStock(ctx, req); err != nil {
					return nil, err
				}

				logger.Log("stock reduced", req)
			}

			logger.Log("stock processed, proceeding", order)
			return order, nil
		}
	}
}

func processPaymentMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (invoice any, err error) {
			defer logger.Log("msg", "payment processed", "invoice", invoice, "err", err)

			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			order := res.(*Order)
			logger.Log("processing payment", order)

			processPayment := makeProcessPaymentEndpoint()
			return processPayment(ctx, order)
		}
	}
}
