package database

import (
	"database/sql"
	"fmt"
	"log"

	"ADS4/internal/config"

	_ "github.com/lib/pq" // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}
//https://www.sqlite.org/pragma.html#pragma_synchronous 
func NewDB(cfg config.Config) (*DB, error) {
	log.Println("Connecting to database...")

	if cfg.DBtype == "sqlite" {
		db, err := sql.Open("sqlite3", fmt.Sprintf("file:"+cfg.DataDir+"/%s.db?cache=shared&_journal_mode=WAL", cfg.DBName)) //ADS4.db
		if err != nil {
			return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
		}

		log.Println("SQLite database connected successfully")
		return &DB{db}, nil
	}

	if cfg.DBtype == "postgres" {
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
		}

		// Ping the database to ensure connection is established
		if err := db.Ping(); err != nil {
			return nil, err
		}
		log.Println("Postgres database connected successfully")

		return &DB{db}, nil
	}
	return nil, fmt.Errorf("unsupported database type: %s", cfg.DBtype)
}
