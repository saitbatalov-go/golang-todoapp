package ports_in

import (
	"context"
)


type TasksService interface {
	CreateTask(
		ctx context.Context,
		in CreateTaskParams,
	) (CreateTaskResult, error)

	GetTasks(
		ctx context.Context,
		in GetTasksParams,
	) (GetTasksResult, error)

	GetTask(
		ctx context.Context,
		in GetTaskParams,
	) (GetTaskResult, error)

	DeleteTask(
		ctx context.Context,
		in DeleteTaskParams,
	) (DeleteTaskResult, error)

	PatchTask(
		ctx context.Context,
		in PatchTaskParams,
	) (PatchTaskResult, error)
}
