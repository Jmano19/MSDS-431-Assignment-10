package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"strings"
)

func initializeCollection(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS personal_collection (
            movie_id INTEGER PRIMARY KEY,
            date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            personal_rating DOUBLE,
            watched BOOLEAN DEFAULT FALSE,
            FOREIGN KEY (movie_id) REFERENCES movies(movie_id),
            CHECK (personal_rating IS NULL OR
                   (personal_rating >= 0 AND personal_rating <= 10))
        )
    `)

	return err
}

func addToCollection(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n========== ADD MOVIE TO COLLECTION ==========")

	movieID := readInt(reader, "Enter movie ID: ", 1, 999999)

	var title string
	err := db.QueryRow(
		"SELECT name FROM movies WHERE movie_id = ?",
		movieID,
	).Scan(&title)

	if err == sql.ErrNoRows {
		fmt.Println("Movie not found.")
		return
	}
	if err != nil {
		fmt.Printf("Failed to find movie: %v\n", err)
		return
	}

	var exists bool
	err = db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM personal_collection WHERE movie_id = ?
        )
    `, movieID).Scan(&exists)

	if err != nil {
		fmt.Printf("Failed to check collection: %v\n", err)
		return
	}

	if exists {
		fmt.Printf("%s is already in your collection.\n", title)
		return
	}

	ratingInput := readLine(reader, "Personal rating (0-10, or press Enter to skip): ")

	var rating interface{}

	if ratingInput != "" {
		var parsedRating float64
		if _, err := fmt.Sscanf(ratingInput, "%f", &parsedRating); err != nil ||
			parsedRating < 0 || parsedRating > 10 {
			fmt.Println("Invalid rating. Please enter a number from 0 to 10.")
			return
		}
		rating = parsedRating
	}

	watchedInput := strings.ToLower(
		readLine(reader, "Have you watched this movie? (y/n): "),
	)

	var watched bool
	switch watchedInput {
	case "y", "yes":
		watched = true
	case "n", "no":
		watched = false
	default:
		fmt.Println("Please enter y or n.")
		return
	}

	_, err = db.Exec(`
        INSERT INTO personal_collection
            (movie_id, personal_rating, watched)
        VALUES (?, ?, ?)
    `, movieID, rating, watched)

	if err != nil {
		fmt.Printf("Failed to add movie: %v\n", err)
		return
	}

	fmt.Printf("\nSuccessfully added '%s' to your collection!\n", title)
}

func viewCollection(db *sql.DB) {
	fmt.Println("\n========== MY PERSONAL COLLECTION ==========")

	rows, err := db.Query(`
        SELECT
            m.movie_id,
            m.name,
            m.year,
            pc.personal_rating,
            pc.watched,
            pc.date_added
        FROM personal_collection pc
        JOIN movies m ON pc.movie_id = m.movie_id
        ORDER BY pc.date_added DESC, m.name
    `)

	if err != nil {
		fmt.Printf("Failed to retrieve collection: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0

	fmt.Printf(
		"\n%-8s %-35s %-7s %-10s %-10s\n",
		"ID", "TITLE", "YEAR", "MY RATING", "WATCHED",
	)
	fmt.Println("--------------------------------------------------------------------------")

	for rows.Next() {
		var (
			movieID   int
			title     string
			year      sql.NullInt64
			rating    sql.NullFloat64
			watched   bool
			dateAdded sql.NullTime
		)

		if err := rows.Scan(
			&movieID,
			&title,
			&year,
			&rating,
			&watched,
			&dateAdded,
		); err != nil {
			fmt.Printf("Error reading collection: %v\n", err)
			return
		}

		yearText := "Unknown"
		if year.Valid {
			yearText = fmt.Sprintf("%d", year.Int64)
		}

		ratingText := "Not rated"
		if rating.Valid {
			ratingText = fmt.Sprintf("%.1f", rating.Float64)
		}

		watchedText := "No"
		if watched {
			watchedText = "Yes"
		}

		fmt.Printf(
			"%-8d %-35.35s %-7s %-10s %-10s\n",
			movieID, title, yearText, ratingText, watchedText,
		)
		count++
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error processing collection: %v\n", err)
		return
	}

	if count == 0 {
		fmt.Println("Your collection is currently empty.")
	} else {
		fmt.Printf("\nTotal movies in your collection: %d\n", count)
	}
}

func removeFromCollection(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n========== REMOVE MOVIE FROM COLLECTION ==========")

	movieID := readInt(reader, "Enter movie ID to remove: ", 1, 999999)

	result, err := db.Exec(
		"DELETE FROM personal_collection WHERE movie_id = ?",
		movieID,
	)

	if err != nil {
		fmt.Printf("Failed to remove movie: %v\n", err)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Could not confirm removal: %v\n", err)
		return
	}

	if affected == 0 {
		fmt.Println("That movie is not in your collection.")
	} else {
		fmt.Println("Movie successfully removed from your collection.")
	}
}

func manageCollection(db *sql.DB, reader *bufio.Reader) {
	if err := initializeCollection(db); err != nil {
		fmt.Printf("Could not initialize personal collection: %v\n", err)
		return
	}

	for {
		fmt.Println("\n========================================")
		fmt.Println("          PERSONAL COLLECTION")
		fmt.Println("========================================")
		fmt.Println("1. Add Movie")
		fmt.Println("2. View Collection")
		fmt.Println("3. Remove Movie")
		fmt.Println("4. Return to Main Menu")

		choice := readInt(reader, "Select an option (1-4): ", 1, 4)

		switch choice {
		case 1:
			addToCollection(db, reader)
		case 2:
			viewCollection(db)
		case 3:
			removeFromCollection(db, reader)
		case 4:
			return
		}
	}
}
