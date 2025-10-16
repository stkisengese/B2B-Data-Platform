package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// NewDatabaseConnection establishes a new database connection to the SQLite database
func NewDatabaseConnection(path string) (*sqlx.DB, error) {
	// Use WAL journal mode for better concurrency performance with SQLite
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL", path))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}
