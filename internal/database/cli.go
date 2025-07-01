package database

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// CLI handles command-line operations for database management
type CLI struct {
	db *sqlx.DB
}

// NewCLI creates a new CLI instance
func NewCLI(db *sqlx.DB) *CLI {
	return &CLI{db: db}
}

// Run executes the CLI based on command line arguments
func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no command specified")
	}

	command := args[0]
	switch command {
	case "migrate":
		return c.runMigrations()
	case "generate-models":
		return c.generateModels()
	case "status":
		return c.showStatus()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// runMigrations runs all pending migrations
func (c *CLI) runMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Running migrations...")
	if err := RunMigrations(ctx, c.db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// generateModels generates Go models from the current database schema
func (c *CLI) generateModels() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Generating models from database schema...")
	if err := GenerateModelsFromDB(ctx, c.db); err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	log.Println("Models generated successfully")
	return nil
}

// showStatus shows the current migration status
func (c *CLI) showStatus() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	manager := NewMigrationManager(c.db)

	// Get applied migrations
	applied, err := manager.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("=================")

	if len(applied) == 0 {
		fmt.Println("No migrations applied yet.")
		return nil
	}

	for _, migration := range applied {
		fmt.Printf("✓ %s (applied at %s)\n", migration.Name, migration.AppliedAt.Format("2006-01-02 15:04:05"))
	}

	// Check for pending migrations
	appliedMap := make(map[string]bool)
	for _, migration := range applied {
		appliedMap[migration.Name] = true
	}

	// Load migration files to check for pending ones
	migrationFiles, err := manager.LoadMigrationFiles(DefaultMigrationsDir())
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	var pending []string
	for _, migrationFile := range migrationFiles {
		if !appliedMap[migrationFile.Name] {
			pending = append(pending, migrationFile.Name)
		}
	}

	if len(pending) > 0 {
		fmt.Println("\nPending migrations:")
		for _, name := range pending {
			fmt.Printf("  - %s\n", name)
		}
	} else {
		fmt.Println("\nAll migrations are up to date.")
	}

	return nil
}

// RunCLI is a convenience function to run the CLI with the database service
func RunCLI() error {
	// Parse command line flags
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Database CLI Usage:")
		fmt.Println("  migrate         - Run all pending migrations")
		fmt.Println("  generate-models - Generate Go models from database schema")
		fmt.Println("  status          - Show migration status")
		return nil
	}

	// Initialize database service
	dbService := New()
	db := dbService.GetDB()
	defer dbService.Close()

	cli := NewCLI(db)
	return cli.Run(args)
}
