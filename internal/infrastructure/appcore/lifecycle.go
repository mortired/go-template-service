package appcore

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"users/internal/infrastructure/logging"

	"go.uber.org/fx"
)

// LifecycleHooks provides lifecycle hooks for the application
type LifecycleHooks struct {
	logger *logging.Logger
}

// NewLifecycleHooks creates new lifecycle hooks
func NewLifecycleHooks(logger *logging.Logger) *LifecycleHooks {
	return &LifecycleHooks{logger: logger}
}

// OnStart logs application start
func (lh *LifecycleHooks) OnStart(ctx context.Context) error {
	if lh.logger != nil {
		lh.logger.Info("AppCore: Starting application...")
	}
	return nil
}

// OnStop logs application stop
func (lh *LifecycleHooks) OnStop(ctx context.Context) error {
	if lh.logger != nil {
		lh.logger.Info("AppCore: Application stopped gracefully")
	}
	return nil
}

// SetupGracefulShutdown sets up graceful shutdown handling
func SetupGracefulShutdown(lifecycle fx.Lifecycle, logger *logging.Logger) {
	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		if logger != nil {
			logger.Info("AppCore: Received shutdown signal, stopping gracefully...")
		}
		os.Exit(0)
	}()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if logger != nil {
				logger.Info("AppCore: Starting application in background...")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if logger != nil {
				logger.Info("AppCore: Stopping application...")
			}
			return nil
		},
	})
}
