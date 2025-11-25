package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
	"github.com/Egorrrad/avitotechBackendPR/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const MaxReviewers = 2

type PullRequestRepo struct {
	*postgres.Postgres

	statusIDMap   map[domain.PullRequestStatus]int
	statusNameMap map[int]domain.PullRequestStatus
}

func NewPullRequestRepo(pg *postgres.Postgres) (*PullRequestRepo, error) {
	repo := &PullRequestRepo{
		Postgres: pg,
	}

	if err := repo.loadStatusMaps(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to load PR status maps: %w", err)
	}

	return repo, nil
}

func (r *PullRequestRepo) loadStatusMaps(ctx context.Context) error {
	q := r.GetQueryer(ctx)

	sql := "SELECT id, name FROM pr_status"
	rows, err := q.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	r.statusIDMap = make(map[domain.PullRequestStatus]int)
	r.statusNameMap = make(map[int]domain.PullRequestStatus)

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}

		status := domain.PullRequestStatus(name)

		r.statusIDMap[status] = id
		r.statusNameMap[id] = status
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

func (r *PullRequestRepo) toStatusID(status domain.PullRequestStatus) (int, error) {
	id, ok := r.statusIDMap[status]
	if !ok {
		return 0, fmt.Errorf("%w: domain status not mapped to ID: %s", status)
	}
	return id, nil
}

func (r *PullRequestRepo) toStatusName(id int) (domain.PullRequestStatus, error) {
	name, ok := r.statusNameMap[id]
	if !ok {
		return "", fmt.Errorf("internal data error: unknown PR status ID in database: %d", id)
	}
	return name, nil
}

type userInternalID struct {
	ID         int
	ExternalID string
}

func (r *PullRequestRepo) resolveExternalUserIDsToInternalIDs(ctx context.Context, externalIDs []string) ([]userInternalID, error) {
	if len(externalIDs) == 0 {
		return nil, nil
	}

	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Select("id", "user_id").
		From("users").
		Where(squirrel.Eq{"user_id": externalIDs}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []userInternalID
	for rows.Next() {
		var u userInternalID
		if err := rows.Scan(&u.ID, &u.ExternalID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if len(users) != len(externalIDs) {
		return nil, domain.ErrUserNotFound
	}
	return users, nil
}

func (r *PullRequestRepo) mapPullRequestFromRows(rows pgx.Rows) ([]*domain.PullRequest, error) {
	prMap := make(map[int]*domain.PullRequest)

	for rows.Next() {
		var (
			prID               int
			pr                 domain.PullRequest
			authorID           int
			statusID           int
			createdAt          time.Time
			mergedAt           pgtype.Timestamptz
			reviewerExternalID pgtype.Text
		)

		err := rows.Scan(
			&prID,
			&pr.PullRequestID,
			&pr.PullRequestName,
			&authorID,
			&pr.AuthorID,
			&statusID,
			&createdAt,
			&mergedAt,
			&reviewerExternalID,
		)

		if err != nil {
			return nil, err
		}

		statusName, err := r.toStatusName(statusID)
		if err != nil {
			return nil, err
		}

		p, exists := prMap[prID]
		if !exists {
			p = &domain.PullRequest{
				PullRequestID:     pr.PullRequestID,
				PullRequestName:   pr.PullRequestName,
				AuthorID:          pr.AuthorID,
				Status:            statusName,
				CreatedAt:         &createdAt,
				AssignedReviewers: make([]string, 0, MaxReviewers),
			}
			if mergedAt.Valid {
				p.MergedAt = &mergedAt.Time
			} else {
				p.MergedAt = nil
			}
			prMap[prID] = p
		}

		if reviewerExternalID.Valid && reviewerExternalID.String != "" {
			p.AssignedReviewers = append(p.AssignedReviewers, reviewerExternalID.String)
		}
	}

	var result []*domain.PullRequest
	for _, pr := range prMap {
		uniqueReviewers := make(map[string]struct{})
		var cleanReviewers []string
		for _, rev := range pr.AssignedReviewers {
			if _, ok := uniqueReviewers[rev]; !ok {
				uniqueReviewers[rev] = struct{}{}
				cleanReviewers = append(cleanReviewers, rev)
			}
		}
		pr.AssignedReviewers = cleanReviewers
		result = append(result, pr)
	}

	return result, rows.Err()
}

func (r *PullRequestRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	q := r.GetQueryer(ctx)

	authorUsers, err := r.resolveExternalUserIDsToInternalIDs(ctx, []string{pr.AuthorID})
	if errors.Is(err, domain.ErrUserNotFound) || len(authorUsers) == 0 {
		return domain.ErrAuthorNotFound
	}
	if err != nil {
		return err
	}
	authorInternalID := authorUsers[0].ID

	reviewerUsers, err := r.resolveExternalUserIDsToInternalIDs(ctx, pr.AssignedReviewers)
	if err != nil {
		return err
	}

	statusID, err := r.toStatusID(pr.Status)
	if err != nil {
		return err
	}

	var prInternalID int
	sql, args, err := r.Builder.
		Insert("pull_requests").
		Columns("pull_request_id", "pull_request_name", "author_id", "status", "created_at").
		Values(pr.PullRequestID, pr.PullRequestName, authorInternalID, statusID, time.Now()).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return err
	}

	err = q.QueryRow(ctx, sql, args...).Scan(&prInternalID)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return domain.ErrPRAlreadyExists
		}
		return err
	}

	if len(reviewerUsers) > 0 {
		reviewersInsert := r.Builder.Insert("reviewers").Columns("pr_id", "user_id")
		for _, u := range reviewerUsers {
			reviewersInsert = reviewersInsert.Values(prInternalID, u.ID)
		}

		sql, args, err = reviewersInsert.ToSql()
		if err != nil {
			return err
		}
		_, err = q.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PullRequestRepo) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Select(
			"pr.id", "pr.pull_request_id", "pr.pull_request_name",
			"pr.author_id", "author.user_id",
			"pr.status", "pr.created_at", "pr.merged_at",
			"r_user.user_id",
		).
		From("pull_requests pr").
		Join("users author ON pr.author_id = author.id").
		LeftJoin("reviewers reviewer ON pr.id = reviewer.pr_id").
		LeftJoin("users r_user ON reviewer.user_id = r_user.id").
		Where(squirrel.Eq{"pr.pull_request_id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs, err := r.mapPullRequestFromRows(rows)
	if err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, nil
	}

	return prs[0], nil
}

func (r *PullRequestRepo) Update(ctx context.Context, pr *domain.PullRequest) error {
	q := r.GetQueryer(ctx)

	var prInternalID int
	sql, args, err := r.Builder.Select("id").From("pull_requests").Where(squirrel.Eq{"pull_request_id": pr.PullRequestID}).ToSql()
	if err != nil {
		return err
	}
	err = q.QueryRow(ctx, sql, args...).Scan(&prInternalID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrPullRequestNotFound
	}
	if err != nil {
		return err
	}

	reviewerUsers, err := r.resolveExternalUserIDsToInternalIDs(ctx, pr.AssignedReviewers)
	if err != nil {
		return err
	}

	statusID, err := r.toStatusID(pr.Status)
	if err != nil {
		return err
	}

	mergedAt := pgtype.Timestamptz{}
	if pr.MergedAt != nil {
		mergedAt.Time = *pr.MergedAt
		mergedAt.Valid = true
	}

	updateBuilder := r.Builder.
		Update("pull_requests").
		Set("pull_request_name", pr.PullRequestName).
		Set("status", statusID)

	if pr.MergedAt != nil {
		updateBuilder = updateBuilder.Set("merged_at", mergedAt)
	}

	sql, args, err = updateBuilder.
		Where(squirrel.Eq{"id": prInternalID}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = q.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	deleteSQL, deleteArgs, err := r.Builder.Delete("reviewers").Where(squirrel.Eq{"pr_id": prInternalID}).ToSql()
	if err != nil {
		return err
	}
	_, err = q.Exec(ctx, deleteSQL, deleteArgs...)
	if err != nil {
		return err
	}

	if len(reviewerUsers) > 0 {
		reviewersInsert := r.Builder.Insert("reviewers").Columns("pr_id", "user_id")
		for _, u := range reviewerUsers {
			reviewersInsert = reviewersInsert.Values(prInternalID, u.ID)
		}

		insertSQL, insertArgs, err := reviewersInsert.ToSql()
		if err != nil {
			return err
		}
		_, err = q.Exec(ctx, insertSQL, insertArgs...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PullRequestRepo) GetByReviewerID(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	q := r.GetQueryer(ctx)

	reviewerUsers, err := r.resolveExternalUserIDsToInternalIDs(ctx, []string{userID})
	if errors.Is(err, domain.ErrUserNotFound) || len(reviewerUsers) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	reviewerInternalID := reviewerUsers[0].ID

	sql, args, err := r.Builder.
		Select(
			"pr.id", "pr.pull_request_id", "pr.pull_request_name",
			"pr.author_id", "author.user_id",
			"pr.status", "pr.created_at", "pr.merged_at",
			"r_user.user_id",
		).
		From("reviewers r").
		Join("pull_requests pr ON r.pr_id = pr.id").
		Join("users author ON pr.author_id = author.id").
		LeftJoin("reviewers r2 ON pr.id = r2.pr_id").
		LeftJoin("users r_user ON r2.user_id = r_user.id").
		Where(squirrel.Eq{"r.user_id": reviewerInternalID}).
		Distinct().
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.mapPullRequestFromRows(rows)
}

func (r *PullRequestRepo) Exists(ctx context.Context, id string) (bool, error) {
	q := r.GetQueryer(ctx)

	sql, args, err := r.Builder.
		Select("COUNT(*)").
		From("pull_requests").
		Where(squirrel.Eq{"pull_request_id": id}).
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
