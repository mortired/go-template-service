package bootstrap

import (
	"users/internal/repository"
	"users/internal/service"
	"users/internal/service/user"

	logging "github.com/mortired/appsap-logging"
)

func ProvideUserService(repo repository.User, logger *logging.Logger) service.User {
	return user.New(repo, logger)
}
