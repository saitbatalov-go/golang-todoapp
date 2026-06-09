package adapters_in_transport_http

import (
	"net/http"

	core_http_server "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/server"
	ports_in "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/in"
)


type TasksHTTPHandler struct {
	tasksService ports_in.TasksService
}

func NewTasksHTTPHandler(
	tasksService ports_in.TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PatchTask,
		},
	}
}
