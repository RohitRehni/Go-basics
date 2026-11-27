// Topic: transactions | Concept: database | Stage: 05-integrations
// Run:   go run ./05-integrations/database/transactions
// Needs: docker compose -f 05-integrations/compose.yaml up -d --wait
//
// Explore:
//   - A transaction groups statements: all of them commit, or none do (atomicity)
//   - Pattern: tx := Begin; defer tx.Rollback(ctx); ...; tx.Commit(ctx)
//     (Rollback after a successful Commit is a harmless no-op)
//   - pgx.BeginFunc commits if your function returns nil and rolls back on error
//   - CHECK constraints let the database refuse bad data even if the code has a bug
//   - SELECT ... FOR UPDATE locks rows you are about to change (stops lost updates)
//   - Isolation levels: READ COMMITTED (default), REPEATABLE READ, SERIALIZABLE
//     (serializable may fail with a serialization error: retry the whole transaction)
//   - A context deadline cancels a slow query on the server, not just in Go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const hint = "Postgres is not reachable. Start it with: docker compose -f 05-integrations/compose.yaml up -d --wait"

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

// transfer moves money between two accounts inside one transaction.
func transfer(ctx context.Context, conn *pgx.Conn, from, to string, amount int) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // undoes everything unless Commit ran first

	// Lock both rows in a fixed order so two opposite transfers cannot deadlock.
	_, err = tx.Exec(ctx, "SELECT 1 FROM go_accounts WHERE owner IN ($1, $2) ORDER BY owner FOR UPDATE", from, to)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "UPDATE go_accounts SET balance = balance + $1 WHERE owner = $2", amount, to); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "UPDATE go_accounts SET balance = balance - $1 WHERE owner = $2", amount, from); err != nil {
		return err // e.g. the CHECK (balance >= 0) constraint fired
	}
	return tx.Commit(ctx)
}

func balances(ctx context.Context, conn *pgx.Conn) map[string]int {
	rows, err := conn.Query(ctx, "SELECT owner, balance FROM go_accounts")
	must(err, "balances")
	out := map[string]int{}
	var owner string
	var bal int
	_, err = pgx.ForEachRow(rows, []any{&owner, &bal}, func() error {
		out[owner] = bal
		return nil
	})
	must(err, "scan balances")
	return out
}

func main() {
	fmt.Println("Topic: transactions with pgx")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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

	_, err = conn.Exec(ctx, `
		DROP TABLE IF EXISTS go_accounts;
		CREATE TABLE go_accounts (
			owner   TEXT PRIMARY KEY,
			balance INT NOT NULL CHECK (balance >= 0)
		);
		INSERT INTO go_accounts VALUES ('alice', 100), ('bob', 50)`)
	must(err, "setup")
	defer func() { // conn may have been replaced by the reconnect below
		_, err := conn.Exec(context.Background(), "DROP TABLE IF EXISTS go_accounts")
		conn.Close(context.Background())
		fmt.Println("cleanup: dropped go_accounts, error:", err)
	}()
	failed := false
	fmt.Println("start:", balances(ctx, conn))

	// 1. A transfer that succeeds: both updates commit together.
	err = transfer(ctx, conn, "alice", "bob", 30)
	fmt.Println("transfer alice -> bob 30:", err, "|", balances(ctx, conn))

	// 2. A transfer that breaks the CHECK constraint: bob was already credited
	//    inside the transaction, but the rollback undoes it.
	err = transfer(ctx, conn, "alice", "bob", 500)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Printf("transfer alice -> bob 500 failed: SQLSTATE %s (%s)\n", pgErr.Code, pgErr.ConstraintName)
	}
	after := balances(ctx, conn)
	fmt.Println("after the failed transfer (unchanged):", after)
	if after["alice"] != 70 || after["bob"] != 80 {
		failed = true
	}

	// 3. pgx.BeginFunc: commit on nil, roll back on error.
	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "INSERT INTO go_accounts VALUES ('carol', 10)"); err != nil {
			return err
		}
		return errors.New("changed my mind")
	})
	var carol int
	missing := conn.QueryRow(ctx, "SELECT balance FROM go_accounts WHERE owner = 'carol'").Scan(&carol)
	fmt.Println("BeginFunc returned:", err, "| carol exists:", !errors.Is(missing, pgx.ErrNoRows))

	// 4. Isolation level for a read-only report.
	err = pgx.BeginTxFunc(ctx, conn, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly},
		func(tx pgx.Tx) error {
			var total int
			if err := tx.QueryRow(ctx, "SELECT sum(balance) FROM go_accounts").Scan(&total); err != nil {
				return err
			}
			fmt.Println("REPEATABLE READ report, total money in the bank:", total)
			if total != 150 {
				failed = true
			}
			return nil
		})
	must(err, "report")

	// A write inside a READ ONLY transaction is refused. The error also aborts
	// that transaction, so BeginTxFunc rolls it back and returns the error.
	err = pgx.BeginTxFunc(ctx, conn, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "UPDATE go_accounts SET balance = 0")
		return err
	})
	if errors.As(err, &pgErr) {
		fmt.Printf("write inside a READ ONLY transaction refused: SQLSTATE %s\n", pgErr.Code)
	} else {
		failed = true
	}

	// 5. Context deadline cancels a slow query.
	slowCtx, slowCancel := context.WithTimeout(ctx, 300*time.Millisecond)
	start := time.Now()
	_, err = conn.Exec(slowCtx, "SELECT pg_sleep(3)")
	slowCancel()
	fmt.Printf("pg_sleep(3) with a 300ms deadline stopped after %v: %v\n",
		time.Since(start).Round(10*time.Millisecond), errors.Is(err, context.DeadlineExceeded) || err != nil)
	if err == nil {
		failed = true
	}
	// A cancelled query leaves this connection unusable; reconnect for cleanup.
	conn.Close(context.Background())
	conn, err = pgx.Connect(ctx, databaseURL())
	must(err, "reconnect")

	if failed {
		fmt.Println("MISMATCH in balances")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
