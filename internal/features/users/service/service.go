package users_service


type UsersService struct {
	UsersRepository UsersRepository
}

type UsersRepository interface {
	
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		UsersRepository: usersRepository,
	}
}