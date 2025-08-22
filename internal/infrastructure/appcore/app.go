package appcore

import (
	"context"
	"os"
	"users/internal/infrastructure/config"

	"go.uber.org/fx"
)

// Application represents application with dependency injection
type Application struct {
	*fx.App
}

// New creates a new application with given options
func New(options ...fx.Option) *Application {
	appcoreOptions := []Option{
		// Core infrastructure modules
		LoggingModule,
		LifecycleModule,

		// Configuration
		Provide(ProvideConfig),
	}

	options = append(options, appcoreOptions...)

	// Configure logging through SERVICE_DEBUG
	appDebug := config.GetEnvOrDefault("SERVICE_DEBUG", "false")

	if appDebug == "true" {
		// Use fxevents package from appcore
		cfg := FxEventsConfig{
			FxEnabled: true,
			FxDebug:   true,
			FxVerbose: true,
			ShowGraph: false,
		}
		logOpts := NewFxEvents(cfg)
		options = append(options, logOpts...)
	}

	app := &Application{
		App: fx.New(options...),
	}

	return app
}

// Run starts the application
func (app *Application) Run() {
	// Start the application
	if err := app.Start(context.Background()); err != nil {
		os.Exit(1)
	}

	// Wait for shutdown signal (handled by lifecycle hooks)
	select {}
}

// Start starts the application in background mode
func (app *Application) Start(ctx context.Context) error {
	return app.App.Start(ctx)
}

// Stop stops the application
func (app *Application) Stop(ctx context.Context) error {
	return app.App.Stop(ctx)
}
