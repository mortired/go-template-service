package appcore

import (
	"users/internal/infrastructure/appcore/config"
	"users/internal/infrastructure/logging"

	"go.uber.org/zap"
)

func ProvideConfig() *config.Config {
	// Создаем временный logger для логирования ошибок конфигурации
	tempLogger, err := logging.New(logging.Config{
		Level:       "info",
		Environment: "development",
		Service:     "users",
		Version:     "1.0.0",
	})
	if err != nil {
		// Если не удалось создать logger, используем стандартный zap
		zap.L().Error("failed to create logger for config loading", zap.Error(err))
	} else {
		defer tempLogger.Sync()
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		if tempLogger != nil {
			tempLogger.Error("failed to load configuration", zap.Error(err))
		} else {
			zap.L().Error("failed to load configuration", zap.Error(err))
		}
		// В случае ошибки загрузки конфигурации, возвращаем nil
		// fx будет обрабатывать это как ошибку инициализации
		return nil
	}
	return cfg
}
