package cli

import (
	"fmt"
	"log"
	"os"

	"users/internal/infrastructure/migration"
)

// App represents CLI application for migration management
type App struct {
	commands *Commands
}

// NewApp creates a new CLI application
func NewApp(commands *Commands) *App {
	return &App{
		commands: commands,
	}
}

// Run runs the CLI application
func (app *App) Run(args *Args) error {
	switch args.Command {
	case "up":
		return app.commands.Up()

	case "down":
		return app.commands.Down()

	case "force":
		return app.commands.Force(args.Version)

	case "status":
		return app.commands.Status()

	case "version":
		return app.commands.Version()

	case "help":
		app.commands.Help()
		return nil

	default:
		return fmt.Errorf("unknown command: %s", args.Command)
	}
}

// Execute runs the CLI application with argument parsing
func Execute() {
	// Parse arguments
	args, err := ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Argument parsing error: %v\n\n", err)
		ShowUsage()
		os.Exit(1)
	}

	var cfg *migration.Config
	var manager *migration.Migrator

	// If database URL is not specified, use environment variables
	if args.DatabaseURL == "" {
		fmt.Println("🔧 Using environment variables for database configuration...")

		// Create migration manager with automatic environment variable loading
		manager, err = migration.NewWithEnv()
		if err != nil {
			log.Fatalf("❌ Error creating migration manager with environment variables: %v", err)
		}
	} else {
		// Create migration configuration from arguments
		cfg = &migration.Config{
			DatabaseURL:    args.DatabaseURL,
			MigrationsPath: args.MigrationsPath,
		}

		// Create migration manager
		manager, err = migration.New(cfg)
		if err != nil {
			log.Fatalf("❌ Error creating migration manager: %v", err)
		}
	}
	defer manager.Close()

	// Create commands
	commands := NewCommands(manager)

	// Create application
	app := NewApp(commands)

	// Run application
	if err := app.Run(args); err != nil {
		log.Fatalf("❌ Command execution error: %v", err)
	}
}
