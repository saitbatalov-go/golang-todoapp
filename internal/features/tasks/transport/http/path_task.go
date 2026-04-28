package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_request "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
}

func (r *PatchTaskRequest) Validate() error {
	
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf(
				"PATCH invalid `Title` can't be value null: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		titleLength := len([]rune(*r.Title.Value))

		if titleLength < 3 || titleLength > 100 {
			return fmt.Errorf("invalid `Title` must be between 3 and 100: %d:%w", titleLength, core_errors.ErrInvalidArgument)
		}

	}

	if r.Description.Set {
		if r.Description.Value != nil {

			descriptionLength := len([]rune(*r.Description.Value))
			if descriptionLength < 3 || descriptionLength > 1000 {
				return fmt.Errorf("invalid `Description` must be between 3 and 1000: %d:%w", descriptionLength, core_errors.ErrInvalidArgument)
			}

		}
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromLogger(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	id, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'id' path params")
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, id, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
	)
}