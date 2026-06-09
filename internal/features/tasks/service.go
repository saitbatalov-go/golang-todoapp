package tasks_service

import (
	repository "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/out/repository"
)


type TasksService struct {
	tasksRepository repository.TasksRepository
}

func NewTasksService(
	tasksRepository repository.TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
