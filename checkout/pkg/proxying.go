package pkg

import (
	"context"
	"log"

	"github.com/go-kit/kit/endpoint"
)

type subtotalProxy struct {
	next    Service
	product endpoint.Endpoint
}

func MakeSubtotalProxy(next Service, e endpoint.Endpoint) Service {
	return &subtotalProxy{next, e}
}

func (p *subtotalProxy) GetOrders() ([]*Order, error) {
	return p.next.GetOrders()
}

func (p *subtotalProxy) PlaceOrder(ctx context.Context, order *Order) (*Order, error) {
	for _, item := range order.Items {
		p, err := p.product(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		product := p.(Product)
		item.Subtotal = product.Price * float64(item.Qty)
		order.Total += item.Subtotal
	}
	return p.next.PlaceOrder(ctx, order)
}

type stockProxy struct {
	next  Service
	stock endpoint.Endpoint
}

func MakeStockProxy(next Service, e endpoint.Endpoint) Service {
	return &stockProxy{next, e}
}

func (p *stockProxy) GetOrders() ([]*Order, error) {
	return p.next.GetOrders()
}

func (p *stockProxy) PlaceOrder(ctx context.Context, o *Order) (*Order, error) {
	order, err := p.next.PlaceOrder(ctx, o)
	if err != nil {
		return nil, err
	}
	for _, item := range order.Items {
		req := struct {
			ProductID string `json:"product_id"`
			Qty       int    `json:"qty"`
		}{item.ProductID, item.Qty}

		if _, err := p.stock(ctx, req); err != nil {
			return nil, err
		}
	}
	return order, nil
}

type paymentProxy struct {
	next    Service
	payment endpoint.Endpoint
}

func MakePaymentProxy(next Service, e endpoint.Endpoint) Service {
	return &paymentProxy{next, e}
}

func (p *paymentProxy) GetOrders() ([]*Order, error) {
	return p.next.GetOrders()
}

func (p *paymentProxy) PlaceOrder(ctx context.Context, o *Order) (*Order, error) {
	order, err := p.next.PlaceOrder(ctx, o)
	if err != nil {
		return nil, err
	}
	go func(order *Order) {
		invoice, err := p.payment(context.Background(), order)
		log.Printf("invoice: %v, err: %v", invoice, err)
	}(order)
	return order, nil
}
