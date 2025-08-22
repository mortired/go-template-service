package appcore

import (
	"users/internal/infrastructure/logging"
)

// ProvideLogger creates logger instance with optional Elasticsearch integration
func ProvideLogger() (*logging.Logger, error) {
	config := logging.LoadConfigFromEnv()
	return logging.New(config)
}
