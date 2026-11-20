package main

import (
	"bufio"
	"database/sql"
	"strings"
	"testing"
)

func createTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory DuckDB: %v", err)
	}

	db.SetMaxOpenConns(1)

	_, err = db.Exec(`
        CREATE TABLE movies (
            movie_id INTEGER PRIMARY KEY,
            name VARCHAR NOT NULL,
            year INTEGER
        )
    `)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create test movies table: %v", err)
	}

	_, err = db.Exec(`
        INSERT INTO movies (movie_id, name)
        VALUES (1, 'Test Movie'), (2, 'Another Test Movie')
    `)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to insert test movies: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestInitializeCollection(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}

	// Initializing the collection a second time should also succeed.
	if err := initializeCollection(db); err != nil {
		t.Fatalf("Second collection initialization failed: %v", err)
	}

	var count int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM information_schema.tables
        WHERE table_name = 'personal_collection'
    `).Scan(&count)

	if err != nil {
		t.Fatalf("Failed to check collection table: %v", err)
	}

	if count != 1 {
		t.Fatalf("Expected one personal_collection table, got %d", count)
	}
}

func TestCollectionInsertAndRetrieve(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}

	_, err := db.Exec(`
        INSERT INTO personal_collection
            (movie_id, personal_rating, watched)
        VALUES (1, 9.5, TRUE)
    `)
	if err != nil {
		t.Fatalf("Failed to insert collection entry: %v", err)
	}

	var (
		rating  sql.NullFloat64
		watched bool
	)

	err = db.QueryRow(`
        SELECT personal_rating, watched
        FROM personal_collection
        WHERE movie_id = 1
    `).Scan(&rating, &watched)

	if err != nil {
		t.Fatalf("Failed to retrieve collection entry: %v", err)
	}

	if !rating.Valid || rating.Float64 != 9.5 {
		t.Errorf("Expected rating 9.5, got %v", rating)
	}

	if !watched {
		t.Error("Expected watched status to be true")
	}
}

func TestCollectionConstraints(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatalf("Failed to initialize collection: %v", err)
	}

	tests := []struct {
		name  string
		query string
		args  []interface{}
	}{
		{
			name:  "duplicate movie",
			query: "INSERT INTO personal_collection (movie_id) VALUES (?)",
			args:  []interface{}{1},
		},
		{
			name:  "rating above 10",
			query: "INSERT INTO personal_collection (movie_id, personal_rating) VALUES (?, ?)",
			args:  []interface{}{2, 11.0},
		},
		{
			name:  "rating below 0",
			query: "INSERT INTO personal_collection (movie_id, personal_rating) VALUES (?, ?)",
			args:  []interface{}{2, -1.0},
		},
		{
			name:  "nonexistent movie",
			query: "INSERT INTO personal_collection (movie_id) VALUES (?)",
			args:  []interface{}{999},
		},
	}

	// Insert one valid entry so the duplicate test has an existing record.
	_, err := db.Exec(
		"INSERT INTO personal_collection (movie_id) VALUES (?)",
		1,
	)
	if err != nil {
		t.Fatalf("Failed to create initial collection entry: %v", err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := db.Exec(test.query, test.args...)
			if err == nil {
				t.Errorf("Expected an error for %s, but got none", test.name)
			}
		})
	}
}

func TestReadInt(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("invalid\n8\n3\n"))

	result := readInt(reader, "Enter a number: ", 1, 5)

	if result != 3 {
		t.Errorf("Expected 3 after rejecting invalid inputs, got %d", result)
	}
}
