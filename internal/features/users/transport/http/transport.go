package users_transport_http

import (
	"context"
	"net/http"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"
)

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
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
