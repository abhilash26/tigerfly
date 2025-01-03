package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	log.Println("Connected to SQLite database successfully")
	return &DB{conn: conn}, nil
}

func (d *DB) Close() error {
	if err := d.conn.Close(); err != nil {
		return fmt.Errorf("Failed to close database: %w", err)
	}
	log.Println("Database connection closed")
	return nil
}

func (d *DB) Execute(query string, args ...interface{}) error {
	_, err := d.conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("Query execution failed: %w", err)
	}
	log.Println("Query executed successfully")
	return nil
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Query failed: %w", err)
	}
	return rows, nil
}
