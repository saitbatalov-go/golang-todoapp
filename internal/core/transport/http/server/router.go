package core_transport_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/middleware"
)

type ApiVersion string

const (
	ApiVersionV1 ApiVersion = "v1"
	ApiVersionV2 ApiVersion = "v2"
	ApiVersionV3 ApiVersion = "v3"
)

type APIVersionRouter struct {
	*http.ServeMux
	ApiVersion ApiVersion
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(apiVersion ApiVersion, middleware ...core_http_middleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		ApiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		handler := route.WithMiddleware() 

		r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

			handler.ServeHTTP(w, r)
		})
	}
}



func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}