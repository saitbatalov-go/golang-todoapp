package tasks_service

import (
	"context"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
}

func NewTasksService(tasksRepository TasksRepository) TasksService {
	return TasksService{
		tasksRepository: tasksRepository,
	}
}
