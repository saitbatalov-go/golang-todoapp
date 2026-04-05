package user_postgres_repository

import core_postgres_pool "github.com/saitbatalov-go/golang-todoapp/internal/core/repository/postgres/pool"


type UsersRespository struct {
	pool core_postgres_pool.Pool

}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRespository {
	return &UsersRespository{
		pool: pool,
	}
}


