package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB //a global DB connection your whole app will use

func InitDB() { //opens the SQLite file (creates it if it doesn't exist)
	var err error
	DB, err = sql.Open("sqlite3", "./linkvault.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected!")
	createTables()
}

func createTables() { //creates the links table with IF NOT EXISTS so it's safe to run every startup
	query := `
	CREATE TABLE IF NOT EXISTS links (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		original      TEXT NOT NULL,
		short         TEXT NOT NULL UNIQUE,
		clicks        INTEGER DEFAULT 0,
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_accessed DATETIME
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Add last_accessed column if it doesn't exist yet (for existing databases)
	DB.Exec("ALTER TABLE links ADD COLUMN last_accessed DATETIME")

	fmt.Println("Tables ready!")
}