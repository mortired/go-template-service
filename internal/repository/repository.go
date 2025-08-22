package repository

import "users/internal/model/user"

type User interface {
	ListUsers(filter user.Filter) (user.UsersResponse, error)
	CreateUser(user user.CreateUserRequest) (user.UserResponse, error)
}
