package tasks_service

import (
	"context"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

func (s *TasksService) GetTasks(ctx context.Context, limit, offset *int) ([]domain.Task, error) {

	if limit != nil && *limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be greater than 0")
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	return tasks, nil
}
