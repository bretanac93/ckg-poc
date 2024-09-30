package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func Open() (*sql.DB, func() error, error) {
	connString := "host=localhost port=5432 user=commerce password=commerce dbname=commercedb sslmode=disable"

	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, nil, fmt.Errorf("error while connecting to db: %w", err)
	}

	slog.Info("openned connection to postgres")

	closeFunc := func() error {
		return conn.Close()
	}

	if err := createSchema(conn); err != nil {
		return nil, nil, fmt.Errorf("error while creating schema: %w", err)
	}

	slog.Info("schema created")

	return conn, closeFunc, nil
}

func createSchema(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id uuid primary key,
			created_at timestamp default current_timestamp,
			payload jsonb
		)
	`)
	if err != nil {
		return fmt.Errorf("error while creating orders table: %w", err)
	}

	return nil
}
