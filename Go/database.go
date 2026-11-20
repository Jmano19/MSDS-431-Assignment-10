package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/duckdb/duckdb-go/v2"
)

func connectDatabase() (*sql.DB, error) {
	databasePath := filepath.Join("data", "database", "movies.duckdb")

	if _, err := os.Stat(databasePath); err != nil {
		return nil, fmt.Errorf("database file not found at %s: %w", databasePath, err)
	}

	db, err := sql.Open("duckdb", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DuckDB database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to DuckDB: %w", err)
	}

	return db, nil
}
