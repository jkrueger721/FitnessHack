package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

// DefaultMigrationsDir returns the default migrations directory path
func DefaultMigrationsDir() string {
	return "internal/database/migrations"
}

// RunMigrations runs all pending migrations from SQL files
func RunMigrations(ctx context.Context, db *sqlx.DB) error {
	return RunMigrationsFromDir(ctx, db, DefaultMigrationsDir())
}

// RunMigrationsFromDir runs migrations from a specific directory
func RunMigrationsFromDir(ctx context.Context, db *sqlx.DB, migrationsDir string) error {
	manager := NewMigrationManager(db)
	return manager.RunMigrations(ctx, migrationsDir)
}

// GenerateModelsFromDB generates Go models from the current database schema
func GenerateModelsFromDB(ctx context.Context, db *sqlx.DB) error {
	manager := NewMigrationManager(db)
	outputPath := filepath.Join("internal", "database", "models.go")
	return manager.GenerateModels(ctx, outputPath)
}

// CreateMigrationFile creates a new migration file with the given name
func CreateMigrationFile(name, sql string) error {
	migrationsDir := DefaultMigrationsDir()

	// Ensure migrations directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create the migration file
	filename := filepath.Join(migrationsDir, name+".sql")
	if err := os.WriteFile(filename, []byte(sql), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	log.Printf("Created migration file: %s", filename)
	return nil
}
