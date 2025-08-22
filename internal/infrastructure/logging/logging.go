package logging

import (
	"context"
	"fmt"
	"users/internal/infrastructure/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds logger configuration
type Config struct {
	Level       string
	Environment string
	Service     string
	Version     string
}

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// New creates a new logger instance
func New(config Config) (*Logger, error) {
	// Load Elasticsearch configuration (cached, loaded only once)
	esConfig, err := GetElasticsearchConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load Elasticsearch configuration: %w", err)
	}

	// Use Elasticsearch-enabled logger if configured
	if esConfig != nil && esConfig.Enabled {
		return NewLoggerWithElasticsearch(config, *esConfig)
	}

	// Otherwise use standard logger
	return NewStandardLogger(config)
}

// NewStandardLogger creates a standard logger without Elasticsearch
func NewStandardLogger(config Config) (*Logger, error) {
	var zapConfig zap.Config

	switch config.Environment {
	case "production":
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "@timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.TimeKey = "@timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Set log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Create logger
	zapLogger, err := zapConfig.Build(
		zap.Fields(
			zap.String("service", config.Service),
			zap.String("version", config.Version),
			zap.String("environment", config.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger}, nil
}

// WithContext creates a new logger with context values
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	fields := []zap.Field{}

	// Add trace_id if present
	if traceID, ok := ctx.Value("trace_id").(string); ok && traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// Add request_id if present
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	// Add user_id if present
	if userID, ok := ctx.Value("user_id").(string); ok && userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	if len(fields) == 0 {
		return l
	}

	return &Logger{l.Logger.With(fields...)}
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return &Logger{l.Logger.With(zapFields...)}
}

// LoadConfigFromEnv loads logger configuration from environment variables
func LoadConfigFromEnv() Config {
	return Config{
		Level:       config.GetEnvOrDefault("LOG_LEVEL", "INFO"),
		Environment: config.GetEnvOrDefault("ENVIRONMENT", "development"),
		Service:     config.GetEnvOrDefault("SERVICE_NAME", "users"),
		Version:     config.GetEnvOrDefault("SERVICE_VERSION", "1.0.0"),
	}
}
