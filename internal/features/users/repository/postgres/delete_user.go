package user_postgres_repository

import (
	"context"

	"fmt"

	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
)


func (s *UsersRespository) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM todoapp.users
		WHERE id = $1
	`

	cmdTag, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user not found with id: %d: %w",
			id,
			core_errors.ErrNotFound,
		)
	} 

	return nil
}