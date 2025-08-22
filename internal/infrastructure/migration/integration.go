package migration

import (
	"context"
	"fmt"
	"log"
	"time"
)

// IntegrationOptions options for migration integration with application
type IntegrationOptions struct {
	AutoMigrate   bool          // Automatically apply migrations on startup
	Timeout       time.Duration // Timeout for migrations
	RetryAttempts int           // Number of reconnection attempts
	RetryInterval time.Duration // Interval between attempts
}

// DefaultIntegrationOptions returns default options
func DefaultIntegrationOptions() *IntegrationOptions {
	return &IntegrationOptions{
		AutoMigrate:   false,            // Disabled by default
		Timeout:       30 * time.Second, // 30 seconds
		RetryAttempts: 3,                // 3 attempts
		RetryInterval: 5 * time.Second,  // 5 seconds between attempts
	}
}

// IntegrationManager manages migration integration with application
type IntegrationManager struct {
	manager *Migrator
	options *IntegrationOptions
}

// NewIntegrationManager creates a new integration manager
func NewIntegrationManager(manager *Migrator, options *IntegrationOptions) *IntegrationManager {
	if options == nil {
		options = DefaultIntegrationOptions()
	}

	return &IntegrationManager{
		manager: manager,
		options: options,
	}
}

// EnsureMigrations applies migrations if necessary
func (im *IntegrationManager) EnsureMigrations(ctx context.Context) error {
	if !im.options.AutoMigrate {
		log.Println("📝 Automatic migrations disabled")
		return nil
	}

	log.Println("🚀 Checking and applying migrations...")

	// Check migration status
	if err := im.manager.Status(); err != nil {
		return fmt.Errorf("error checking migration status: %w", err)
	}

	// Apply migrations with timeout
	ctx, cancel := context.WithTimeout(ctx, im.options.Timeout)
	defer cancel()

	// Create channel for result
	result := make(chan error, 1)

	go func() {
		result <- im.manager.Up()
	}()

	// Wait for result or timeout
	select {
	case err := <-result:
		if err != nil {
			return fmt.Errorf("error applying migrations: %w", err)
		}
		log.Println("✅ Migrations successfully applied")
		return nil

	case <-ctx.Done():
		return fmt.Errorf("migration timeout: %w", ctx.Err())
	}
}

// WaitForDatabase waits for database availability
func (im *IntegrationManager) WaitForDatabase(ctx context.Context) error {
	log.Println("⏳ Waiting for database availability...")

	for attempt := 1; attempt <= im.options.RetryAttempts; attempt++ {
		// Check migration status (this will test DB connection)
		if err := im.manager.Status(); err == nil {
			log.Println("✅ Database available")
			return nil
		}

		log.Printf("⚠️  Attempt %d/%d: database unavailable", attempt, im.options.RetryAttempts)

		if attempt < im.options.RetryAttempts {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(im.options.RetryInterval):
				continue
			}
		}
	}

	return fmt.Errorf("database unavailable after %d attempts", im.options.RetryAttempts)
}

// HealthCheck checks migration system health
func (im *IntegrationManager) HealthCheck(ctx context.Context) error {
	// Check migration status
	if err := im.manager.Status(); err != nil {
		return fmt.Errorf("error checking migration status: %w", err)
	}

	// Check version
	version, dirty, err := im.manager.Version()
	if err != nil {
		return fmt.Errorf("error getting migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database in unstable state (dirty)")
	}

	log.Printf("✅ Migrations healthy, version: %d", version)
	return nil
}
