package users_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"
)

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		fullName string,
		phoneNumber *string,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		patch domain.UserPatch,
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
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
			// Middleware: []core_http_middleware.Middleware{
			// 	core_http_middleware.Dummy("get users Middleware"),
			// },
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
