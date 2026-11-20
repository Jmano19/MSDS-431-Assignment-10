
from pathlib import Path

import duckdb
import polars as pl


# Project paths
PROJECT_ROOT = Path(__file__).resolve().parent.parent
RAW_DATA_DIR = PROJECT_ROOT / "data" / "raw"
DATABASE_DIR = PROJECT_ROOT / "data" / "database"
DATABASE_PATH = DATABASE_DIR / "movies.duckdb"

# Input files
CSV_FILES = {
    "movies": "aan520_w5a_movies2000After.csv",
    "genres": "aan520_w5a_genres2000After.csv",
    "actors": "aan520_w5a_actors2000After.csv",
    "roles": "aan520_w5a_roles2000After.csv",
}


def load_csv_files() -> dict[str, pl.DataFrame]:
    """Load the four source CSV files into Polars DataFrames."""
    dataframes = {}

    for table_name, filename in CSV_FILES.items():
        file_path = RAW_DATA_DIR / filename

        if not file_path.exists():
            raise FileNotFoundError(f"Missing input file: {file_path}")

        df = pl.read_csv(file_path, infer_schema_length=10000)
        dataframes[table_name] = df

        print(f"Loaded {table_name}: {df.height:,} rows")

    return dataframes


def create_schema(connection: duckdb.DuckDBPyConnection) -> None:
    """Create the relational database schema."""

    connection.execute("""
        CREATE TABLE movies (
            movie_id INTEGER PRIMARY KEY,
            name VARCHAR NOT NULL,
            year INTEGER,
            rank DOUBLE
        )
    """)

    connection.execute("""
        CREATE TABLE genres (
            movie_id INTEGER NOT NULL,
            genre VARCHAR NOT NULL,
            PRIMARY KEY (movie_id, genre),
            FOREIGN KEY (movie_id) REFERENCES movies(movie_id)
        )
    """)

    connection.execute("""
        CREATE TABLE actors (
            actor_id INTEGER PRIMARY KEY,
            first_name VARCHAR,
            last_name VARCHAR,
            gender VARCHAR
        )
    """)

    connection.execute("""
        CREATE TABLE roles (
            role_id INTEGER PRIMARY KEY,
            actor_id INTEGER NOT NULL,
            movie_id INTEGER NOT NULL,
            role VARCHAR,
            FOREIGN KEY (actor_id) REFERENCES actors(actor_id),
            FOREIGN KEY (movie_id) REFERENCES movies(movie_id)
        )
    """)



def insert_data(
    connection: duckdb.DuckDBPyConnection,
    dataframes: dict[str, pl.DataFrame],
) -> None:
    """Insert Polars DataFrames into DuckDB tables."""

    for table_name, dataframe in dataframes.items():
        if table_name == "roles":
            dataframe = dataframe.with_row_index("role_id", offset=1)

        connection.register("source_dataframe", dataframe)

        columns = ", ".join(dataframe.columns)

        connection.execute(
            f"INSERT INTO {table_name} ({columns}) "
            f"SELECT {columns} FROM source_dataframe"
        )

        connection.unregister("source_dataframe")

        print(f"Inserted {dataframe.height:,} rows into {table_name}")



def validate_database(connection: duckdb.DuckDBPyConnection) -> None:
    """Check the final database record counts."""

    print("\nDatabase validation:")

    for table_name in CSV_FILES:
        count = connection.execute(
            f"SELECT COUNT(*) FROM {table_name}"
        ).fetchone()[0]

        print(f"{table_name}: {count:,} rows")


def main() -> None:
    DATABASE_DIR.mkdir(parents=True, exist_ok=True)

    print("Loading CSV files with Polars...")
    dataframes = load_csv_files()

    # Remove an existing database to allow a clean rebuild.
    if DATABASE_PATH.exists():
        DATABASE_PATH.unlink()

    print("\nCreating DuckDB database...")
    connection = duckdb.connect(str(DATABASE_PATH))

    try:
        create_schema(connection)
        insert_data(connection, dataframes)
        validate_database(connection)
    finally:
        connection.close()

    print(f"\nDatabase created successfully: {DATABASE_PATH}")


if __name__ == "__main__":
    main()
