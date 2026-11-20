package main

import (
	"bufio"
	"database/sql"
	"strings"
	"testing"
)

func TestAddToCollection(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(strings.NewReader("1\n9.5\ny\n"))
	addToCollection(db, reader)

	var rating sql.NullFloat64
	var watched bool

	err := db.QueryRow(`
        SELECT personal_rating, watched
        FROM personal_collection
        WHERE movie_id = 1
    `).Scan(&rating, &watched)

	if err != nil {
		t.Fatalf("Expected movie to be added: %v", err)
	}

	if !rating.Valid || rating.Float64 != 9.5 {
		t.Errorf("Expected rating 9.5, got %v", rating)
	}

	if !watched {
		t.Error("Expected watched status to be true")
	}
}

func TestAddToCollectionWithoutRating(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(strings.NewReader("2\n\nn\n"))
	addToCollection(db, reader)

	var rating sql.NullFloat64
	var watched bool

	err := db.QueryRow(`
        SELECT personal_rating, watched
        FROM personal_collection
        WHERE movie_id = 2
    `).Scan(&rating, &watched)

	if err != nil {
		t.Fatalf("Expected movie to be added: %v", err)
	}

	if rating.Valid {
		t.Errorf("Expected no personal rating, got %v", rating.Float64)
	}

	if watched {
		t.Error("Expected watched status to be false")
	}
}

func TestAddToCollectionRejectsInvalidEntries(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "nonexistent movie",
			input: "999\n",
		},
		{
			name:  "rating above 10",
			input: "1\n11\n",
		},
		{
			name:  "negative rating",
			input: "1\n-1\n",
		},
		{
			name:  "invalid watched response",
			input: "1\n8\nmaybe\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := createTestDatabase(t)

			if err := initializeCollection(db); err != nil {
				t.Fatal(err)
			}

			reader := bufio.NewReader(strings.NewReader(test.input))
			addToCollection(db, reader)

			var count int
			err := db.QueryRow(
				"SELECT COUNT(*) FROM personal_collection",
			).Scan(&count)

			if err != nil {
				t.Fatal(err)
			}

			if count != 0 {
				t.Errorf("Expected no collection entries, got %d", count)
			}
		})
	}
}

func TestAddToCollectionRejectsDuplicate(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatal(err)
	}

	_, err := db.Exec(
		"INSERT INTO personal_collection (movie_id) VALUES (1)",
	)
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(strings.NewReader("1\n"))
	addToCollection(db, reader)

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM personal_collection WHERE movie_id = 1",
	).Scan(&count)

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Errorf("Expected one entry, got %d", count)
	}
}

func TestViewCollection(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatal(err)
	}

	_, err := db.Exec(`
        INSERT INTO personal_collection
            (movie_id, personal_rating, watched)
        VALUES (1, 9.5, TRUE)
    `)
	if err != nil {
		t.Fatal(err)
	}

	// Should successfully retrieve and display the populated collection.
	viewCollection(db)
}

func TestViewEmptyCollection(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatal(err)
	}

	// Should handle an empty collection without errors.
	viewCollection(db)
}

func TestRemoveFromCollection(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatal(err)
	}

	_, err := db.Exec(
		"INSERT INTO personal_collection (movie_id) VALUES (1)",
	)
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(strings.NewReader("1\n"))
	removeFromCollection(db, reader)

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM personal_collection WHERE movie_id = 1",
	).Scan(&count)

	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Errorf("Expected movie to be removed, got %d entries", count)
	}
}

func TestRemoveNonexistentCollectionEntry(t *testing.T) {
	db := createTestDatabase(t)

	if err := initializeCollection(db); err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(strings.NewReader("1\n"))
	removeFromCollection(db, reader)

	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM personal_collection",
	).Scan(&count)

	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Errorf("Expected empty collection, got %d entries", count)
	}
}

func TestManageCollectionExit(t *testing.T) {
	db := createTestDatabase(t)

	reader := bufio.NewReader(strings.NewReader("4\n"))
	manageCollection(db, reader)

	var count int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM information_schema.tables
        WHERE table_name = 'personal_collection'
    `).Scan(&count)

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Error("Expected manageCollection to initialize the collection table")
	}
}
