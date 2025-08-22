package bootstrap

import (
	"users/internal/controller/user"
	"users/internal/service"
)

func ProvideUserController(service service.User) *user.Controller {
	return user.New(service)
}
