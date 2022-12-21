package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

type ServiceMiddleware func(Greeter) Greeter

type proxyingmw struct {
	next Greeter
	ctx  context.Context
	auth endpoint.Endpoint
}

func (m *proxyingmw) Greet(name string) map[string]string {
	res, _ := m.auth(m.ctx, nil)
	return m.next.Greet(res.(string))
}

func ProxyingMiddleware(ctx context.Context, addr string) ServiceMiddleware {
	return func(g Greeter) Greeter {
		return &proxyingmw{g, ctx, makeProxy(addr)}
	}
}

func makeProxy(addr string) endpoint.Endpoint {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	url, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	url.Path += "/peek"

	return kithttp.NewClient(
		"GET",
		url,
		kithttp.EncodeJSONRequest,
		decodeResponse,
	).Endpoint()
}

func decodeResponse(ctx context.Context, r *http.Response) (any, error) {
	var name string
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		return nil, err
	}
	return name, nil
}
