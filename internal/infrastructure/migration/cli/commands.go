package cli

import (
	"fmt"

	"users/internal/infrastructure/migration"
)

// Commands contains all available migration commands
type Commands struct {
	migrator *migration.Migrator
}

// NewCommands creates a new commands instance
func NewCommands(migrator *migration.Migrator) *Commands {
	return &Commands{
		migrator: migrator,
	}
}

// Up applies all migrations
func (c *Commands) Up() error {
	return c.migrator.Up()
}

// Down rolls back all migrations
func (c *Commands) Down() error {
	return c.migrator.Down()
}

// Force forcefully sets migration version
func (c *Commands) Force(version int) error {
	return c.migrator.Force(version)
}

// Status shows migration status
func (c *Commands) Status() error {
	return c.migrator.Status()
}

// Version shows current migration version
func (c *Commands) Version() error {
	version, dirty, err := c.migrator.Version()
	if err != nil {
		return fmt.Errorf("error getting version: %w", err)
	}

	fmt.Printf("📊 Current version: %d\n", version)
	if dirty {
		fmt.Println("⚠️  Database in unstable state (dirty)")
	} else {
		fmt.Println("✅ Database in stable state")
	}

	return nil
}

// Help shows command help
func (c *Commands) Help() {
	ShowUsage()
}
