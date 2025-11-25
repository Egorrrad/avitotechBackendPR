package postgres

import (
	"context"
	"errors"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
	"github.com/Egorrrad/avitotechBackendPR/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

type teamInternalID struct {
	ID   int
	Name string
}

func (r *UserRepo) resolveTeamNameToInternalID(ctx context.Context, teamName string) (*teamInternalID, error) {
	q := r.GetQueryer(ctx)
	var t teamInternalID

	sql, args, err := r.Builder.
		Select("id", "name").
		From("teams").
		Where(squirrel.Eq{"name": teamName}).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = q.QueryRow(ctx, sql, args...).Scan(&t.ID, &t.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTeamNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *UserRepo) UpsertBatch(ctx context.Context, users []domain.User) error {
	if len(users) == 0 {
		return nil
	}

	q := r.GetQueryer(ctx)

	teamName := users[0].TeamName
	team, err := r.resolveTeamNameToInternalID(ctx, teamName)
	if err != nil {
		return err
	}

	upsert := r.Builder.
		Insert("users").
		Columns("user_id", "username", "is_active")

	userExternalToInternalID := make(map[string]int)

	for _, u := range users {
		upsert = upsert.Values(u.UserID, u.Username, u.IsActive)
	}

	sql, args, err := upsert.
		Suffix("ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, is_active = EXCLUDED.is_active RETURNING id, user_id").
		ToSql()

	if err != nil {
		return err
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var internalID int
	var externalID string
	for rows.Next() {
		if err = rows.Scan(&internalID, &externalID); err != nil {
			return err
		}
		userExternalToInternalID[externalID] = internalID
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	var userInternalIDs []int
	for _, u := range users {
		userInternalIDs = append(userInternalIDs, userExternalToInternalID[u.UserID])
	}

	deleteSQL, deleteArgs, err := r.Builder.
		Delete("team_member").
		Where(squirrel.Eq{"user_id": userInternalIDs}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = q.Exec(ctx, deleteSQL, deleteArgs...)
	if err != nil {
		return err
	}

	linkInsert := r.Builder.Insert("team_member").Columns("team_id", "user_id")
	for _, internalID := range userInternalIDs {
		linkInsert = linkInsert.Values(team.ID, internalID)
	}

	insertSQL, insertArgs, err := linkInsert.ToSql()
	if err != nil {
		return err
	}

	_, err = q.Exec(ctx, insertSQL, insertArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Select("u.user_id", "u.username", "u.is_active", "t.name").
		From("users u").
		LeftJoin("team_member tm ON u.id = tm.user_id").
		LeftJoin("teams t ON tm.team_id = t.id").
		Where(squirrel.Eq{"u.user_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user domain.User
	var teamName pgtype.Text
	err = q.QueryRow(ctx, sql, args...).Scan(&user.UserID, &user.Username, &user.IsActive, &teamName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if teamName.Valid {
		user.TeamName = teamName.String
	} else {
		user.TeamName = ""
	}

	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Update("users").
		Set("username", user.Username).
		Set("is_active", user.IsActive).
		Where(squirrel.Eq{"user_id": user.UserID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = q.Exec(ctx, sql, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}

	return err
}

func (r *UserRepo) GetByTeamActive(ctx context.Context, teamName string) ([]domain.User, error) {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Select("u.user_id", "u.username", "u.is_active", "t.name").
		From("users u").
		Join("team_member tm ON u.id = tm.user_id").
		Join("teams t ON tm.team_id = t.id").
		Where(squirrel.Eq{"t.name": teamName, "u.is_active": true}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var tName string
		err := rows.Scan(&user.UserID, &user.Username, &user.IsActive, &tName)
		if err != nil {
			return nil, err
		}
		user.TeamName = tName
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}
