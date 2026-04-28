package tasks_postgres_repository

import (
	"context"
	"fmt"
)


func (r *TasksRepository) DeleteTask(ctx context.Context, taskID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM todoapp.tasks
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}