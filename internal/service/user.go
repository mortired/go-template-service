package service

import "users/internal/model/user"

type User interface {
	ListUsers(filter user.Filter) (user.UsersResponse, error)
	CreateUser(req user.CreateUserRequest) (user.UserResponse, error)
}
