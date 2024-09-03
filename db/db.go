package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewConnection creates a new connection to the PostgreSQL database via Supabase
func NewConnection() (*gorm.DB, error) {
	// Get the Supabase connection URL from the environment variables
	dsn := os.Getenv("SUPABASE_DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("SUPABASE_DB_URL environment variable not set")
	}

	fmt.Println("Connecting to database with DSN:", dsn) // Debugging print

	// Open the connection to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}
