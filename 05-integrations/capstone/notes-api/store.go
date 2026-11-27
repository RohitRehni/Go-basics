package main

// store.go: the data layer. Only this file knows SQL.

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound lets the handlers map "no such note" to a 404 without knowing about pgx.
var ErrNotFound = errors.New("note not found")

type Note struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Body      string    `json:"body" db:"body"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type NoteStore struct {
	pool *pgxpool.Pool // safe for concurrent use by many handlers
}

func NewNoteStore(ctx context.Context, url string) (*NoteStore, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 5
	connectCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(connectCtx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, err
	}
	return &NoteStore{pool: pool}, nil
}

func (s *NoteStore) Close() { s.pool.Close() }

func (s *NoteStore) CreateSchema(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		DROP TABLE IF EXISTS go_notes;
		CREATE TABLE go_notes (
			id         SERIAL PRIMARY KEY,
			title      TEXT NOT NULL,
			body       TEXT NOT NULL DEFAULT '',
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	return err
}

func (s *NoteStore) DropSchema(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "DROP TABLE IF EXISTS go_notes")
	return err
}

func (s *NoteStore) Create(ctx context.Context, title, body string) (Note, error) {
	rows, _ := s.pool.Query(ctx,
		"INSERT INTO go_notes (title, body) VALUES ($1, $2) RETURNING id, title, body, updated_at", title, body)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Note])
}

func (s *NoteStore) Get(ctx context.Context, id int) (Note, error) {
	rows, _ := s.pool.Query(ctx, "SELECT id, title, body, updated_at FROM go_notes WHERE id = $1", id)
	n, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Note])
	if errors.Is(err, pgx.ErrNoRows) {
		return Note{}, ErrNotFound
	}
	return n, err
}

func (s *NoteStore) List(ctx context.Context) ([]Note, error) {
	rows, _ := s.pool.Query(ctx, "SELECT id, title, body, updated_at FROM go_notes ORDER BY id")
	return pgx.CollectRows(rows, pgx.RowToStructByName[Note])
}

func (s *NoteStore) Update(ctx context.Context, id int, title, body string) (Note, error) {
	rows, _ := s.pool.Query(ctx,
		"UPDATE go_notes SET title = $2, body = $3, updated_at = now() WHERE id = $1 RETURNING id, title, body, updated_at",
		id, title, body)
	n, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Note])
	if errors.Is(err, pgx.ErrNoRows) {
		return Note{}, ErrNotFound
	}
	return n, err
}

func (s *NoteStore) Delete(ctx context.Context, id int) error {
	tag, err := s.pool.Exec(ctx, "DELETE FROM go_notes WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
