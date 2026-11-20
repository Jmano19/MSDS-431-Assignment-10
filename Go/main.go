package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("       MOVIE DATABASE APPLICATION")
	fmt.Println("========================================")

	db, err := connectDatabase()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	fmt.Println("\nConnected to the movie database successfully!")

	fmt.Println("\n========== DATABASE SUMMARY ==========")

	tables := []string{"movies", "genres", "actors", "roles"}

	for _, table := range tables {
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)

		if err := db.QueryRow(query).Scan(&count); err != nil {
			log.Fatalf("Failed to count %s: %v", table, err)
		}

		fmt.Printf("%-10s %10d records\n", table, count)
	}

	runMenu(db)
}
