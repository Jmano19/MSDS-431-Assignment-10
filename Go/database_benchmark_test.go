package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
)

func openBenchmarkDB(b *testing.B) *sql.DB {
	b.Helper()

	db, err := sql.Open("duckdb", "../data/database/movies.duckdb")
	if err != nil {
		b.Fatalf("Failed to open benchmark database: %v", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		b.Fatalf("Failed to connect to benchmark database: %v", err)
	}

	b.Cleanup(func() {
		db.Close()
	})

	return db
}

func BenchmarkMovieSearch(b *testing.B) {
	db := openBenchmarkDB(b)

	query := `
        SELECT movie_id, name, year, rank
        FROM movies
        WHERE LOWER(name) LIKE LOWER(?)
        ORDER BY rank DESC NULLS LAST, name
        LIMIT 20
    `

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rows, err := db.Query(query, "%matrix%")
		if err != nil {
			b.Fatal(err)
		}

		for rows.Next() {
			var (
				id   int
				name string
				year sql.NullInt64
				rank sql.NullFloat64
			)

			if err := rows.Scan(&id, &name, &year, &rank); err != nil {
				rows.Close()
				b.Fatal(err)
			}
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			b.Fatal(err)
		}

		rows.Close()
	}
}

func BenchmarkGenreAnalysis(b *testing.B) {
	db := openBenchmarkDB(b)

	query := `
        SELECT
            g.genre,
            COUNT(*) AS movie_count,
            AVG(m.rank) AS average_rank,
            MAX(m.rank) AS maximum_rank,
            MIN(m.rank) AS minimum_rank
        FROM genres g
        JOIN movies m ON g.movie_id = m.movie_id
        GROUP BY g.genre
        ORDER BY average_rank DESC
    `

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rows, err := db.Query(query)
		if err != nil {
			b.Fatal(err)
		}

		for rows.Next() {
			var (
				genre string
				count int
				avg   sql.NullFloat64
				max   sql.NullFloat64
				min   sql.NullFloat64
			)

			if err := rows.Scan(&genre, &count, &avg, &max, &min); err != nil {
				rows.Close()
				b.Fatal(err)
			}
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			b.Fatal(err)
		}

		rows.Close()
	}
}

func BenchmarkActorSearch(b *testing.B) {
	db := openBenchmarkDB(b)

	query := `
        SELECT DISTINCT
            a.actor_id,
            a.first_name,
            a.last_name
        FROM actors a
        JOIN roles r ON a.actor_id = r.actor_id
        WHERE LOWER(a.first_name) LIKE LOWER(?)
           OR LOWER(a.last_name) LIKE LOWER(?)
        ORDER BY a.last_name, a.first_name
        LIMIT 20
    `

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rows, err := db.Query(query, "%tom%", "%tom%")
		if err != nil {
			b.Fatal(err)
		}

		for rows.Next() {
			var (
				id        int
				firstName string
				lastName  string
			)

			if err := rows.Scan(&id, &firstName, &lastName); err != nil {
				rows.Close()
				b.Fatal(err)
			}
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			b.Fatal(err)
		}

		rows.Close()
	}
}

func Example_benchmarkMetrics() {
	fmt.Println("Go benchmarks report execution time, memory, and allocations.")
	// Output: Go benchmarks report execution time, memory, and allocations.
}
