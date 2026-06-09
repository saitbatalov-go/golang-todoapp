package tasks_service

import (
	"context"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	ports_in "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/out/repository"
)


func (s *TasksService) CreateTask(
	ctx context.Context,
	params ports_in.CreateTaskParams,
) (ports_in.CreateTaskResult, error) {
	task := domain.CreateTask(
		params.Title,
		params.Description,
		params.AuthorUserID,
	)

	if err := task.Validate(); err != nil {
		return ports_in.CreateTaskResult{}, fmt.Errorf("validate task domain: %w", err)
	}

	repoParams := ports_out_repository.NewSaveTaskParams(task)
	repoResult, err := s.tasksRepository.SaveTask(ctx, repoParams)
	if err != nil {
		return ports_in.CreateTaskResult{}, fmt.Errorf("save task in repository: %w", err)
	}

	return ports_in.NewCreateTaskResult(
		repoResult.Task,
	), nil
}
