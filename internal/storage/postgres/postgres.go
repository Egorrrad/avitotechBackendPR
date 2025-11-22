package postgres

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Storage struct {
	DB *sql.DB
}

func OpenDB(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	slog.Info("Connection to the database successfully!")
	store := &Storage{
		DB: db,
	}

	return store, nil
}

func CloseDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return err
	}

	slog.Info("Connection to the database closed successfully!")
	return nil
}
