package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	kitamqp "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

type MessageBroker interface {
	Broadcast(context.Context, string, any) error
}

type rabbitmq struct {
	conn *amqp.Connection
}

func NewRabbitMQBroker(user, pass, url string) (MessageBroker, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", user, pass, url))
	if err != nil {
		return nil, err
	}
	return &rabbitmq{conn}, nil
}

func (r *rabbitmq) Broadcast(ctx context.Context, event string, data any) error {
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}

	if err := channel.ExchangeDeclare(
		event,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	q, err := channel.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return err
	}

	endpoint := kitamqp.NewPublisher(
		channel,
		&q,
		r.encodeJSONRequest,
		r.nopResponseDecoder,
		kitamqp.PublisherDeliverer(kitamqp.SendAndForgetDeliverer),
		kitamqp.PublisherBefore(kitamqp.SetPublishExchange(event)),
	).Endpoint()

	_, err = endpoint(ctx, data)
	return err
}

func (r *rabbitmq) encodeJSONRequest(ctx context.Context, p *amqp.Publishing, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	p.Body = b
	p.ContentType = "application/json"

	return nil
}

func (r *rabbitmq) nopResponseDecoder(ctx context.Context, d *amqp.Delivery) (any, error) {
	return nil, nil
}
