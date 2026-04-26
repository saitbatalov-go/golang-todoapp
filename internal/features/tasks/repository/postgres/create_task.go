package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

func (r *TasksRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id
	`
	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Completed, task.CreatedAt, task.CompletedAt, task.AuthorUserID)
	
	var taskModel TaskModel
	if err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	); err != nil {
		return domain.Task{}, fmt.Errorf("scan create task model: %w", err)
	}

	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CompletedAt,
		taskModel.CreatedAt,
		taskModel.AuthorUserID,
	), nil
}
