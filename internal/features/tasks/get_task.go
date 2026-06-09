package tasks_service

import (
	"context"
	"fmt"

	ports_in "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/out/repository"
)


func (s *TasksService) GetTask(
	ctx context.Context,
	params ports_in.GetTaskParams,
) (ports_in.GetTaskResult, error) {
	repoParams := ports_out_repository.NewGetTaskParams(params.ID)
	repoResult, err := s.tasksRepository.GetTask(ctx, repoParams)
	if err != nil {
		return ports_in.GetTaskResult{}, fmt.Errorf("get task from repository: %w", err)
	}

	return ports_in.NewGetTaskResult(
		repoResult.Task,
	), nil
}
