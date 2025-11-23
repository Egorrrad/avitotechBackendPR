package postgres

import "github.com/Egorrrad/avitotechBackendPR/pkg/postgres"

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}
