package bootstrap

import (
	"users/internal/infrastructure/logging"
	"users/internal/repository"
	"users/internal/service"
	"users/internal/service/user"
)

func ProvideUserService(repo repository.User, logger *logging.Logger) service.User {
	return user.New(repo, logger)
}
