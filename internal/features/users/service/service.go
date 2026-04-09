package users_service

import (
	"context"

	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)


type UsersService struct {
	UsersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}


func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		UsersRepository: usersRepository,
	}
}