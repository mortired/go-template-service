package tracing

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	TraceIDKey    = "trace_id"
	TraceIDHeader = "X-Trace-ID"
)

// TraceIDMiddleware adds unique TraceID to each request
func TraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Generate TraceID if it doesn't exist
			traceID := c.Request().Header.Get(TraceIDHeader)
			if traceID == "" {
				traceID = generateTraceID()
			}

			// Save in context
			c.Set(TraceIDKey, traceID)

			// Add to response headers
			c.Response().Header().Set(TraceIDHeader, traceID)

			return next(c)
		}
	}
}

// GetTraceID gets TraceID from context
func GetTraceID(c echo.Context) string {
	if traceID, ok := c.Get(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// TraceIDLoggerMiddleware adds TraceID to logger context
func TraceIDLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get TraceID from context
			traceID := GetTraceID(c)

			// Add TraceID to logger context so it can be used in format string
			c.Set("trace_id", traceID)

			return next(c)
		}
	}
}

// CustomLoggerMiddleware creates a custom logger that includes TraceID
func CustomLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start)

			// Get TraceID from context
			traceID := GetTraceID(c)

			// Log request details with TraceID
			log.Printf(`{"time":"%s","id":"%s","trace_id":"%s","remote_ip":"%s","host":"%s","method":"%s","uri":"%s","user_agent":"%s","status":%d,"error":"%s","latency":%d,"latency_human":"%s","bytes_in":%d,"bytes_out":%d}`+"\n",
				time.Now().Format(time.RFC3339),
				c.Response().Header().Get("X-Request-ID"),
				traceID,
				c.RealIP(),
				c.Request().Host,
				c.Request().Method,
				c.Request().URL.Path,
				c.Request().UserAgent(),
				c.Response().Status,
				"",
				latency.Nanoseconds(),
				latency.String(),
				c.Request().ContentLength,
				c.Response().Size,
			)

			return err
		}
	}
}

// generateTraceID generates unique TraceID
func generateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
