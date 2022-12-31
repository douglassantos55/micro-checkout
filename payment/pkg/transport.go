package pkg

import (
	"context"
	"encoding/json"
	"net/http"

	kitamqp "github.com/go-kit/kit/transport/amqp"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
)

func MakeHTTPServer(svc Service, logger log.Logger) http.Handler {
	router := httprouter.New()

	getMethodsHandler := kithttp.NewServer(
		loggingMiddleware("GetMethods", logger)(makeGetMethodsEndpoint(svc)),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("GET", "/methods", getMethodsHandler)

	getInvoicesHandler := kithttp.NewServer(
		loggingMiddleware("GetInvoices", logger)(makeGetInvoicesEndpoint(svc)),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("GET", "/invoices", getInvoicesHandler)

	return router
}

func MakeAMQPSubscriber(svc Service, logger log.Logger) {
	subscriber := kitamqp.NewSubscriber(
		loggingMiddleware("ProcessPayment", logger)(makeProcessPaymentEndpoint(svc)),
		decodeProcessPaymentRequest(logger),
		kitamqp.EncodeJSONResponse,
	)

	conn, err := amqp.Dial("amqp://guest:guest@messaging-service/")
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

	q, err := channel.QueueDeclare("payment-order-placed", true, false, true, false, nil)
	if err != nil {
		panic(err)
	}

	if err := channel.QueueBind(q.Name, "", "order-placed", false, nil); err != nil {
		panic(err)
	}

	msgs, err := channel.Consume(q.Name, "", true, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			logger.Log("msg received", string(d.Body))
			handler := subscriber.ServeDelivery(channel)
			handler(&d)
			logger.Log("reply delivered")
		}
	}()

	logger.Log("waiting for messages")
	<-forever
}

func decodeProcessPaymentRequest(logger log.Logger) kitamqp.DecodeRequestFunc {
	return func(ctx context.Context, r *amqp.Delivery) (any, error) {
		var order Order
		logger.Log("request", string(r.Body))

		if err := json.Unmarshal(r.Body, &order); err != nil {
			return nil, err
		}

		logger.Log("parsed", order)
		return order, nil
	}
}
