package main

import (
	"flag"
	"fmt"
	"log"

	"fitness-hack/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Parse command line flags
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Database Migration CLI")
		fmt.Println("======================")
		fmt.Println("Usage:")
		fmt.Println("  go run migrate.go                    - Run all pending migrations")
		fmt.Println("  go run migrate.go status             - Show migration status")
		fmt.Println("  go run migrate.go generate-models    - Generate Go models from database schema")
		fmt.Println("  go run migrate.go create-migration <name or filename> - Create a new migration file")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run migrate.go create-migration add user profiles")
		fmt.Println("  go run migrate.go create-migration add_user_profiles.sql")
		fmt.Println("  go run migrate.go create-migration add-user-profiles")
		return
	}

	// Initialize database service
	dbService := database.New()
	db := dbService.GetDB()
	defer dbService.Close()

	// Create CLI instance
	cli := database.NewCLI(db)

	// Run the command
	if err := cli.Run(args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
