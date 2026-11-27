// Topic: postgres-crud | Concept: database | Stage: 05-integrations
// Run:   go run ./05-integrations/database/postgres-crud
// Needs: docker compose -f 05-integrations/compose.yaml up -d --wait
//
// Explore:
//   - pgx is the go-to PostgreSQL driver: pgx.Connect for one connection,
//     pgxpool for a pool shared by many goroutines (see capstone/notes-api)
//   - Always pass a context with a timeout: a hung database should not hang your program
//   - Placeholders $1, $2 send values separately from the SQL: no SQL injection
//   - INSERT ... RETURNING gets the generated id in the same round trip
//   - pgx.CollectRows + pgx.RowToStructByName maps columns to struct fields
//   - pgx.ErrNoRows is returned by QueryRow().Scan when nothing matches
//   - CommandTag.RowsAffected() tells you whether UPDATE / DELETE touched anything
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const hint = "Postgres is not reachable. Start it with: docker compose -f 05-integrations/compose.yaml up -d --wait"

type Book struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
	Year   int    `db:"year"`
}

func databaseURL() string {
	if u := os.Getenv("DATABASE_URL"); u != "" {
		return u
	}
	return "postgres://learning:learning@127.0.0.1:55432/learning"
}

func must(err error, what string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", what, err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Topic: PostgreSQL CRUD with pgx")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectCtx, connectCancel := context.WithTimeout(ctx, 3*time.Second)
	conn, err := pgx.Connect(connectCtx, databaseURL())
	connectCancel()
	if err == nil {
		err = conn.Ping(ctx)
	}
	if err != nil {
		fmt.Println(hint)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var version string
	must(conn.QueryRow(ctx, "SELECT current_setting('server_version')").Scan(&version), "version")
	fmt.Println("connected to PostgreSQL", version)

	// --- Schema (dropped again at the end so reruns start clean) ---
	_, err = conn.Exec(ctx, `
		DROP TABLE IF EXISTS go_books;
		CREATE TABLE go_books (
			id     SERIAL PRIMARY KEY,
			title  TEXT NOT NULL,
			author TEXT NOT NULL,
			year   INT  NOT NULL
		)`)
	must(err, "create table")
	defer func() {
		_, _ = conn.Exec(context.Background(), "DROP TABLE IF EXISTS go_books")
		fmt.Println("cleanup: dropped go_books")
	}()

	// --- Create ---
	var id int
	err = conn.QueryRow(ctx,
		"INSERT INTO go_books (title, author, year) VALUES ($1, $2, $3) RETURNING id",
		"The Go Programming Language", "Donovan & Kernighan", 2015).Scan(&id)
	must(err, "insert")
	fmt.Println("CREATE: inserted book with id", id)

	batch := &pgx.Batch{} // several statements in one round trip
	for _, b := range []Book{
		{Title: "Programming Perl", Author: "Wall", Year: 1991},
		{Title: "Fluent Python", Author: "Ramalho", Year: 2015},
		{Title: "The Rust Programming Language", Author: "Klabnik & Nichols", Year: 2018},
	} {
		batch.Queue("INSERT INTO go_books (title, author, year) VALUES ($1, $2, $3)", b.Title, b.Author, b.Year)
	}
	must(conn.SendBatch(ctx, batch).Close(), "batch insert")
	fmt.Println("CREATE: batch-inserted 3 more books")

	// --- Read many ---
	rows, err := conn.Query(ctx, "SELECT id, title, author, year FROM go_books WHERE year >= $1 ORDER BY year, id", 2000)
	must(err, "query")
	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	must(err, "collect rows")
	fmt.Println("READ: books from 2000 onwards:")
	for _, b := range books {
		fmt.Printf("  #%d %-32s %-20s %d\n", b.ID, b.Title, b.Author, b.Year)
	}

	// --- Read one, and the "not found" case ---
	var title string
	must(conn.QueryRow(ctx, "SELECT title FROM go_books WHERE id = $1", id).Scan(&title), "read one")
	fmt.Println("READ: book", id, "is", title)
	err = conn.QueryRow(ctx, "SELECT title FROM go_books WHERE id = $1", 9999).Scan(&title)
	fmt.Println("READ: book 9999 -> pgx.ErrNoRows:", errors.Is(err, pgx.ErrNoRows))

	// --- Placeholders keep hostile input as plain data ---
	hostile := "x'; DROP TABLE go_books; --"
	var count int
	must(conn.QueryRow(ctx, "SELECT count(*) FROM go_books WHERE title = $1", hostile).Scan(&count), "safe query")
	fmt.Println("SAFE: searching for a SQL-injection string matched", count, "rows; the table is still here")

	// --- Update ---
	tag, err := conn.Exec(ctx, "UPDATE go_books SET year = year + 1 WHERE author LIKE $1", "%Ramalho%")
	must(err, "update")
	fmt.Println("UPDATE: rows affected:", tag.RowsAffected())

	// --- Delete ---
	tag, err = conn.Exec(ctx, "DELETE FROM go_books WHERE year < $1", 2000)
	must(err, "delete")
	fmt.Println("DELETE: rows affected:", tag.RowsAffected())

	must(conn.QueryRow(ctx, "SELECT count(*) FROM go_books").Scan(&count), "count")
	fmt.Println("books left:", count)
	if count != 3 {
		fmt.Println("MISMATCH: expected 3 books")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
