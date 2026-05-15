package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect(dbPath string) error {
    var err error
    DB, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        return fmt.Errorf("error opening database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        return fmt.Errorf("error connecting to database: %v", err)
    }

    if err = createTables(); err != nil {
        return fmt.Errorf("error creating tables: %v", err)
    }

    fmt.Println("SQLite database connected successfully")
    return nil
}

func createTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
        `CREATE TABLE IF NOT EXISTS links (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER REFERENCES users(id),
            original_url TEXT NOT NULL,
            short_code TEXT UNIQUE NOT NULL,
            clicks INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
        `CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code)`,
        `CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id)`,
    }

    for _, query := range queries {
        if _, err := DB.Exec(query); err != nil {
            return err
        }
    }

    return nil
}
