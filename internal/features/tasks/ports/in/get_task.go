package ports_in

import (
	"github.com/google/uuid"
	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)


type GetTaskParams struct {
	ID uuid.UUID
}

func NewGetTaskParams(
	id uuid.UUID,
) GetTaskParams {
	return GetTaskParams{
		ID: id,
	}
}

type GetTaskResult struct {
	Task domain.Task
}

func NewGetTaskResult(
	task domain.Task,
) GetTaskResult {
	return GetTaskResult{
		Task: task,
	}
}
