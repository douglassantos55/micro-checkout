package pkg

import (
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/julienschmidt/httprouter"
)

func MakeHTTPServer(svc Service, logger log.Logger) http.Handler {
	router := httprouter.New()

	getMethodsHandler := kithttp.NewServer(
		loggingMiddleware(logger)(makeGetMethodsEndpoint(svc)),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
	router.Handler("GET", "/", getMethodsHandler)

	return router
}
