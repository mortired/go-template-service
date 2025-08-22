package config

import (
	"fmt"
	"users/internal/infrastructure/config"

	"go.uber.org/zap"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type LoggingConfig struct {
	Level     string
	FxDebug   bool
	FxVerbose bool
	ShowGraph bool
	FxEnabled bool // Enable/disable fx logging completely
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	if err := config.LoadEnv(); err != nil {
		zap.L().Error("failed to load .env file", zap.Error(err))
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: config.GetEnvAsIntOrDefault("SERVER_PORT", 8080),
			Host: config.GetEnvOrDefault("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     config.GetEnvOrDefault("DB_HOST", "localhost"),
			Port:     config.GetEnvAsIntOrDefault("DB_PORT", 5432),
			Username: config.GetEnvOrDefault("DB_USERNAME", "postgres"),
			Password: config.GetEnvOrDefault("DB_PASSWORD", ""),
			Database: config.GetEnvOrDefault("DB_NAME", "users"),
		},
		Logging: LoggingConfig{
			Level:     config.GetEnvOrDefault("LOG_LEVEL", "info"),
			FxDebug:   config.GetEnvAsBoolOrDefault("FX_DEBUG", false),
			FxVerbose: config.GetEnvAsBoolOrDefault("FX_VERBOSE", false),
			ShowGraph: config.GetEnvAsBoolOrDefault("FX_SHOW_GRAPH", false),
			FxEnabled: config.GetEnvAsBoolOrDefault("FX_ENABLED", true), // Enabled by default
		},
	}, nil
}
