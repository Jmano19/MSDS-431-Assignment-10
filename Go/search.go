package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"strings"
)

func searchMovies(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n========== SEARCH MOVIES ==========")

	searchTerm := readLine(reader, "Enter a movie title (or part of a title): ")

	if searchTerm == "" {
		fmt.Println("Search term cannot be empty.")
		return
	}

	query := `
        SELECT movie_id, name, year, rank
        FROM movies
        WHERE LOWER(name) LIKE LOWER(?)
        ORDER BY rank DESC NULLS LAST, name
        LIMIT 20
    `

	rows, err := db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		fmt.Printf("Movie search failed: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0

	fmt.Println("\n========== SEARCH RESULTS ==========")

	for rows.Next() {
		var (
			movieID int
			name    string
			year    sql.NullInt64
			rank    sql.NullFloat64
		)

		if err := rows.Scan(&movieID, &name, &year, &rank); err != nil {
			fmt.Printf("Error reading movie: %v\n", err)
			return
		}

		yearText := "Unknown"
		if year.Valid {
			yearText = fmt.Sprintf("%d", year.Int64)
		}

		rankText := "Not available"
		if rank.Valid {
			rankText = fmt.Sprintf("%.1f", rank.Float64)
		}

		fmt.Printf(
			"ID: %-8d | Year: %-7s | Rank: %-6s | %s\n",
			movieID,
			yearText,
			rankText,
			strings.TrimSpace(name),
		)

		count++
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error processing search results: %v\n", err)
		return
	}

	if count == 0 {
		fmt.Println("No matching movies found.")
	} else {
		fmt.Printf("\nFound %d matching movie(s). Showing up to 20 results.\n", count)
	}
}
