package main

import (
	"bufio"
	"database/sql"
	"fmt"
)

func showTopMovies(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n========== TOP-RANKED MOVIES ==========")

	limit := readInt(reader, "How many movies would you like to see (5-50)? ", 5, 50)

	rows, err := db.Query(`
        SELECT movie_id, name, year, rank
        FROM movies
        WHERE rank IS NOT NULL
        ORDER BY rank DESC, name
        LIMIT ?
    `, limit)

	if err != nil {
		fmt.Printf("Failed to retrieve top-ranked movies: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("\nTOP %d MOVIES BY SOURCE RANK\n", limit)
	fmt.Println("---------------------------------------------------------------")

	count := 0

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
			"%2d. ID: %-8d | Year: %-7s | Rank: %-6s | %s\n",
			count+1,
			movieID,
			yearText,
			rankText,
			name,
		)

		count++
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error processing movies: %v\n", err)
		return
	}

	if count == 0 {
		fmt.Println("No ranked movies found.")
	} else {
		fmt.Printf("\nDisplayed %d movies.\n", count)
	}
}
