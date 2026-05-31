package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_request "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers godoc
// @Summary     Получение пользователей
// @Description Получение пользователей существующих в системе
// @Tags        Users
// @Produces    json
// @Param       limit query int false "Limit"
// @Param       offset query int false "Offset"
// @Success     200 {object} GetUsersResponse
// @Failure     400 {object} core_http_response.ErrorResponse "Плохой запрос"
// @Failure     500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router      /users [get]
func (h *UserHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromLogger(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit offset query params")
		return
	}

	usersDomains, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(usersDomains))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {

	const (
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'lmimit' query params:%w", err)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query params:%w", err)
	}
	return limit, offset, nil
}
