package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, taskID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM todoapp.tasks
		WHERE id = $1
	`
	cmdTag, err := r.pool.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"task not found with id: %d: %w",
			taskID,
			core_errors.ErrNotFound,
		)
	}
	return nil
}
