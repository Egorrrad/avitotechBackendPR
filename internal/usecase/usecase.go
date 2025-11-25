package usecase

import (
	"context"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

type (
	PullRequestRepo interface {
		Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
		GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
		Update(ctx context.Context, pr *domain.PullRequest) error
		GetByReviewerID(ctx context.Context, userID string) ([]*domain.PullRequest, error)
		Exists(ctx context.Context, id string) (bool, error)
	}

	TeamRepo interface {
		Create(ctx context.Context, team *domain.Team) error
		GetByName(ctx context.Context, name string) (*domain.Team, error)
		Exists(ctx context.Context, name string) (bool, error)
	}

	UserRepo interface {
		UpsertBatch(ctx context.Context, users []domain.User) error
		GetByID(ctx context.Context, id string) (*domain.User, error)
		Update(ctx context.Context, user *domain.User) error
		GetByTeamActive(ctx context.Context, teamName string) ([]domain.User, error)
	}
)

type Service struct {
	teams TeamRepo
	users UserRepo
	pr    PullRequestRepo
}

func NewService(team TeamRepo, users UserRepo, pr PullRequestRepo) *Service {
	return &Service{
		teams: team,
		users: users,
		pr:    pr,
	}
}
