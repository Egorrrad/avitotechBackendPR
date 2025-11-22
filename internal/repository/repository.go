package repository

import "github.com/Egorrrad/avitotechBackendPR/internal/storage"

type Repository struct {
	Storage storage.DataStorage
}

func NewRepository(storage storage.DataStorage) Repository {
	return Repository{
		Storage: storage,
	}
}
