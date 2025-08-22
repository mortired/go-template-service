package appcore

import (
	"time"
	"users/internal/infrastructure/logging"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// FxEventsConfig structure for fx events configuration
type FxEventsConfig struct {
	FxEnabled bool
	FxDebug   bool
	FxVerbose bool
	ShowGraph bool
}

// NewFxEvents creates and returns fx options for events configuration
func NewFxEvents(cfg FxEventsConfig) []fx.Option {
	var options []fx.Option

	// Determine which logger to use
	if !cfg.FxEnabled {
		// Completely disable fx logging
		options = append(options, fx.WithLogger(fx.NopLogger))
	} else if cfg.FxDebug || cfg.FxVerbose || cfg.ShowGraph {
		// Use custom logger for debugging
		options = append(options, fx.WithLogger(NewFxDebugLogger))
	}
	// If FxEnabled=true, but other flags are false, use standard fx logger

	// Add timeout options for debugging
	if cfg.FxDebug || cfg.FxVerbose {
		// Set reasonable timeouts for debug mode instead of 0
		options = append(options, fx.StartTimeout(30*time.Second))
		options = append(options, fx.StopTimeout(30*time.Second))
	}

	return options
}

// NewFxDebugLogger creates custom logger for fx debugging
func NewFxDebugLogger(logger *logging.Logger) fxevent.Logger {
	return &fxDebugLogger{logger: logger}
}

// fxDebugLogger - custom logger for fx debugging
type fxDebugLogger struct {
	logger *logging.Logger
}

func (l *fxDebugLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		if l.logger != nil {
			l.logger.Debug("Starting", zap.String("function", e.FunctionName))
		}
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			if l.logger != nil {
				l.logger.Error("Start failed", zap.String("function", e.FunctionName), zap.Error(e.Err))
			}
		} else {
			if l.logger != nil {
				l.logger.Debug("Started", zap.String("function", e.FunctionName), zap.Duration("runtime", e.Runtime))
			}
		}
	case *fxevent.OnStopExecuting:
		if l.logger != nil {
			l.logger.Debug("Stopping", zap.String("function", e.FunctionName))
		}
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			if l.logger != nil {
				l.logger.Error("Stop failed", zap.String("function", e.FunctionName), zap.Error(e.Err))
			}
		} else {
			if l.logger != nil {
				l.logger.Debug("Stopped", zap.String("function", e.FunctionName))
			}
		}
	case *fxevent.Supplied:
		if l.logger != nil {
			l.logger.Debug("Supplied", zap.String("type", e.TypeName))
		}
	case *fxevent.Provided:
		if l.logger != nil {
			l.logger.Debug("Provided", zap.String("constructor", e.ConstructorName))
		}
	case *fxevent.Invoking:
		if l.logger != nil {
			l.logger.Debug("Invoking", zap.String("function", e.FunctionName))
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			if l.logger != nil {
				l.logger.Error("Invoke failed", zap.String("function", e.FunctionName), zap.Error(e.Err))
			}
		} else {
			if l.logger != nil {
				l.logger.Debug("Invoked", zap.String("function", e.FunctionName))
			}
		}
	default:
		// Ignore other events to reduce noise
	}
}
