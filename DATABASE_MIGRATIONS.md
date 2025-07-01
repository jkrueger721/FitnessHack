# Database Migrations Guide

This guide covers the database migration system for the fitness application. The system uses SQL files for version-controlled database schema changes and automatically generates Go models to match the database structure.

## Overview

The migration system provides:
- **Versioned migrations** using SQL files
- **Automatic model generation** from database schema
- **CLI interface** for running migrations and generating models
- **Transaction safety** with automatic rollback on failure
- **Migration tracking** to prevent duplicate applications

## Quick Start

```bash
# Run all pending migrations
go run cmd/main.go migrate

# Check migration status
go run cmd/main.go status

# Generate Go models from current schema
go run cmd/main.go generate-models
```

## Migration Files

### Current Migrations

Located in `internal/database/migrations/`:

1. **001_create_users_table.sql** - Creates the users table with authentication fields
2. **002_create_workouts_table.sql** - Creates the workouts table linked to users
3. **003_create_exercises_table.sql** - Creates the exercises table with exercise metadata
4. **004_create_workout_exercises_table.sql** - Creates the junction table for workout-exercise relationships
5. **005_create_workout_sessions_table.sql** - Creates the workout sessions table for tracking completed workouts

## Migration Naming Convention

Migrations follow this naming pattern:
```
{sequence_number}_{description}.sql
```

Examples:
- `001_create_users_table.sql`
- `002_add_user_profile_fields.sql`
- `003_create_indexes.sql`

## Migration File Structure

Each migration file should include:

```sql
-- Migration: {filename}
-- Description: Brief description of what this migration does
-- Date: YYYY-MM-DD

-- SQL statements here
CREATE TABLE IF NOT EXISTS table_name (
    -- table definition
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column);

-- Comments for documentation
COMMENT ON TABLE table_name IS 'Description of the table';
COMMENT ON COLUMN table_name.column IS 'Description of the column';
```

## Running Migrations

### Using the CLI

```bash
# Run all pending migrations
go run cmd/main.go migrate

# Check migration status
go run cmd/main.go status

# Generate models from current schema
go run cmd/main.go generate-models
```

### Programmatically

```go
import "fitness-hack/internal/database"

// Initialize database service
dbService := database.New()
db := dbService.GetDB()
defer dbService.Close()

ctx := context.Background()

// Run migrations
if err := database.RunMigrations(ctx, db); err != nil {
    log.Fatal("Failed to run migrations:", err)
}

// Generate models
if err := database.GenerateModelsFromDB(ctx, db); err != nil {
    log.Fatal("Failed to generate models:", err)
}
```

## Creating New Migrations

### Method 1: Create SQL File Manually

1. Create a new `.sql` file in the `internal/database/migrations/` directory
2. Follow the naming convention: `{next_sequence_number}_{description}.sql`
3. Include the migration header with description and date
4. Write your SQL statements
5. Add appropriate indexes and comments

Example:
```sql
-- Migration: 006_add_user_profiles
-- Description: Adds user profile information table
-- Date: 2024-01-02

CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    height_cm INTEGER,
    weight_kg DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);

COMMENT ON TABLE user_profiles IS 'Stores additional user profile information';
COMMENT ON COLUMN user_profiles.bio IS 'User biography or description';
COMMENT ON COLUMN user_profiles.height_cm IS 'User height in centimeters';
COMMENT ON COLUMN user_profiles.weight_kg IS 'User weight in kilograms';
```

### Method 2: Using the Helper Function

```go
// Create a new migration file programmatically
err := database.CreateMigrationFile("006_add_user_profiles", `
    CREATE TABLE IF NOT EXISTS user_profiles (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        bio TEXT,
        height_cm INTEGER,
        weight_kg DECIMAL(5,2),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );
`)
```

## Migration Best Practices

### 1. Always Use IF NOT EXISTS
```sql
-- Good
CREATE TABLE IF NOT EXISTS table_name (...);
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column);

-- Avoid
CREATE TABLE table_name (...);
```

### 2. Include Proper Foreign Key Constraints
```sql
-- Good
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

-- Avoid
user_id UUID NOT NULL,
```

### 3. Add Indexes for Performance
```sql
-- Add indexes for frequently queried columns
CREATE INDEX IF NOT EXISTS idx_table_column ON table_name(column);
CREATE INDEX IF NOT EXISTS idx_table_compound ON table_name(col1, col2);
```

### 4. Include Comments for Documentation
```sql
COMMENT ON TABLE table_name IS 'Description of the table';
COMMENT ON COLUMN table_name.column IS 'Description of the column';
```

### 5. Use Appropriate Data Types
```sql
-- Good
weight_kg DECIMAL(5,2),  -- Precise decimal for weight
duration_seconds INTEGER, -- Integer for seconds
created_at TIMESTAMP WITH TIME ZONE, -- Timezone-aware timestamps

-- Avoid
weight_kg FLOAT, -- Imprecise for financial/measurement data
duration VARCHAR(10), -- String for numeric data
created_at VARCHAR(20), -- String for date/time data
```

### 6. Handle Rollbacks (if needed)
For complex migrations, consider adding rollback logic:
```sql
-- Migration: 007_add_complex_feature
-- Description: Adds complex feature with rollback support
-- Date: 2024-01-03

-- Create new table
CREATE TABLE IF NOT EXISTS complex_feature (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL
);

-- Add rollback comment (for manual rollback if needed)
-- ROLLBACK: DROP TABLE IF EXISTS complex_feature;
```

## Model Generation

The system automatically generates Go models from your database schema:

### Generated Model Example
```go
// Users represents the users table
type Users struct {
    ID           string    `db:"id" json:"id"` // Primary key
    Email        string    `db:"email" json:"email"` // Unique
    Username     string    `db:"username" json:"username"` // Unique
    PasswordHash string    `db:"password_hash" json:"password_hash"`
    FirstName    *string   `db:"first_name" json:"first_name"`
    LastName     *string   `db:"last_name" json:"last_name"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for Users
func (Users) TableName() string {
    return "users"
}
```

### Type Mapping

| PostgreSQL | Go Type |
|------------|---------|
| `uuid` | `string` |
| `varchar`, `text` | `string` |
| `integer`, `int4` | `int` |
| `bigint`, `int8` | `int64` |
| `decimal`, `numeric` | `decimal.Decimal` |
| `boolean` | `bool` |
| `timestamp with time zone` | `time.Time` |
| `json`, `jsonb` | `json.RawMessage` |

## Migration Tracking

The system automatically tracks applied migrations in a `migrations` table:

```sql
CREATE TABLE migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

This ensures:
- Migrations are only applied once
- Migration history is preserved
- Rollback information is available

## Environment-Specific Migrations

For environment-specific migrations, you can use conditional logic:

```sql
-- Migration: 008_add_test_data
-- Description: Adds test data for development environment
-- Date: 2024-01-04

-- Only run in development/test environments
DO $$
BEGIN
    IF current_setting('app.environment') = 'development' THEN
        INSERT INTO exercises (name, description, muscle_group) VALUES
        ('Push-ups', 'Basic push-ups', 'Chest'),
        ('Squats', 'Basic squats', 'Legs');
    END IF;
END $$;
```

## Troubleshooting

### Migration Fails
1. Check the migration logs for specific error messages
2. Verify SQL syntax is correct
3. Ensure all referenced tables/columns exist
4. Check for constraint violations

### Migration Already Applied
- The system prevents duplicate migrations
- Check the `migrations` table to see applied migrations
- Use `go run cmd/main.go status` to see current status

### Model Generation Issues
- Ensure the database connection is working
- Verify the schema is accessible
- Check that tables exist before generating models

## Project Structure

```
fitness-hack/
├── DATABASE_MIGRATIONS.md          # This file
├── internal/
│   └── database/
│       ├── database.go             # Enhanced database service
│       ├── migration.go            # Core migration system
│       ├── migrations.go           # Migration utilities
│       ├── models.go               # Generated Go models
│       ├── cli.go                  # Command-line interface
│       ├── README.md               # Database documentation
│       └── migrations/             # SQL migration files
│           ├── 001_create_users_table.sql
│           ├── 002_create_workouts_table.sql
│           ├── 003_create_exercises_table.sql
│           ├── 004_create_workout_exercises_table.sql
│           └── 005_create_workout_sessions_table.sql
```

## Dependencies

The migration system requires:
- PostgreSQL database
- `github.com/jmoiron/sqlx` for database operations
- `github.com/shopspring/decimal` for precise decimal arithmetic

## Related Documentation

- `internal/database/README.md` - Database service documentation
- `internal/database/migrations/README.md` - Detailed migration examples 