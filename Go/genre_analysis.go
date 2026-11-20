package main

import (
	"database/sql"
	"fmt"
)

func showGenreAnalysis(db *sql.DB) {
	fmt.Println("\n========== GENRE RANK ANALYSIS ==========")

	rows, err := db.Query(`
        SELECT
            g.genre,
            COUNT(*) AS movie_count,
            COUNT(m.rank) AS ranked_movies,
            AVG(m.rank) AS average_rank,
            MAX(m.rank) AS highest_rank,
            MIN(m.rank) AS lowest_rank
        FROM genres g
        JOIN movies m ON g.movie_id = m.movie_id
        GROUP BY g.genre
        ORDER BY average_rank DESC NULLS LAST, g.genre
    `)

	if err != nil {
		fmt.Printf("Failed to retrieve genre analysis: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf(
		"\n%-15s %10s %12s %12s %12s %12s\n",
		"GENRE",
		"MOVIES",
		"RANKED",
		"AVG RANK",
		"HIGHEST",
		"LOWEST",
	)
	fmt.Println("--------------------------------------------------------------------------")

	for rows.Next() {
		var (
			genre       string
			movieCount  int
			rankedCount int
			average     sql.NullFloat64
			highest     sql.NullFloat64
			lowest      sql.NullFloat64
		)

		if err := rows.Scan(
			&genre,
			&movieCount,
			&rankedCount,
			&average,
			&highest,
			&lowest,
		); err != nil {
			fmt.Printf("Error reading genre analysis: %v\n", err)
			return
		}

		averageText := "N/A"
		if average.Valid {
			averageText = fmt.Sprintf("%.2f", average.Float64)
		}

		highestText := "N/A"
		if highest.Valid {
			highestText = fmt.Sprintf("%.2f", highest.Float64)
		}

		lowestText := "N/A"
		if lowest.Valid {
			lowestText = fmt.Sprintf("%.2f", lowest.Float64)
		}

		fmt.Printf(
			"%-15s %10d %12d %12s %12s %12s\n",
			genre,
			movieCount,
			rankedCount,
			averageText,
			highestText,
			lowestText,
		)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error processing genre analysis: %v\n", err)
		return
	}

	fmt.Println("\nAnalysis uses the source dataset's rank field.")
	fmt.Println("Movies with missing ranks are excluded from rank calculations.")
}
