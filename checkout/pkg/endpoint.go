package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kitamqp "github.com/go-kit/kit/transport/amqp"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/streadway/amqp"
)

func makeGetOrdersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		return svc.GetOrders()
	}
}

func makePlaceOrderEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, r any) (any, error) {
		order, ok := r.(Order)
		if !ok {
			return nil, fmt.Errorf("could not convert data into order: %t", r)
		}
		return svc.PlaceOrder(order)
	}
}

func makeReduceStockEndpoint() endpoint.Endpoint {
	url, err := url.Parse("http://product-service/reduce-stock")
	if err != nil {
		panic(err)
	}

	return kithttp.NewClient(
		"POST",
		url,
		kithttp.EncodeJSONRequest,
		func(ctx context.Context, r *http.Response) (any, error) {
			if r.StatusCode == 500 {
				return nil, fmt.Errorf("could not reduce stock")
			}
			return nil, nil
		},
	).Endpoint()
}

func makeProcessPaymentEndpoint() endpoint.Endpoint {
	conn, err := amqp.Dial("amqp://guest:guest@messaging-service/")
	if err != nil {
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	queue, err := channel.QueueDeclare("orders", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	return kitamqp.NewPublisher(
		channel,
		&queue,
		encodeProcessPaymentRequest,
		decodeProcessPaymentResponse,
		kitamqp.PublisherBefore(setPublishKey(queue.Name)),
		kitamqp.PublisherDeliverer(kitamqp.SendAndForgetDeliverer),
	).Endpoint()
}

func setPublishKey(key string) kitamqp.RequestFunc {
	return func(ctx context.Context, p *amqp.Publishing, d *amqp.Delivery) context.Context {
		return context.WithValue(ctx, kitamqp.ContextKeyPublishKey, key)
	}
}

func encodeProcessPaymentRequest(ctx context.Context, p *amqp.Publishing, r any) error {
	order, ok := r.(*Order)
	if !ok {
		return fmt.Errorf("could not parse r into order: %v", r)
	}
	b, err := json.Marshal(order)
	if err != nil {
		return err
	}
	p.Body = b
	p.ContentType = "application/json"
	return nil
}

func decodeProcessPaymentResponse(ctx context.Context, d *amqp.Delivery) (any, error) {
	return nil, nil
}
