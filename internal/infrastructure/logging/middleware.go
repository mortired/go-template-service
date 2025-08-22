package logging

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggerMiddleware creates Echo middleware that uses Zap logger
func LoggerMiddleware(logger *Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Create context with trace_id
			ctx := context.WithValue(c.Request().Context(), "trace_id", c.Get("trace_id"))

			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start)

			// Create logger with context
			log := logger.WithContext(ctx)

			// Log request details
			log.Info("HTTP Request",
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().URL.Path),
				zap.Int("status", c.Response().Status),
				zap.Duration("latency", latency),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
				zap.Int64("bytes_in", c.Request().ContentLength),
				zap.Int64("bytes_out", c.Response().Size),
			)

			return err
		}
	}
}

// ErrorLoggerMiddleware creates middleware for logging errors
func ErrorLoggerMiddleware(logger *Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err != nil {
				ctx := context.WithValue(c.Request().Context(), "trace_id", c.Get("trace_id"))
				log := logger.WithContext(ctx)

				log.Error("HTTP Error",
					zap.Error(err),
					zap.String("method", c.Request().Method),
					zap.String("uri", c.Request().URL.Path),
					zap.Int("status", c.Response().Status),
				)
			}

			return err
		}
	}
}
