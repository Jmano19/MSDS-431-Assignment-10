package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"strings"
)

func showMovieDetails(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n========== VIEW MOVIE DETAILS ==========")

	movieID := readInt(reader, "Enter a movie ID (1-999999): ", 1, 999999)

	var (
		name string
		year sql.NullInt64
		rank sql.NullFloat64
	)

	err := db.QueryRow(`
        SELECT name, year, rank
        FROM movies
        WHERE movie_id = ?
    `, movieID).Scan(&name, &year, &rank)

	if err == sql.ErrNoRows {
		fmt.Printf("No movie found with ID %d.\n", movieID)
		return
	}

	if err != nil {
		fmt.Printf("Failed to retrieve movie: %v\n", err)
		return
	}

	fmt.Println("\n========== MOVIE INFORMATION ==========")
	fmt.Printf("Movie ID: %d\n", movieID)
	fmt.Printf("Title:    %s\n", name)

	if year.Valid {
		fmt.Printf("Year:     %d\n", year.Int64)
	} else {
		fmt.Println("Year:     Unknown")
	}

	if rank.Valid {
		fmt.Printf("Rank:     %.2f\n", rank.Float64)
	} else {
		fmt.Println("Rank:     Not available")
	}

	genreRows, err := db.Query(`
        SELECT genre
        FROM genres
        WHERE movie_id = ?
        ORDER BY genre
    `, movieID)

	if err != nil {
		fmt.Printf("Failed to retrieve genres: %v\n", err)
		return
	}

	var genres []string

	for genreRows.Next() {
		var genre string
		if err := genreRows.Scan(&genre); err != nil {
			genreRows.Close()
			fmt.Printf("Error reading genre: %v\n", err)
			return
		}
		genres = append(genres, genre)
	}

	if err := genreRows.Err(); err != nil {
		genreRows.Close()
		fmt.Printf("Error processing genres: %v\n", err)
		return
	}
	genreRows.Close()

	fmt.Println("\nGENRES")
	if len(genres) == 0 {
		fmt.Println("No genres listed.")
	} else {
		fmt.Println(strings.Join(genres, ", "))
	}

	actorRows, err := db.Query(`
        SELECT a.first_name, a.last_name, r.role
        FROM roles r
        JOIN actors a ON r.actor_id = a.actor_id
        WHERE r.movie_id = ?
        ORDER BY a.last_name, a.first_name
        LIMIT 30
    `, movieID)

	if err != nil {
		fmt.Printf("Failed to retrieve cast: %v\n", err)
		return
	}
	defer actorRows.Close()

	fmt.Println("\nCAST (UP TO 30 ACTORS)")
	fmt.Println("----------------------------------------")

	actorCount := 0

	for actorRows.Next() {
		var (
			firstName string
			lastName  string
			role      sql.NullString
		)

		if err := actorRows.Scan(&firstName, &lastName, &role); err != nil {
			fmt.Printf("Error reading cast member: %v\n", err)
			return
		}

		roleText := "Not listed"
		if role.Valid && strings.TrimSpace(role.String) != "" {
			roleText = role.String
		}

		fmt.Printf("%-30s | %s\n", firstName+" "+lastName, roleText)
		actorCount++
	}

	if err := actorRows.Err(); err != nil {
		fmt.Printf("Error processing cast: %v\n", err)
		return
	}

	if actorCount == 0 {
		fmt.Println("No cast information available.")
	}
}
