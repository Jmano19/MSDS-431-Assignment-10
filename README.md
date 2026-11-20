# Movie Database Application — MSDS 431 Assignment 10

A movie database application built with **Python, Polars, DuckDB, and Go**. The project uses Polars to ingest and prepare CSV datasets, DuckDB as a relational database, and a Go command-line application to explore movie information, analyze genres, search actors, and manage a personal movie collection.

## Project Overview

The application combines Python-based data engineering with a Go-based database interface.

**Technology stack**
- **Python 3.11** — Data ingestion and database validation
- **Polars** — CSV loading and DataFrame processing
- **DuckDB** — Relational database storage and SQL queries
- **Go 1.27.1** — Interactive command-line application
- **duckdb-go/v2** — Go database driver
- **CGO and GCC** — Required for the Go DuckDB driver

## Project Structure

```text
MSDS-431-Assignment-10/
│
├── Go/
│   ├── main.go
│   ├── database.go
│   ├── menu.go
│   ├── search.go
│   ├── genres.go
│   ├── actors.go
│   ├── top_movies.go
│   ├── genre_analysis.go
│   ├── movie_details.go
│   ├── collection.go
│   ├── database_test.go
│   ├── collection_workflow_test.go
│   └── database_benchmark_test.go
│
├── Python/
│   ├── ingest_data.py
│   └── validate_database.py
│
├── data/
│   ├── raw/
│   │   ├── aan520_w5a_movies2000After.csv
│   │   ├── aan520_w5a_genres2000After.csv
│   │   ├── aan520_w5a_actors2000After.csv
│   │   └── aan520_w5a_roles2000After.csv
│   │
│   └── database/
│       └── movies.duckdb
│
├── results/
│   └── go_benchmarks.txt
│
├── go.mod
├── go.sum
└── README.md
```

The DuckDB database, Python virtual environment, compiled executable, and generated coverage profiles are excluded from version control.

## Database and Data Ingestion

The application uses four CSV datasets:

| Dataset | Records |
|---|---:|
| Movies | 11,502 |
| Genres | 18,548 |
| Actors | 130,002 |
| Roles | 203,027 |
| **Total** | **363,079** |

### Database schema

The DuckDB database contains four relational tables:

**Movies**
- `movie_id` — Primary key
- `name` — Movie title
- `year` — Release year
- `rank` — Source dataset ranking field

**Genres**
- `movie_id` — Foreign key referencing movies
- `genre` — Genre name
- Composite primary key: `(movie_id, genre)`

**Actors**
- `actor_id` — Primary key
- `first_name`
- `last_name`
- `gender`

**Roles**
- `role_id` — Generated unique identifier
- `actor_id` — Foreign key referencing actors
- `movie_id` — Foreign key referencing movies
- `role` — Character or role name

### Ingestion workflow

The Python ingestion script (`Python/ingest_data.py`) performs the following operations:

1. Loads the four source CSV files into Polars DataFrames.
2. Creates the DuckDB database and relational schema.
3. Generates sequential role IDs for the roles dataset.
4. Registers Polars DataFrames with DuckDB.
5. Inserts the data into the corresponding relational tables.
6. Reports the final record counts.

**Important:** The ingestion script deletes and recreates an existing database to ensure a clean rebuild. Running it again will erase any saved personal collection entries. Back up the database before rebuilding if you want to preserve them.

### Database validation

The validation script (`Python/validate_database.py`) checks:
- Expected record counts
- Unique role IDs
- Missing role names
- Orphaned actor references
- Orphaned movie references
- Orphaned genre references

The source data contains 17,386 missing or blank role names. These are retained as source-data conditions rather than treated as ingestion failures.

## Go Application Features

The Go application connects to DuckDB and provides an interactive command-line menu.

| Menu option | Feature |
|---|---|
| 1 | Search Movies |
| 2 | Explore Genres |
| 3 | Find Movies by Actor |
| 4 | Top-Ranked Movies |
| 5 | Genre Rating Analysis |
| 6 | View Movie Details |
| 7 | Manage Personal Collection |
| 8 | Exit |

### Feature descriptions

**Movie Search**
- Search movie titles using case-insensitive substring matching.
- Display matching movie records.

**Genre Exploration**
- Browse available genres.
- View movie counts by genre.
- Explore movies associated with a selected genre.

**Actor Search**
- Search actors by name.
- View movies associated with a selected actor.

**Top-Ranked Movies**
- Retrieve movies ordered by the source dataset's `rank` field.
- Allow the user to specify the number of results.

**Genre Rating Analysis**
- Group movies by genre.
- Calculate aggregate rank statistics, including average, minimum, and maximum values.

**Movie Details**
- Retrieve detailed information for a selected movie.
- Display associated genres and cast information.

**Personal Movie Collection**
- Add movies to a persistent personal collection.
- Assign an optional personal rating from 0 to 10.
- Track whether a movie has been watched.
- View saved collection entries.
- Remove movies from the collection.

The personal collection is stored in a separate `personal_collection` table in DuckDB. It persists between application sessions.

## Installation and Setup

### Prerequisites

- Python 3.11 or compatible version
- Go 1.27.1 or compatible version
- GCC and G++ toolchain for CGO
- Git

### 1. Clone the repository

```bash
git clone https://github.com/Jmano19/MSDS-431-Assignment-10.git
cd MSDS-431-Assignment-10
```

### 2. Set up Python

Create a virtual environment:

```powershell
py -3.11 -m venv .venv
```

Install the required packages:

```powershell
.\.venv\Scripts\python.exe -m pip install duckdb polars
```

### 3. Build the DuckDB database

Run the ingestion script:

```powershell
.\.venv\Scripts\python.exe .\Python\ingest_data.py
```

This creates the database at:

```text
data/database/movies.duckdb
```

### 4. Validate the database

```powershell
.\.venv\Scripts\python.exe .\Python\validate_database.py
```

The validation script should report that all database integrity checks passed.

### 5. Configure CGO for Go

The DuckDB Go driver requires CGO. Ensure GCC and G++ are installed and available on your system PATH.

For the Windows development environment used for this project:

```powershell
$gccPath = "C:\mingw64\bin"
$env:Path = "$gccPath;$env:Path"
$env:CC = "gcc"
$env:CXX = "g++"
$env:CGO_ENABLED = "1"
```

If your compiler is installed in a different location, update `$gccPath` accordingly.

### 6. Run the Go application

From the project root:

```powershell
go run ./Go
```

Alternatively, build and run the executable:

```powershell
go build -o .\Go\movie_app.exe .\Go
.\Go\movie_app.exe
```

## Testing

The project includes Go unit and workflow tests.

Run all tests:

```powershell
go test ./...
```

Run tests with verbose output:

```powershell
go test -v ./Go
```

Run tests with coverage:

```powershell
go test ./Go -coverprofile=results/coverage.out
go tool cover -func results/coverage.out
```

Run static analysis:

```powershell
go vet ./...
```

Build the application:

```powershell
go build -o .\Go\movie_app.exe .\Go
```

The latest test run passed, with **22.3% overall Go statement coverage**. The collection workflow functions have particularly strong coverage:

| Function | Coverage |
|---|---:|
| initializeCollection | 100.0% |
| addToCollection | 91.0% |
| viewCollection | 87.0% |
| removeFromCollection | 84.0% |
| manageCollection | 85.7% |
| readLine | 100.0% |
| readInt | 100.0% |

## Performance Benchmarks

The Go application includes database benchmarks for movie search, genre analysis, and actor search.

Run the benchmarks:

```powershell
go test -run "^$" -bench . -benchmem -benchtime=100x ./...
```

The benchmark results below were collected on:

- **Operating system:** Windows
- **Architecture:** amd64
- **Processor:** AMD Ryzen 7 5800X 8-Core Processor
- **Go version:** 1.27.1

### Results

| Benchmark | Time per operation | Memory per operation | Allocations |
|---|---:|---:|---:|
| Movie Search | 2.807 ms | 2,940 B | 103 |
| Genre Analysis | 4.585 ms | 5,249 B | 284 |
| Actor Search | 24.803 ms | 4,005 B | 234 |

These measurements include Go database-driver overhead, query execution, and result scanning. They represent end-to-end benchmark timings rather than isolated DuckDB engine execution times.

The full benchmark output is available in `results/go_benchmarks.txt`.

## Generative AI Assistance

Generative AI tools were used as development assistants during this assignment. ChatGPT assisted with implementation planning, Go code troubleshooting, test development, interpreting test and benchmark results, and documentation preparation.

All generated suggestions and code changes were reviewed, executed, and tested by the student. The student performed the final integration, database validation, application testing, and benchmark runs.

## GenAI Tools

Generative AI tools were used as development assistants throughout this assignment. OpenAI ChatGPT supported implementation planning, Go and Python code development, DuckDB integration, debugging, test development, benchmark analysis, and project documentation.

Specific areas of assistance included:

- **Database implementation:** Planning the DuckDB schema and helping develop the Python data-ingestion and validation scripts.
- **Go application development:** Reviewing and troubleshooting database connections, SQL queries, application features, and personal collection functionality.
- **Testing and validation:** Developing test cases for database operations and collection workflows, interpreting test results, and identifying issues.
- **Performance analysis:** Supporting the creation and interpretation of Go benchmarks for movie search, genre analysis, and actor search.
- **Troubleshooting:** Helping resolve Windows, CGO, GCC, DuckDB driver, and Go build configuration issues.
- **Documentation:** Assisting with the README, project organization, setup instructions, and benchmark reporting.

All generated suggestions and code were reviewed and tested by the student. The student performed the implementation steps, executed the application and validation scripts, reviewed test and benchmark results, and verified the final functionality.

Additional details about the use of generative AI are documented in [`docs/genai_tools.md`](docs/genai_tools.md) and [`docs/genai_conversation_log.md`](docs/genai_conversation_log.md).

## Future Improvements

Potential enhancements include:
- Expanding automated test coverage for the interactive search and analysis functions.
- Adding more advanced movie filtering and sorting.
- Improving query performance for large actor and role searches.
- Adding export functionality for personal movie collections.
- Introducing a graphical or web-based user interface.

## Repository

GitHub: https://github.com/Jmano19/MSDS-431-Assignment-10

---

**Course:** Northwestern University — MSDS 431  
**Assignment:** Assignment 10  
**Project:** Movie Database Application