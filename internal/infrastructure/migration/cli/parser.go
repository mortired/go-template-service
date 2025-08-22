package cli

import (
	"flag"
	"fmt"
)

// Args contains command line arguments
type Args struct {
	DatabaseURL    string
	MigrationsPath string
	Command        string
	Version        int
}

// ParseArgs parses command line arguments
func ParseArgs() (*Args, error) {
	// Parse command line flags
	var (
		databaseURL    = flag.String("database", "", "Database connection URL (e.g., postgres://user:pass@localhost:5432/dbname?sslmode=disable)")
		migrationsPath = flag.String("path", "./migrations", "Path to migrations folder")
		command        = flag.String("command", "up", "Command: up, down, force, status, version, help")
		version        = flag.Int("version", 0, "Version for force command")
	)
	flag.Parse()

	// Check command validity
	if !isValidCommand(*command) {
		return nil, fmt.Errorf("unknown command: %s. Supported commands: up, down, force, status, version, help", *command)
	}

	// Check version for force command
	if *command == "force" && *version == 0 {
		return nil, fmt.Errorf("force command requires version flag -version")
	}

	return &Args{
		DatabaseURL:    *databaseURL,
		MigrationsPath: *migrationsPath,
		Command:        *command,
		Version:        *version,
	}, nil
}

// isValidCommand checks command validity
func isValidCommand(cmd string) bool {
	validCommands := map[string]bool{
		"up":      true,
		"down":    true,
		"force":   true,
		"status":  true,
		"version": true,
		"help":    true,
	}
	return validCommands[cmd]
}

// ShowUsage shows complete usage information
func ShowUsage() {
	fmt.Print(`
🚀 Users Migration Tool

Usage:
  users-migrate [flags]

Available Commands:
  up       - Apply all migrations
  down     - Rollback all migrations
  force    - Force set version (requires -version)
  status   - Show migration status
  version  - Show current version
  help     - Show this help

Flags:
  -database string
        Database connection URL (optional)
        If not specified, environment variables from .env file are used
        Example: postgres://user:pass@localhost:5432/dbname?sslmode=disable
  -path string
        Path to migrations folder (default: ./migrations)
  -command string
        Command to execute (default: up)
  -version int
        Version for force command

Usage Examples:
  # Using environment variables (recommended)
  users-migrate -command=up
  users-migrate -command=status
  
  # Using direct URL
  users-migrate -database="postgres://postgres:pass@localhost:5432/users" -command=up
  users-migrate -database="postgres://postgres:pass@localhost:5432/users" -command=status
  users-migrate -database="postgres://postgres:pass@localhost:5432/users" -command=force -version=1

Environment Variables (if -database is not specified):
  DB_HOST=localhost
  DB_PORT=5432
  DB_USERNAME=postgres
  DB_PASSWORD=password
  DB_NAME=users
  MIGRATIONS_PATH=./migrations
`)
}
