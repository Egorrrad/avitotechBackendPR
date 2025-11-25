package postgres

import (
	"context"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
	"github.com/Egorrrad/avitotechBackendPR/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgtype"
)

type TeamRepo struct {
	*postgres.Postgres
}

func NewTeamRepo(pg *postgres.Postgres) *TeamRepo {
	return &TeamRepo{pg}
}

func (r *TeamRepo) Create(ctx context.Context, team *domain.Team) error {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Insert("teams").
		Columns("name").
		Values(team.TeamName).
		ToSql()

	if err != nil {
		return err
	}

	_, err = q.Exec(ctx, sql, args...)

	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return domain.ErrTeamAlreadyExists
		}
		return err
	}

	return nil
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Select("t.name", "u.user_id", "u.username", "u.is_active").
		From("teams t").
		LeftJoin("team_member tm ON t.id = tm.team_id").
		LeftJoin("users u ON tm.user_id = u.id").
		Where(squirrel.Eq{"t.name": name}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var team *domain.Team

	var userID pgtype.Text
	var username pgtype.Text
	var isActive pgtype.Bool

	teamFound := false

	for rows.Next() {
		var teamName string

		err := rows.Scan(&teamName, &userID, &username, &isActive)
		if err != nil {
			return nil, err
		}

		if team == nil {
			team = &domain.Team{
				TeamName: teamName,
				Members:  make([]domain.TeamMember, 0),
			}
		}

		if userID.Valid {
			team.Members = append(team.Members, domain.TeamMember{
				UserID:   userID.String,
				Username: username.String,
				IsActive: isActive.Bool,
			})
		}
		teamFound = true
	}

	if !teamFound {
		return nil, domain.ErrTeamNotFound
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return team, nil
}

func (r *TeamRepo) Exists(ctx context.Context, name string) (bool, error) {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Select("COUNT(*)").
		From("teams").
		Where(squirrel.Eq{"name": name}).
		ToSql()
	if err != nil {
		return false, err
	}

	var count int
	err = q.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
