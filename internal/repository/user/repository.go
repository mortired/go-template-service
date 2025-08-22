package user

import (
	"database/sql"
	"fmt"
	"strings"
	"users/internal/infrastructure/postgres"
	"users/internal/model/user"
	"users/internal/repository"
)

type Repository struct {
	db *sql.DB
}

func New(pg *postgres.DB) repository.User {
	return &Repository{
		db: pg.DB,
	}
}

func (r *Repository) ListUsers(filter user.Filter) (user.UsersResponse, error) {
	query := "SELECT id, name, email FROM users.users"
	var args []interface{}
	var conditions []string
	argIndex := 1

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *filter.ID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY id ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := make(user.UsersResponse, 0)
	for rows.Next() {
		var u user.UserResponse
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	return users, nil
}

func (r *Repository) CreateUser(req user.CreateUserRequest) (user.UserResponse, error) {
	query := "INSERT INTO users.users (name, email) VALUES ($1, $2) RETURNING id, name, email"

	var user user.UserResponse
	err := r.db.QueryRow(query, req.Name, req.Email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return user, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
