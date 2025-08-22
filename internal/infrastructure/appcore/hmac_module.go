package appcore

import (
	"encoding/json"
	"fmt"
	"users/internal/infrastructure/config"
	"users/internal/infrastructure/logging"
	"users/internal/infrastructure/middleware"
	"users/internal/infrastructure/response"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// HMACConfig structure for HMAC configuration
// Allows infrastructure to not depend on application packages
type HMACConfig struct {
	ClientSecrets []middleware.HMACClientSecret
	RouteRights   middleware.HMACRouteRights
	Algorithm     string
	MaxAge        int
	Required      bool
}

// ProvideHMACConfig creates HMAC configuration from environment variables
func ProvideHMACConfig(logger *logging.Logger) (HMACConfig, error) {
	// Load environment variables from .env file
	if err := config.LoadEnv(); err != nil {
		// Log error but don't panic, as there might be system variables
		if logger != nil {
			logger.Warn("Failed to load .env file", zap.String("error", err.Error()))
		}
	}

	// Load client secrets
	clientSecretsStr := config.GetEnvOrDefault("HMAC_CLIENT_SECRETS", "[]")
	var clientSecrets []middleware.HMACClientSecret
	if err := json.Unmarshal([]byte(clientSecretsStr), &clientSecrets); err != nil {
		problem := response.InternalError(fmt.Sprintf("Error parsing HMAC_CLIENT_SECRETS: %v", err), "hmac_config")
		if logger != nil {
			logger.Error("HMAC Configuration Error", zap.String("problem", fmt.Sprintf("%+v", problem)))
		}
		return HMACConfig{}, fmt.Errorf("failed to parse HMAC_CLIENT_SECRETS: %v", err)
	}

	// Load route access rights
	routeRightsStr := config.GetEnvOrDefault("HMAC_ROUTE_RIGHTS", "{}")
	var routeRights middleware.HMACRouteRights
	if err := json.Unmarshal([]byte(routeRightsStr), &routeRights); err != nil {
		problem := response.InternalError(fmt.Sprintf("Error parsing HMAC_ROUTE_RIGHTS: %v", err), "hmac_config")
		if logger != nil {
			logger.Error("HMAC Configuration Error", zap.String("problem", fmt.Sprintf("%+v", problem)))
		}
		return HMACConfig{}, fmt.Errorf("failed to parse HMAC_ROUTE_RIGHTS: %v", err)
	}

	return HMACConfig{
		ClientSecrets: clientSecrets,
		RouteRights:   routeRights,
		Algorithm:     config.GetEnvOrDefault("HMAC_ALGORITHM", "sha256"),
		MaxAge:        config.GetEnvAsIntOrDefault("HMAC_MAX_AGE", 300), // 5 minutes by default
		Required:      config.GetEnvAsBoolOrDefault("HMAC_REQUIRED", true),
	}, nil
}

// ProvideHMACMiddleware creates HMAC middleware with configuration
func ProvideHMACMiddleware(cfg HMACConfig) middleware.HMACConfig {
	return &simpleHMACConfig{
		clientSecrets: cfg.ClientSecrets,
		routeRights:   cfg.RouteRights,
		algorithm:     cfg.Algorithm,
		maxAge:        cfg.MaxAge,
		required:      cfg.Required,
	}
}

// simpleHMACConfig simple implementation of HMACConfig
type simpleHMACConfig struct {
	clientSecrets []middleware.HMACClientSecret
	routeRights   middleware.HMACRouteRights
	algorithm     string
	maxAge        int
	required      bool
}

func (c *simpleHMACConfig) GetClientSecrets() []middleware.HMACClientSecret {
	return c.clientSecrets
}

func (c *simpleHMACConfig) GetRouteRights() middleware.HMACRouteRights {
	return c.routeRights
}

func (c *simpleHMACConfig) GetAlgorithm() string {
	return c.algorithm
}

func (c *simpleHMACConfig) GetMaxAge() int {
	return c.maxAge
}

func (c *simpleHMACConfig) IsRequired() bool {
	return c.required
}

// SetupHMACMiddleware configures HMAC middleware
func SetupHMACMiddleware(e *echo.Echo, config middleware.HMACConfig, logger *logging.Logger) {
	if !config.IsRequired() {
		if logger != nil {
			logger.Info("HMAC middleware is disabled")
		}
		return
	}

	hmacMiddleware := middleware.HMACAuthMiddleware(config)
	e.Use(echo.WrapMiddleware(hmacMiddleware))
	if logger != nil {
		logger.Info("HMAC middleware configured successfully")
	}
}
