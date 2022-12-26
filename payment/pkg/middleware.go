package pkg

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Service middlewares
type loggingmw struct {
	logger log.Logger
	next   Service
}

func LoggingMiddleware(logger log.Logger, svc Service) Service {
	return &loggingmw{logger, svc}
}

func (m *loggingmw) GetPaymentMethods() (methods []*PaymentMethod, err error) {
	defer func(timestamp time.Time) {
		m.logger.Log(
			"method", "GetPaymentMethods",
			"methods", methods,
			"error", err,
			"took", time.Since(timestamp),
		)
	}(time.Now())

	return m.next.GetPaymentMethods()
}

// Endpoints middlewares
func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req any) (res any, err error) {
			defer logger.Log("msg", "endpoint executed", "req", req, "res", res, "err", err)
			return next(ctx, req)
		}
	}
}
