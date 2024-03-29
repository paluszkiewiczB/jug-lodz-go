package data

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"time"
	"todo/internal/todo"
)

type SQLiteTaskStorage struct {
	db *sql.DB
}

func NewSQLiteTaskStorage(file string) (*SQLiteTaskStorage, error) {
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite database from file: %s, %w", file, err)
	}
	return &SQLiteTaskStorage{db: db}, nil
}

func (s *SQLiteTaskStorage) Upsert(ctx context.Context, t todo.Task) (stored todo.Task, err error) {
	wDead := fromDeadline(t.Deadline)
	rows, err := s.db.QueryContext(ctx, `
			INSERT INTO tasks (id, title, deadline, done)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(id)
			DO UPDATE SET title = ?, deadline = ?, done = ?
			RETURNING id, title, done, deadline
		`,
		t.ID, t.Title, wDead, t.Done,
		t.Title, wDead, t.Done,
	)
	if err != nil {
		return todo.Task{}, fmt.Errorf("upserting task: %v, %w", t, err)
	}

	if !rows.Next() {
		return todo.Task{}, fmt.Errorf("expected to receive upserted task but no rows were returned")
	}

	ret, err := scanTask(rows)
	if err != nil {
		return todo.Task{}, fmt.Errorf("scanning upserted task, %w", err)
	}

	return ret, nil
}

func scanTask(rows *sql.Rows) (todo.Task, error) {
	ret := todo.Task{}
	retDead := sql.NullString{}
	err := rows.Scan(&ret.ID, &ret.Title, &ret.Done, &retDead)
	if err != nil {
		return todo.Task{}, fmt.Errorf("scanning task from row, %w", err)
	}

	ret.Deadline = toDeadline(retDead)
	return ret, nil
}

func toDeadline(d sql.NullString) *time.Time {
	if !d.Valid {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339, d.String)
	if err != nil {
		panic(fmt.Sprintf("failed to parse the duration from string: %s", d.String))
	}

	return &parsed
}

func (s *SQLiteTaskStorage) List(ctx context.Context, filter *todo.TaskFilter) ([]todo.Task, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, done, deadline
		FROM tasks
		WHERE id like ?
	`,
		id(filter),
	)

	if err != nil {
		return nil, fmt.Errorf("listing tasks with filter: %v, %w", filter, err)
	}

	out := make([]todo.Task, 0)
	for rows.Next() {
		if rows.Err() != nil {
			return out, fmt.Errorf("iterating over results of listing tasks with filter: %v, %w", filter, err)
		}

		t, err := scanTask(rows)
		if err != nil {
			return out, fmt.Errorf("scanning task from row, %w", err)
		}

		out = append(out, t)
	}

	return out, nil
}

func id(f *todo.TaskFilter) string {
	if f == nil {
		return "%"
	}

	if f.ID == nil {
		return "%"
	}
	return string(*f.ID)
}

func (s *SQLiteTaskStorage) Initialize() error {
	// create table if not exists
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL,
			deadline TEXT
		)
	`)
	return err
}

func fromDeadline(d *time.Time) sql.NullString {
	if d == nil {
		return sql.NullString{}
	}

	return sql.NullString{Valid: true, String: d.Format(time.RFC3339)}
}
