
from pathlib import Path

import duckdb


PROJECT_ROOT = Path(__file__).resolve().parent.parent
DATABASE_PATH = PROJECT_ROOT / "data" / "database" / "movies.duckdb"


def main():
    connection = duckdb.connect(str(DATABASE_PATH), read_only=True)

    try:
        print("========== TABLE COUNTS ==========")

        for table in ["movies", "genres", "actors", "roles"]:
            count = connection.execute(
                f"SELECT COUNT(*) FROM {table}"
            ).fetchone()[0]

            print(f"{table}: {count:,}")

        print("\n========== ROLE VALIDATION ==========")

        missing_roles = connection.execute("""
            SELECT COUNT(*)
            FROM roles
            WHERE role IS NULL OR TRIM(role) = ''
        """).fetchone()[0]

        unique_role_ids = connection.execute("""
            SELECT COUNT(DISTINCT role_id)
            FROM roles
        """).fetchone()[0]

        total_roles = connection.execute("""
            SELECT COUNT(*)
            FROM roles
        """).fetchone()[0]

        print(f"Missing role names: {missing_roles:,}")
        print(f"Unique role IDs: {unique_role_ids:,}")
        print(f"Total role records: {total_roles:,}")

        print("\n========== FOREIGN KEY VALIDATION ==========")

        orphan_actors = connection.execute("""
            SELECT COUNT(*)
            FROM roles r
            LEFT JOIN actors a ON r.actor_id = a.actor_id
            WHERE a.actor_id IS NULL
        """).fetchone()[0]

        orphan_movies = connection.execute("""
            SELECT COUNT(*)
            FROM roles r
            LEFT JOIN movies m ON r.movie_id = m.movie_id
            WHERE m.movie_id IS NULL
        """).fetchone()[0]

        orphan_genres = connection.execute("""
            SELECT COUNT(*)
            FROM genres g
            LEFT JOIN movies m ON g.movie_id = m.movie_id
            WHERE m.movie_id IS NULL
        """).fetchone()[0]

        print(f"Roles with missing actors: {orphan_actors}")
        print(f"Roles with missing movies: {orphan_movies}")
        print(f"Genres with missing movies: {orphan_genres}")

        print("\n========== SAMPLE JOIN ==========")

        results = connection.execute("""
            SELECT
                m.name,
                m.year,
                m.rank,
                g.genre
            FROM movies m
            JOIN genres g ON m.movie_id = g.movie_id
            ORDER BY m.rank DESC NULLS LAST
            LIMIT 10
        """).fetchall()

        for name, year, rank, genre in results:
            print(f"{name} | {year} | {rank} | {genre}")

        print("\n========== VALIDATION SUMMARY ==========")

        checks = {
            "Expected role count": total_roles == 203027,
            "All role IDs unique": unique_role_ids == total_roles,
            "Expected missing role names": missing_roles == 17386,
            "No orphaned actors": orphan_actors == 0,
            "No orphaned movies": orphan_movies == 0,
            "No orphaned genres": orphan_genres == 0,
        }

        for description, passed in checks.items():
            print(f"{'PASS' if passed else 'FAIL'}: {description}")

        if not all(checks.values()):
            raise RuntimeError("One or more database validation checks failed.")

        print("\nAll database validation checks passed!")

    finally:
        connection.close()


if __name__ == "__main__":
    main()
