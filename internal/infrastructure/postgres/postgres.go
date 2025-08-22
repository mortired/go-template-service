package postgres

import (
	"database/sql"
	"fmt"
	"users/internal/infrastructure/config"
	"users/internal/infrastructure/logging"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Config structure for PostgreSQL configuration
// Allows infrastructure to not depend on application packages
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// DB manages PostgreSQL connection
type DB struct {
	*sql.DB
}

// New creates a new PostgreSQL connection
func New(cfg Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Check connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db}, nil
}

// Close closes database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// ProvidePostgresConfig creates PostgreSQL configuration from environment variables
func ProvidePostgresConfig(logger *logging.Logger) Config {
	// Load environment variables from .env file
	if err := config.LoadEnv(); err != nil {
		// Log error but don't panic, as there might be system variables
		if logger != nil {
			logger.Warn("Failed to load .env file", zap.String("error", err.Error()))
		}
	}

	return Config{
		Host:     config.GetEnvOrDefault("DB_HOST", "localhost"),
		Port:     config.GetEnvAsIntOrDefault("DB_PORT", 5432),
		Username: config.GetEnvOrDefault("DB_USERNAME", "postgres"),
		Password: config.GetEnvOrDefault("DB_PASSWORD", ""),
		Database: config.GetEnvOrDefault("DB_NAME", "users"),
	}
}
