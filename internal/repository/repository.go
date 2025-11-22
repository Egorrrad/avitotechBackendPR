package repository

import (
	"database/sql"
	"log/slog"

	"github.com/Egorrrad/avitotechBackendPR/internal/repository/postgres"
)

type DataStorage interface {
	PullRequestStore
	TeamStore
	UserStore
}

type PullRequestStore interface {
}

type TeamStore interface {
}

type UserStore interface {
}

func NewDataStorage(dsn string) (DataStorage, *sql.DB) {
	store, err := postgres.OpenDB(dsn)
	if err != nil {
		slog.Error(err.Error())
	}
	return store, store.DB
}

type Repository struct {
	Storage DataStorage
}

func NewRepository(storage DataStorage) Repository {
	return Repository{
		Storage: storage,
	}
}
