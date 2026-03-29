package users_transport_http

import (
	"net/http"

	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/server"
)

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
}

func NewUserHTTPHandler(userService UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: userService,
	}
}

func (h *UserHTTPHandler) Routes() []core_transport_server.Route {

	return []core_transport_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
