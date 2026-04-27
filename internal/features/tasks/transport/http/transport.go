package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_transport_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, limit, offset *int) ([]domain.Task, error)
}

func NewTasksHTTPHandler(tasksService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_transport_server.Route {
	return []core_transport_server.Route{{
		Method:  http.MethodPost,
		Path:    "/tasks",
		Handler: h.CreateTask,
	}}
}
