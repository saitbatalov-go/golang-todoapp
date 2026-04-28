package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/saitbatalov-go/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(ctx context.Context, taskID int, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.tasks
		SET 
			title = $1,
			description = $2,
			completed = $3,
			created_at = $4,
			completed_at = $5,
			author_user_id = $6,
			version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
		taskID,
		task.Version,
	)

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"%v:task with id='%d' not found: %w",
				err,
				taskID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}
	taskDomain := domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CompletedAt,
		taskModel.CreatedAt,
		taskModel.AuthorUserID,
	)
	return taskDomain, nil
}
