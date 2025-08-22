package migration

import (
	"fmt"
	"users/internal/infrastructure/config"
)

// LoadConfigFromEnv loads migration configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	// Load environment variables from .env file
	if err := config.LoadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Get database connection parameters
	host := config.GetEnvOrDefault("DB_HOST", "localhost")
	port := config.GetEnvAsIntOrDefault("DB_PORT", 5432)
	username := config.GetEnvOrDefault("DB_USERNAME", "postgres")
	password := config.GetEnvOrDefault("DB_PASSWORD", "")
	database := config.GetEnvOrDefault("DB_NAME", "users")
	migrationsPath := config.GetEnvOrDefault("MIGRATIONS_PATH", "./migrations")

	// Form database connection URL
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		username, password, host, port, database)

	return &Config{
		DatabaseURL:    databaseURL,
		MigrationsPath: migrationsPath,
	}, nil
}

// NewWithEnv creates a new migrator with automatic configuration loading from environment variables
func NewWithEnv() (*Migrator, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return New(cfg)
}
