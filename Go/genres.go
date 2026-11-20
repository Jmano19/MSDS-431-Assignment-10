package main

import (
	"bufio"
	"database/sql"
	"fmt"
)

func exploreGenres(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n========== EXPLORE GENRES ==========")

	rows, err := db.Query(`
        SELECT genre, COUNT(*) AS movie_count
        FROM genres
        GROUP BY genre
        ORDER BY genre
    `)
	if err != nil {
		fmt.Printf("Failed to retrieve genres: %v\n", err)
		return
	}

	var genres []string

	fmt.Println("\nAVAILABLE GENRES")
	fmt.Println("----------------------------------------")

	for rows.Next() {
		var genre string
		var count int

		if err := rows.Scan(&genre, &count); err != nil {
			rows.Close()
			fmt.Printf("Error reading genres: %v\n", err)
			return
		}

		genres = append(genres, genre)
		fmt.Printf("%2d. %-20s %5d movies\n", len(genres), genre, count)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		fmt.Printf("Error processing genres: %v\n", err)
		return
	}
	rows.Close()

	if len(genres) == 0 {
		fmt.Println("No genres found.")
		return
	}

	choice := readInt(reader, "\nSelect a genre (0 to return): ", 0, len(genres))
	if choice == 0 {
		return
	}

	selectedGenre := genres[choice-1]

	movieRows, err := db.Query(`
        SELECT m.movie_id, m.name, m.year, m.rank
        FROM movies m
        JOIN genres g ON m.movie_id = g.movie_id
        WHERE g.genre = ?
        ORDER BY m.rank DESC NULLS LAST, m.name
        LIMIT 20
    `, selectedGenre)
	if err != nil {
		fmt.Printf("Failed to retrieve movies: %v\n", err)
		return
	}
	defer movieRows.Close()

	fmt.Printf("\n========== %s MOVIES ==========\n", selectedGenre)

	count := 0

	for movieRows.Next() {
		var (
			movieID int
			name    string
			year    sql.NullInt64
			rank    sql.NullFloat64
		)

		if err := movieRows.Scan(&movieID, &name, &year, &rank); err != nil {
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
			movieID, yearText, rankText, name,
		)
		count++
	}

	if err := movieRows.Err(); err != nil {
		fmt.Printf("Error processing movies: %v\n", err)
		return
	}

	if count == 0 {
		fmt.Println("No movies found for this genre.")
	} else {
		fmt.Printf("\nShowing %d movies (maximum 20).\n", count)
	}
}
