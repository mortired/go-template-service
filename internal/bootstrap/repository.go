package bootstrap

import (
	"users/internal/infrastructure/postgres"
	"users/internal/repository"
	userRepo "users/internal/repository/user"
)

func ProvideUserRepository(pg *postgres.DB) repository.User {
	return userRepo.New(pg)
}
