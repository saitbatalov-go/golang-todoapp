package tasks_transport_http

import (
	"net/http"
	"time"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_request "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserID int `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse struct {
	Task TaskDTOResponse
	Approximate time.Duration
}

// CreateTask godoc
// @Summary     Создание задачи
// @Description Создание задачи
// @Tags        Tasks
// @Produces    json
// @Param       request body CreateTaskRequest true "Запрос на создание задачи"
// @Success     201 {object} CreateTaskResponse
// @Failure     400 {object} core_http_response.ErrorResponse "Плохой запрос"
// @Failure     500 {object} core_http_response.ErrorResponse "Внутренняя ошибка сервера"
// @Router      /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromLogger(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	task, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response:= taskDTOFromDomain(task)

	responseHandler.JSONResponse(response, http.StatusCreated)


}