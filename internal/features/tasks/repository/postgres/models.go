package tasks_postgres_repository

import (
	"time"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)


type TaskModel struct {
	ID          int
	Version     int
	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
	AuthorUserID int
}

func tasksModelToDomains(tasks []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		domains[i] = domain.NewTask(
			task.ID,
			task.Version,
			task.Title,
			task.Description,
			task.Completed,
			task.CompletedAt,
			task.CreatedAt,
			task.AuthorUserID,
		)
	}
	return domains
}