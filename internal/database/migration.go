package database

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
)

// Migration represents a database migration
type Migration struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	AppliedAt time.Time `db:"applied_at"`
}

// MigrationFile represents a migration file
type MigrationFile struct {
	Name     string
	Path     string
	SQL      string
	Filename string
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db *sqlx.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sqlx.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// InitMigrationsTable creates the migrations table if it doesn't exist
func (m *MigrationManager) InitMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := m.db.ExecContext(ctx, query)
	return err
}

// GetAppliedMigrations returns all applied migrations
func (m *MigrationManager) GetAppliedMigrations(ctx context.Context) ([]Migration, error) {
	var migrations []Migration
	query := `SELECT id, name, applied_at FROM migrations ORDER BY id ASC`
	err := m.db.SelectContext(ctx, &migrations, query)
	return migrations, err
}

// ApplyMigration applies a single migration
func (m *MigrationManager) ApplyMigration(ctx context.Context, name, sql string) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute the migration SQL
	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", name, err)
	}

	// Record the migration
	_, err = tx.ExecContext(ctx, "INSERT INTO migrations (name) VALUES ($1)", name)
	if err != nil {
		return fmt.Errorf("failed to record migration %s: %w", name, err)
	}

	return tx.Commit()
}

// LoadMigrationFiles loads migration SQL files from the migrations directory
func (m *MigrationManager) LoadMigrationFiles(migrationsDir string) ([]MigrationFile, error) {
	var migrationFiles []MigrationFile

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Printf("Migrations directory does not exist: %s", migrationsDir)
		return migrationFiles, nil
	}

	// Walk through the migrations directory
	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-SQL files
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".sql") {
			return nil
		}

		// Read the SQL file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		// Extract migration name from filename (remove .sql extension)
		name := strings.TrimSuffix(d.Name(), ".sql")

		migrationFile := MigrationFile{
			Name:     name,
			Path:     path,
			SQL:      string(content),
			Filename: d.Name(),
		}

		migrationFiles = append(migrationFiles, migrationFile)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk migrations directory: %w", err)
	}

	// Sort migration files by name to ensure proper order
	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Name < migrationFiles[j].Name
	})

	return migrationFiles, nil
}

// RunMigrations runs all pending migrations from SQL files
func (m *MigrationManager) RunMigrations(ctx context.Context, migrationsDir string) error {
	// Initialize migrations table
	if err := m.InitMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Load migration files
	migrationFiles, err := m.LoadMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	if len(migrationFiles) == 0 {
		log.Println("No migration files found")
		return nil
	}

	// Get applied migrations
	applied, err := m.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, migration := range applied {
		appliedMap[migration.Name] = true
	}

	// Apply pending migrations
	for _, migrationFile := range migrationFiles {
		if !appliedMap[migrationFile.Name] {
			log.Printf("Applying migration: %s", migrationFile.Name)
			if err := m.ApplyMigration(ctx, migrationFile.Name, migrationFile.SQL); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migrationFile.Name, err)
			}
			log.Printf("Applied migration: %s", migrationFile.Name)
		}
	}

	return nil
}

// GenerateModels generates Go models from the current database schema
func (m *MigrationManager) GenerateModels(ctx context.Context, outputPath string) error {
	// Get all tables
	tables, err := m.getTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Generate models for each table
	var models []TableModel
	for _, table := range tables {
		columns, err := m.getColumns(ctx, table)
		if err != nil {
			return fmt.Errorf("failed to get columns for table %s: %w", table, err)
		}

		model := TableModel{
			Name:    table,
			Columns: columns,
		}
		models = append(models, model)
	}

	// Generate the Go file
	return m.generateGoFile(models, outputPath)
}

// TableModel represents a database table for model generation
type TableModel struct {
	Name    string
	Columns []Column
}

// Column represents a database column
type Column struct {
	Name       string
	Type       string
	IsNullable bool
	IsPrimary  bool
	IsUnique   bool
	Default    *string
}

// getTables returns all table names in the current schema
func (m *MigrationManager) getTables(ctx context.Context) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = current_schema() 
		AND table_type = 'BASE TABLE'
		AND table_name != 'migrations'
		ORDER BY table_name
	`
	var tables []string
	err := m.db.SelectContext(ctx, &tables, query)
	return tables, err
}

// getColumns returns column information for a table
func (m *MigrationManager) getColumns(ctx context.Context, tableName string) ([]Column, error) {
	query := `
		SELECT 
			c.column_name,
			c.data_type,
			c.is_nullable,
			CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END as is_primary,
			CASE WHEN u.column_name IS NOT NULL THEN true ELSE false END as is_unique,
			c.column_default
		FROM information_schema.columns c
		LEFT JOIN (
			SELECT ku.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage ku ON tc.constraint_name = ku.constraint_name
			WHERE tc.constraint_type = 'PRIMARY KEY' AND ku.table_name = $1
		) pk ON c.column_name = pk.column_name
		LEFT JOIN (
			SELECT ku.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage ku ON tc.constraint_name = ku.constraint_name
			WHERE tc.constraint_type = 'UNIQUE' AND ku.table_name = $1
		) u ON c.column_name = u.column_name
		WHERE c.table_name = $1
		ORDER BY c.ordinal_position
	`

	rows, err := m.db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var col Column
		var isNullable, isPrimary, isUnique string
		var defaultVal *string

		err := rows.Scan(&col.Name, &col.Type, &isNullable, &isPrimary, &isUnique, &defaultVal)
		if err != nil {
			return nil, err
		}

		col.IsNullable = isNullable == "YES"
		col.IsPrimary = isPrimary == "true"
		col.IsUnique = isUnique == "true"
		col.Default = defaultVal
		col.Type = m.mapSQLTypeToGoType(col.Type)

		columns = append(columns, col)
	}

	return columns, nil
}

// mapSQLTypeToGoType maps PostgreSQL types to Go types
func (m *MigrationManager) mapSQLTypeToGoType(sqlType string) string {
	switch strings.ToLower(sqlType) {
	case "uuid":
		return "string"
	case "varchar", "text", "char":
		return "string"
	case "integer", "int", "int4":
		return "int"
	case "bigint", "int8":
		return "int64"
	case "smallint", "int2":
		return "int16"
	case "decimal", "numeric":
		return "decimal.Decimal"
	case "real", "float4":
		return "float32"
	case "double precision", "float8":
		return "float64"
	case "boolean", "bool":
		return "bool"
	case "timestamp with time zone", "timestamptz":
		return "time.Time"
	case "timestamp without time zone", "timestamp":
		return "time.Time"
	case "date":
		return "time.Time"
	case "json", "jsonb":
		return "json.RawMessage"
	default:
		return "interface{}"
	}
}

// generateGoFile generates the Go models file
func (m *MigrationManager) generateGoFile(models []TableModel, outputPath string) error {
	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create template with functions
	funcMap := template.FuncMap{
		"title": strings.Title,
		"snake": m.toSnakeCase,
	}

	// Parse and execute the template
	tmpl, err := template.New("models").Funcs(funcMap).Parse(modelsTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		Models []TableModel
		Time   string
	}{
		Models: models,
		Time:   time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	log.Printf("Generated models file: %s", outputPath)
	return nil
}

// toSnakeCase converts camelCase to snake_case
func (m *MigrationManager) toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// modelsTemplate is the Go template for generating model files
const modelsTemplate = `// Code generated by migration system on {{.Time}}
// DO NOT EDIT THIS FILE MANUALLY

package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

{{range .Models}}
// {{.Name | title}} represents the {{.Name}} table
type {{.Name | title}} struct {
{{range .Columns}}	{{.Name | title}} {{.Type}} ` + "`" + `db:"{{.Name}}" json:"{{.Name | snake}}"` + "`" + `{{if .IsPrimary}} // Primary key{{end}}{{if .IsUnique}} // Unique{{end}}{{if .Default}} // Default: {{.Default}}{{end}}
{{end}}}

// TableName returns the table name for {{.Name | title}}
func ({{.Name | title}}) TableName() string {
	return "{{.Name}}"
}

// Scan implements the sql.Scanner interface for {{.Name | title}}
func (m *{{.Name | title}}) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	default:
		return fmt.Errorf("cannot scan %T into {{.Name | title}}", value)
	}
}

// Value implements the driver.Valuer interface for {{.Name | title}}
func (m {{.Name | title}}) Value() (driver.Value, error) {
	return json.Marshal(m)
}
{{end}}

// Custom types for better type safety
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return fmt.Errorf("cannot scan %T into JSONMap", value)
	}
}
`
