package web_transport_http

import core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHTTPHandler(webService WebService) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

func (h *WebHTTPHandler) Routes() []core_transport_server.Route {
	return []core_transport_server.Route{
		{
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
