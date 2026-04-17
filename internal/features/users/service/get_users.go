package users_service

import (
	"context"
	"fmt"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {

	if limit != nil && *limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be greater than 0")
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return users, nil
}
