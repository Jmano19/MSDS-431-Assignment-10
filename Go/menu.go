package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt(reader *bufio.Reader, prompt string, min int, max int) int {
	for {
		input := readLine(reader, prompt)
		value, err := strconv.Atoi(input)

		if err == nil && value >= min && value <= max {
			return value
		}

		fmt.Printf("Please enter a number between %d and %d.\n", min, max)
	}
}

func runMenu(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n========================================")
		fmt.Println("             MOVIE DATABASE")
		fmt.Println("========================================")
		fmt.Println("1. Search Movies")
		fmt.Println("2. Explore Genres")
		fmt.Println("3. Find Movies by Actor")
		fmt.Println("4. Top-Ranked Movies")
		fmt.Println("5. Genre Rating Analysis")
		fmt.Println("6. View Movie Details")
		fmt.Println("7. Manage Personal Collection")
		fmt.Println("8. Exit")
		fmt.Println("----------------------------------------")

		choice := readInt(reader, "Select an option (1-8): ", 1, 8)

		switch choice {
		case 1:
			searchMovies(db, reader)
		case 2:
			exploreGenres(db, reader)
		case 3:
			findMoviesByActor(db, reader)
		case 4:
			showTopMovies(db, reader)
		case 5:
			showGenreAnalysis(db)
		case 6:
			showMovieDetails(db, reader)
		case 7:
			manageCollection(db, reader)
		case 8:
			fmt.Println("\nThank you for using the Movie Database Application!")
			return
		}
	}
}
