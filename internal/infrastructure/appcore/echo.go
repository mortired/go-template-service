package appcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"users/internal/infrastructure/appcore/config"
	"users/internal/infrastructure/logging"
	"users/internal/infrastructure/response"
	"users/internal/infrastructure/tracing"

	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

// ProvideEchoServer creates HTTP Echo server
func ProvideEchoServer() *echo.Echo {
	e := echo.New()

	// Configure middleware with custom logger that includes TraceID
	e.Use(tracing.TraceIDMiddleware())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Configure HTTP error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var statusCode int
		var message string

		if he, ok := err.(*echo.HTTPError); ok {
			statusCode = he.Code
			if he.Message != nil {
				message = fmt.Sprintf("%v", he.Message)
			} else {
				message = http.StatusText(he.Code)
			}
		} else {
			statusCode = http.StatusInternalServerError
			message = "Internal server error"
		}

		// Create problem response
		problem := response.NewProblem(statusCode, message).
			WithType(response.TypeNotFound).
			WithInstance(c.Request().URL.Path).
			WithTraceID(tracing.GetTraceID(c))

		// Set appropriate problem type based on status code
		switch statusCode {
		case http.StatusNotFound:
			problem.WithType(response.TypeNotFound)
		case http.StatusUnauthorized:
			problem.WithType(response.TypeUnauthorized)
		case http.StatusForbidden:
			problem.WithType(response.TypeForbidden)
		case http.StatusBadRequest:
			problem.WithType(response.TypeInvalidRequest)
		case http.StatusConflict:
			problem.WithType(response.TypeConflict)
		default:
			problem.WithType(response.TypeInternalError)
		}

		// Send JSON response
		if err := c.JSON(statusCode, problem); err != nil {
			c.Logger().Error(err)
		}
	}

	return e
}

// SetupEchoMiddleware sets up Echo middleware with logger
func SetupEchoMiddleware(e *echo.Echo, logger *logging.Logger) {
	// Only error logging remains for critical errors
	e.Use(logging.ErrorLoggerMiddleware(logger))
}

// StartEchoServer starts HTTP server
func StartEchoServer(lifecycle fx.Lifecycle, e *echo.Echo, cfg *config.Config, logger *logging.Logger) {
	port := strconv.Itoa(cfg.Server.Port)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
					log.Printf("ERROR: HTTP server error: %v", err)
				}
			}()
			if logger != nil {
				logger.Info("HTTP server started", zap.String("port", port))
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if logger != nil {
				logger.Info("Stopping HTTP server...")
			}
			if err := e.Shutdown(ctx); err != nil {
				log.Printf("ERROR: Failed to shutdown HTTP server gracefully: %v", err)
				return err
			}
			if logger != nil {
				logger.Info("HTTP server stopped gracefully")
			}
			return nil
		},
	})
}

// StopEchoServer gracefully stops HTTP server
func StopEchoServer(e *echo.Echo, logger *logging.Logger) {
	if e != nil {
		if logger != nil {
			logger.Info("Stopping HTTP server...")
		}
		if err := e.Shutdown(context.Background()); err != nil {
			log.Printf("ERROR: Failed to shutdown HTTP server gracefully: %v", err)
		} else {
			if logger != nil {
				logger.Info("HTTP server stopped gracefully")
			}
		}
	}
}
