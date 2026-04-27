package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)


func (r *TasksRepository) GetTasks(ctx context.Context, limit, offset *int) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		if err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		); err != nil {
			return nil, fmt.Errorf("scan row tasks: %w", err)
		}
		taskModels = append(taskModels, taskModel)
	}

	tasksDomains := tasksModelToDomains(taskModels)
	return tasksDomains, nil
}