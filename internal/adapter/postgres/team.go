package postgres

import (
	"github.com/Egorrrad/avitotechBackendPR/pkg/postgres"
)

type TeamRepo struct {
	*postgres.Postgres
}

func NewTeamRepo(pg *postgres.Postgres) *TeamRepo {
	return &TeamRepo{pg}
}
