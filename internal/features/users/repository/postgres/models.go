package user_postgres_repository

import "github.com/saitbatalov-go/golang-todoapp/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}


func userDomainsFromUserModels(userModels []UserModel) []domain.User {
	var userDomains []domain.User
	for _, userModel := range userModels {
		userDomains = append(userDomains, domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber))
	}
	return userDomains
}