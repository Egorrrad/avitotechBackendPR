package postgres

import "github.com/Egorrrad/avitotechBackendPR/pkg/postgres"

type PullRequestRepo struct {
	*postgres.Postgres
}

func NewPullRequestRepo(pg *postgres.Postgres) *PullRequestRepo {
	return &PullRequestRepo{pg}
}
