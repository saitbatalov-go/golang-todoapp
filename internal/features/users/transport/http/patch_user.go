package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_request "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
)

type PatchUserRequest struct {
	FullName    string `json:"full_name" `
	PhoneNumber string `json:"phone_number"`
}

func (h *UserHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromLogger(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	log.Debug(
		fmt.Sprintf(
			"PatchUserREequest fields:\nFullName: %s'\nPhoneNumber:'%s'",
			request.FullName,
			request.PhoneNumber,
		),
	)

	rw.WriteHeader(http.StatusOK)
	
}
