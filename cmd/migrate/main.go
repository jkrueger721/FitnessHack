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
		fmt.Println("  migrate         - Run all pending migrations")
		fmt.Println("  status          - Show migration status")
		fmt.Println("  generate-models - Generate Go models from database schema")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run cmd/migrate/main.go migrate")
		fmt.Println("  go run cmd/migrate/main.go status")
		fmt.Println("  go run cmd/migrate/main.go generate-models")
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
