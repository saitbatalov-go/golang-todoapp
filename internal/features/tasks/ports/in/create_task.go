package ports_in

import (
	"github.com/google/uuid"
	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)


type CreateTaskParams struct {
	Title        string
	Description  *string
	AuthorUserID uuid.UUID
}

func NewCreateTaskParams(
	title string,
	description *string,
	authorUserID uuid.UUID,
) CreateTaskParams {
	return CreateTaskParams{
		Title:        title,
		Description:  description,
		AuthorUserID: authorUserID,
	}
}


type CreateTaskResult struct {
	Task domain.Task
}

func NewCreateTaskResult(
	task domain.Task,
) CreateTaskResult {
	return CreateTaskResult{
		Task: task,
	}
}
