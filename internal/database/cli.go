package database

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	case "create-migration":
		if len(args) < 2 {
			return fmt.Errorf("usage: create-migration <name or filename>. Example: create-migration add_user_profiles.sql")
		}
		return c.createMigration(args[1])
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// createMigration creates a new migration file with proper naming standards
func (c *CLI) createMigration(input string) error {
	// Remove .sql extension if present, and clean/format the name
	name := input
	if strings.HasSuffix(strings.ToLower(name), ".sql") {
		name = name[:len(name)-4]
	}
	name = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	name = strings.ReplaceAll(name, "-", "_")
	// Remove any special characters
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, name)
	if name == "" {
		return fmt.Errorf("invalid migration name")
	}

	// Get the next migration number
	nextNumber, err := c.getNextMigrationNumber()
	if err != nil {
		return fmt.Errorf("failed to get next migration number: %w", err)
	}

	// Create the filename
	filename := fmt.Sprintf("%03d_%s.sql", nextNumber, name)
	filepath := filepath.Join(DefaultMigrationsDir(), filename)

	// Create the migration content
	content := fmt.Sprintf(`-- Migration: %s
-- Description: %s
-- Date: %s

-- Add your SQL statements here
-- Example:
-- CREATE TABLE IF NOT EXISTS table_name (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );

-- CREATE INDEX IF NOT EXISTS idx_table_name_column ON table_name(column);
`, filename, strings.ReplaceAll(name, "_", " "), time.Now().Format("2006-01-02"))

	// Create the file
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	fmt.Printf("Created migration file: %s\n", filepath)
	fmt.Printf("Edit the file to add your SQL statements.\n")
	return nil
}

// getNextMigrationNumber determines the next migration number by reading existing files
func (c *CLI) getNextMigrationNumber() (int, error) {
	migrationsDir := DefaultMigrationsDir()

	// Read the migrations directory
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	maxNumber := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Extract number from filename (e.g., "001_create_users_table.sql" -> 1)
		if len(entry.Name()) >= 3 {
			var fileNumber int
			if num, err := fmt.Sscanf(entry.Name()[:3], "%d", &fileNumber); err == nil && num > 0 {
				// Keep track of the highest number
				if fileNumber > maxNumber {
					maxNumber = fileNumber
				}
			}
		}
	}

	return maxNumber + 1, nil
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
		fmt.Println("  migrate                    - Run all pending migrations")
		fmt.Println("  generate-models            - Generate Go models from database schema")
		fmt.Println("  status                     - Show migration status")
		fmt.Println("  create-migration <name or filename> - Create a new migration file (e.g. add_user_profiles.sql or \"add user profiles\")")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  create-migration add user profiles")
		fmt.Println("  create-migration add_user_profiles.sql")
		fmt.Println("  create-migration add-user-profiles")
		return nil
	}

	// Initialize database service
	dbService := New()
	db := dbService.GetDB()
	defer dbService.Close()

	cli := NewCLI(db)
	return cli.Run(args)
}
