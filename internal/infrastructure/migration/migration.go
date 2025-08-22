package migration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Config configuration for migrations
type Config struct {
	DatabaseURL    string
	MigrationsPath string
}

// Migrator manages database migrations
type Migrator struct {
	config  *Config
	migrate *migrate.Migrate
}

// New creates a new migrator
func New(cfg *Config) (*Migrator, error) {
	// Check if migrations folder exists
	if _, err := os.Stat(cfg.MigrationsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations folder not found: %s", cfg.MigrationsPath)
	}

	// Create absolute path to migrations
	absPath, err := filepath.Abs(cfg.MigrationsPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}

	// Form URL for file migrations
	// Replace backslashes with forward slashes for correct URL on Windows
	fileURL := filepath.ToSlash(absPath)
	migrationsURL := fmt.Sprintf("file://%s", fileURL)

	// Create migrator
	m, err := migrate.New(migrationsURL, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating migrator: %w", err)
	}

	return &Migrator{
		config:  cfg,
		migrate: m,
	}, nil
}

// Up applies all migrations
func (m *Migrator) Up() error {
	log.Println("🚀 Applying migrations...")

	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	log.Println("✅ Migrations successfully applied")
	return nil
}

// Down rolls back all migrations
func (m *Migrator) Down() error {
	log.Println("🔄 Rolling back all migrations...")

	if err := m.migrate.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error rolling back migrations: %w", err)
	}

	log.Println("✅ Migrations successfully rolled back")
	return nil
}

// Force sets migration version
func (m *Migrator) Force(version int) error {
	log.Printf("🔧 Force setting migration version: %d", version)

	if err := m.migrate.Force(version); err != nil {
		return fmt.Errorf("error setting version: %w", err)
	}

	log.Printf("✅ Migration version set: %d", version)
	return nil
}

// Version returns current migration version
func (m *Migrator) Version() (int, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		return 0, false, fmt.Errorf("error getting version: %w", err)
	}

	return int(version), dirty, nil
}

// Status shows migration status
func (m *Migrator) Status() error {
	log.Println("📊 Migration status:")

	version, dirty, err := m.Version()
	if err != nil {
		return err
	}

	log.Printf("   Current version: %d", version)
	if dirty {
		log.Println("   ⚠️  Database in unstable state (dirty)")
	} else {
		log.Println("   ✅ Database in stable state")
	}

	return nil
}

// Close closes database connection
func (m *Migrator) Close() error {
	if m.migrate != nil {
		_, err := m.migrate.Close()
		return err
	}
	return nil
}
