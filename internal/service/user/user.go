package user

import (
	"context"
	"users/internal/infrastructure/logging"
	"users/internal/model/user"
	"users/internal/repository"
	"users/internal/service"

	"go.uber.org/zap"
)

type Service struct {
	repo   repository.User
	logger *logging.Logger
}

func New(repo repository.User, logger *logging.Logger) service.User {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) ListUsers(filter user.Filter) (user.UsersResponse, error) {
	ctx := context.Background()
	log := s.logger.WithContext(ctx)

	log.Info("ListUsers service called",
		zap.Any("filter", filter),
	)

	users, err := s.repo.ListUsers(filter)
	if err != nil {
		log.Error("Failed to list users from repository",
			zap.Error(err),
			zap.Any("filter", filter),
		)
		return user.UsersResponse{}, err
	}

	log.Info("Users listed successfully",
		zap.Int("count", len(users)),
	)

	return users, nil
}

func (s *Service) CreateUser(req user.CreateUserRequest) (user.UserResponse, error) {
	ctx := context.Background()
	log := s.logger.WithContext(ctx)

	log.Info("CreateUser service called",
		zap.String("email", req.Email),
		zap.String("name", req.Name),
	)

	createdUser, err := s.repo.CreateUser(req)
	if err != nil {
		log.Error("Failed to create user in repository",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("name", req.Name),
		)
		return user.UserResponse{}, err
	}

	log.Info("User created successfully",
		zap.Int("user_id", int(createdUser.ID)),
		zap.String("email", createdUser.Email),
	)

	return createdUser, nil
}
