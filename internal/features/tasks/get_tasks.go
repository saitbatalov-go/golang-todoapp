package tasks_service

import (
	"context"
	"fmt"

	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"

	ports_in "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/out/repository"
)


func (s *TasksService) GetTasks(
	ctx context.Context,
	params ports_in.GetTasksParams,
) (ports_in.GetTasksResult, error) {
	if params.Limit != nil && *params.Limit < 0 {
		return ports_in.GetTasksResult{}, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if params.Offset != nil && *params.Offset < 0 {
		return ports_in.GetTasksResult{}, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	repoParams := ports_out_repository.NewGetTasksParams(
		params.UserID,
		params.Limit,
		params.Offset,
	)
	repoResult, err := s.tasksRepository.GetTasks(ctx, repoParams)
	if err != nil {
		return ports_in.GetTasksResult{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	return ports_in.NewGetTasksResult(
		repoResult.Tasks,
	), nil
}
