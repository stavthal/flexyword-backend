package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Make a new connection to an SQLite database
// TODO: Add support for other databases (PostgreSQL)
func NewConnection() (*gorm.DB) {
	db, err := gorm.Open(sqlite.Open("./dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}