package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	kitamqp "github.com/go-kit/kit/transport/amqp"
	"github.com/go-kit/log"
	"github.com/streadway/amqp"
)

func NewMessageBroker(logger log.Logger, user, pass, url string) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", user, pass, url))
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	if err := channel.ExchangeDeclare(
		"order-placed",
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		panic(err)
	}

	q, err := channel.QueueDeclare("customer-order-placed", true, false, true, false, nil)
	if err != nil {
		panic(err)
	}

	if err := channel.QueueBind(q.Name, "", "order-placed", false, nil); err != nil {
		panic(err)
	}

	msgs, err := channel.Consume(q.Name, "", false, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	subscriber := kitamqp.NewSubscriber(
		func(ctx context.Context, r any) (any, error) {
			logger.Log("param", r)
			return nil, nil
		},
		func(ctx context.Context, d *amqp.Delivery) (any, error) {
			var order struct{}
			if err := json.Unmarshal(d.Body, &order); err != nil {
				return nil, err
			}
			return order, nil
		},
		kitamqp.EncodeNopResponse,
	)

	var forever chan struct{}
	go func() {
		for d := range msgs {
			logger.Log("msg received", string(d.Body))
			handler := subscriber.ServeDelivery(channel)
			handler(&d)
		}
	}()

	logger.Log("waiting for messages")
	<-forever
}
