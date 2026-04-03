package users_service

import (
	"context"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

func (s *UsersService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	
	return s.UsersRepository.CreateUser(ctx, user)
}