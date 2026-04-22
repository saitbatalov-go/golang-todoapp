package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/saitbatalov-go/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRespository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.users
		SET 
			full_name = $1,
			phone_number = $2,
			version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING id, version, full_name, phone_number
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		id,
		user.Version,
	)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	); err != nil {
		// if user by id not found and version not match (changed by another user)
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id: %d concurrency accessed: %w",
				err,
				core_errors.ErrConflict,
			)
		}

		return domain.User{}, fmt.Errorf("pathc user: scan row: %w", err)
	}

	userDomain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)

	return userDomain, nil
}
