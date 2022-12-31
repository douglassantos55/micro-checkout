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

type proxyingmw struct {
	next           Service
	getProduct     endpoint.Endpoint
	reduceStock    endpoint.Endpoint
	processPayment endpoint.Endpoint
}

func NewProxyingMiddleware(next Service, getProduct, reduceStock, processPayment endpoint.Endpoint) Service {
	return &proxyingmw{next, getProduct, reduceStock, processPayment}
}

func (m *proxyingmw) GetOrders() ([]*Order, error) {
	return m.next.GetOrders()
}

func (m *proxyingmw) calculateOrderTotal(ctx context.Context, order *Order) error {
	for _, item := range order.Items {
		p, err := m.getProduct(ctx, item.ProductID)
		if err != nil {
			return err
		}
		product := p.(Product)
		item.Subtotal = product.Price * float64(item.Qty)
		order.Total += item.Subtotal
	}
	return nil
}

func (m *proxyingmw) reduceProductsStock(ctx context.Context, order *Order) error {
	for _, item := range order.Items {
		req := struct {
			ProductID string `json:"product_id"`
			Qty       int    `json:"qty"`
		}{item.ProductID, item.Qty}

		if _, err := m.reduceStock(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func (m *proxyingmw) PlaceOrder(ctx context.Context, order *Order) (*Order, error) {
	if err := m.calculateOrderTotal(ctx, order); err != nil {
		return nil, err
	}

	saved, err := m.next.PlaceOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	if err := m.reduceProductsStock(ctx, saved); err != nil {
		return nil, err
	}

	go func(order *Order) {
		invoice, err := m.processPayment(context.Background(), saved)
	}(saved)

	return saved, nil
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
