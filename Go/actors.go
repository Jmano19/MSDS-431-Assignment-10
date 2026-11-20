package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"strings"
)

func findMoviesByActor(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n========== FIND MOVIES BY ACTOR ==========")

	searchTerm := readLine(reader, "Enter an actor's first or last name: ")

	if searchTerm == "" {
		fmt.Println("Search term cannot be empty.")
		return
	}

	actorRows, err := db.Query(`
        SELECT actor_id, first_name, last_name
        FROM actors
        WHERE LOWER(first_name) LIKE LOWER(?)
           OR LOWER(last_name) LIKE LOWER(?)
        ORDER BY last_name, first_name
        LIMIT 20
    `, "%"+searchTerm+"%", "%"+searchTerm+"%")

	if err != nil {
		fmt.Printf("Actor search failed: %v\n", err)
		return
	}
	defer actorRows.Close()

	type Actor struct {
		ID        int
		FirstName string
		LastName  string
	}

	var actors []Actor

	fmt.Println("\n========== ACTOR SEARCH RESULTS ==========")

	for actorRows.Next() {
		var actor Actor

		if err := actorRows.Scan(
			&actor.ID,
			&actor.FirstName,
			&actor.LastName,
		); err != nil {
			fmt.Printf("Error reading actor: %v\n", err)
			return
		}

		actors = append(actors, actor)
		fmt.Printf("%2d. %s %s (ID: %d)\n",
			len(actors),
			actor.FirstName,
			actor.LastName,
			actor.ID,
		)
	}

	if err := actorRows.Err(); err != nil {
		fmt.Printf("Error processing actors: %v\n", err)
		return
	}

	if len(actors) == 0 {
		fmt.Println("No matching actors found.")
		return
	}

	choice := readInt(
		reader,
		"\nSelect an actor (0 to return): ",
		0,
		len(actors),
	)

	if choice == 0 {
		return
	}

	selectedActor := actors[choice-1]

	movieRows, err := db.Query(`
        SELECT DISTINCT m.movie_id, m.name, m.year, m.rank, r.role
        FROM roles r
        JOIN movies m ON r.movie_id = m.movie_id
        WHERE r.actor_id = ?
        ORDER BY m.rank DESC NULLS LAST, m.name
        LIMIT 30
    `, selectedActor.ID)

	if err != nil {
		fmt.Printf("Failed to retrieve actor's movies: %v\n", err)
		return
	}
	defer movieRows.Close()

	fmt.Printf(
		"\n========== MOVIES FEATURING %s ==========\n",
		strings.ToUpper(
			strings.TrimSpace(selectedActor.FirstName+" "+selectedActor.LastName),
		),
	)

	count := 0

	for movieRows.Next() {
		var (
			movieID int
			name    string
			year    sql.NullInt64
			rank    sql.NullFloat64
			role    sql.NullString
		)

		if err := movieRows.Scan(
			&movieID,
			&name,
			&year,
			&rank,
			&role,
		); err != nil {
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

		roleText := "Not listed"
		if role.Valid && strings.TrimSpace(role.String) != "" {
			roleText = role.String
		}

		fmt.Printf(
			"ID: %-8d | Year: %-7s | Rank: %-6s | Movie: %-35s | Role: %s\n",
			movieID,
			yearText,
			rankText,
			name,
			roleText,
		)

		count++
	}

	if err := movieRows.Err(); err != nil {
		fmt.Printf("Error processing movies: %v\n", err)
		return
	}

	if count == 0 {
		fmt.Println("No movies found for this actor.")
	} else {
		fmt.Printf("\nShowing %d movies (maximum 30).\n", count)
	}
}
