package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

func (r *UsersRespository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel
	if err := row.Scan(&userModel.ID, &userModel.Version, &userModel.FullName, &userModel.PhoneNumber); err != nil {
		return domain.User{}, fmt.Errorf("scan row: %w", err)
	}

	userDomain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)

	return userDomain, nil
}
