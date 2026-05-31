package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_request "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swagggertype:"string" example:"Aza Saitbatalov"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swagggertype:"string" example:"+79998887766"`
}

func (r PatchUserRequest) Validate() error {

	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf(
				"PATCH invalid `FullName` can't be value null: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		fullNameLength := len([]rune(*r.FullName.Value))

		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf("invalid `FullName` must be between 3 and 100: %d:%w", fullNameLength, core_errors.ErrInvalidArgument)
		}

	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {

			phoneNumberLength := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLength < 10 || phoneNumberLength > 15 {
				return fmt.Errorf("invalid `PhoneNumber` must be between 10 and 15: %d:%w", phoneNumberLength, core_errors.ErrInvalidArgument)
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("invalid `PhoneNumber` must start with +: %s:%w", *r.PhoneNumber.Value, core_errors.ErrInvalidArgument)
			}
		}

	}
	return nil
}

type PatchUserRespose UserDTOResponse

// PatchUser godoc
// @Summary     Обновление пользователя
// @Description Обновление пользователя существующего в системе по ID
// @Description ### Логика обновления
// @Description 1. **Поле не передано** - `phone_number` игонируется , занчение в бд не меняется
// @Description 2. **Поле передано** - `phone_number` обновляется в бд
// @Description 3. **Поле null**: - `phone_number` очищается в бд (set null)
// @Description 4. Ограничения - `full_name` не может быть пустым null, длина от 3 до 100 символов
// @Tags        Users
// @Produces    json
// @Param       id path int true "ID пользователя"
// @Param       request body PatchUserRequest true "Запрос на обновление пользователя"
// @Success     200 {object} PatchUserRespose
// @Failure     400 {object} core_http_response.ErrorResponse "Плохой запрос"
// @Failure     500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router      /users/{id} [patch]
func (h *UserHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromLogger(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'id' path params")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserRespose(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
