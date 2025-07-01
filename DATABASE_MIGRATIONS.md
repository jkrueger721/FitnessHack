# Database Migrations Guide

This guide covers how to use the database migration CLI for the fitness application.

## Quick Start

### Prerequisites
1. Make sure your database is running:
   ```bash
   docker-compose up -d
   ```

2. Set your database connection environment variables:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=password
   export DB_NAME=fitness_db
   export DB_SSL_MODE=disable
   ```

### Running Migrations

```bash
# Run all pending migrations
go run migrate.go
# or
./migrate
# or
go run cmd/migrate/main.go

# Check migration status
go run migrate.go status
# or
./migrate status
# or
go run cmd/migrate/main.go status

# Generate Go models from current schema
go run migrate.go generate-models
# or
./migrate generate-models
# or
go run cmd/migrate/main.go generate-models

# Create a new migration file (provide a name or filename)
go run migrate.go create-migration "add user profiles"
go run migrate.go create-migration add_user_profiles.sql
go run migrate.go create-migration add-user-profiles
# or
./migrate create-migration "add user profiles"
./migrate create-migration add_user_profiles.sql
./migrate create-migration add-user-profiles
# or
go run cmd/migrate/main.go create-migration "add user profiles"
go run cmd/migrate/main.go create-migration add_user_profiles.sql
go run cmd/migrate/main.go create-migration add-user-profiles
```

## Migration Files

Migrations are stored in `