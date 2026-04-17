package users_transport_http

import (
	"net/http"

	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromLogger(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_utils.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'id' path params")
		return
	}

	err = h.userService.DeleteUser(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.NoContentResponse()
}
