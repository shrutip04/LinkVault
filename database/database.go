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
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		username   TEXT NOT NULL UNIQUE,
		email      TEXT NOT NULL UNIQUE,
		password   TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	linksTable := `
CREATE TABLE IF NOT EXISTS links (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL,
    original      TEXT NOT NULL,
    short         TEXT NOT NULL UNIQUE,
    clicks        INTEGER DEFAULT 0,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_accessed DATETIME,
    expires_at    DATETIME,
    password      TEXT,
    category      TEXT DEFAULT 'General',
    FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = DB.Exec(linksTable)
	if err != nil {
		log.Fatal("Failed to create links table:", err)
	}
	// Migrations for existing databases
	DB.Exec("ALTER TABLE links ADD COLUMN last_accessed DATETIME")
	DB.Exec("ALTER TABLE links ADD COLUMN expires_at DATETIME")
	DB.Exec("ALTER TABLE links ADD COLUMN user_id INTEGER")
	DB.Exec("ALTER TABLE links ADD COLUMN password TEXT")
	DB.Exec("ALTER TABLE links ADD COLUMN category TEXT DEFAULT 'General'")

	fmt.Println("Tables ready!")
}
