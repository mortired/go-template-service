package bootstrap

import (
	"users/internal/repository"
	userRepo "users/internal/repository/user"

	postgres "github.com/mortired/appsap-postgres"
)

func ProvideUserRepository(pg *postgres.DB) repository.User {
	return userRepo.New(pg)
}
