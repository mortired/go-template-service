package appcore

import (
	"log"
	"time"
	"users/internal/infrastructure/tracing"

	"github.com/labstack/echo/v4"
)

// LoggerMiddleware creates a custom logger middleware that includes TraceID
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start)

			// Get TraceID from context
			traceID := tracing.GetTraceID(c)

			// Log request details with TraceID
			log.Printf("[%s] %s %s %d %v",
				traceID,
				c.Request().Method,
				c.Request().URL.Path,
				c.Response().Status,
				latency,
			)

			return err
		}
	}
}
