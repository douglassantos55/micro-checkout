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
			return nil, fmt.Errorf("could not convert data into order: %v", r)
		}
		return svc.PlaceOrder(ctx, &order)
	}
}

func MakeGetProductEndpoint() endpoint.Endpoint {
	url, err := url.Parse("http://product-service/")
	if err != nil {
		panic(err)
	}

	return kithttp.NewClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Request, i any) error {
			r.URL.Path += i.(string)
			return nil
		},
		func(ctx context.Context, r *http.Response) (any, error) {
			var product Product
			if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
				return nil, err
			}
			return product, nil
		},
	).Endpoint()
}

func MakeReduceStockEndpoint() endpoint.Endpoint {
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

func MakeProcessPaymentEndpoint() endpoint.Endpoint {
	conn, err := amqp.Dial("amqp://guest:guest@messaging-service/")
	if err != nil {
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	if err := channel.ExchangeDeclare(
		"order-placed",
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		panic(err)
	}

	// publisher is interested in a reply queue, not subscriber queue. The
	// subscriber queue is specified using the publish key
	replyQueue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(err)
	}

	return kitamqp.NewPublisher(
		channel,
		&replyQueue,
		encodeProcessPaymentRequest,
		decodeProcessPaymentResponse,
		kitamqp.PublisherBefore(kitamqp.SetPublishExchange("order-placed")),
	).Endpoint()
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
	var invoice struct {
		ID     string  `json:"id"`
		Total  float64 `json:"total"`
		Status string  `json:"status"`
	}
	if err := json.Unmarshal(d.Body, &invoice); err != nil {
		return nil, err
	}
	return invoice, nil
}
