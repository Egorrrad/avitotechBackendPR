package storage

import (
	"database/sql"
	"log/slog"

	"github.com/Egorrrad/avitotechBackendPR/internal/storage/postgres"
)

type DataStorage interface {
}

func NewDataStorage(dsn string) (DataStorage, *sql.DB) {
	store, err := postgres.OpenDB(dsn)
	if err != nil {
		slog.Error(err.Error())
	}
	return store, store.DB
}
