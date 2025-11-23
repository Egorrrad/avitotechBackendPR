package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

type Service interface {
	PRService
	TeamService
	UserService
}

type PRService interface {
	CreatePullRequest(ctx context.Context, prId, author, name string) (*domain.PullRequest, error)
	MergePullRequest(ctx context.Context, prId string, mergedAt time.Time) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prId string, id2 string) (*domain.PullRequest, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, teamName string, members []domain.TeamMember) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}

type UserService interface {
	GetPrUserReviewer(ctx context.Context, userId string) ([]*domain.PullRequest, error)
	UpdateUserActive(ctx context.Context, userId string, active bool) (*domain.User, error)
}

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(repository Service) *HTTPHandler {
	return &HTTPHandler{
		service: repository,
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		return
	}
}
