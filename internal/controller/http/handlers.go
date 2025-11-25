package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service Service
	l       logger.Interface
	v       *validator.Validate
}

type Service interface {
	PullRequestService
	TeamService
	UserService
}

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, prID, author, name string) (*domain.PullRequestResponse, error)
	MergePullRequest(ctx context.Context, prID string) (*domain.PullRequestResponse, error)
	ReassignReviewer(ctx context.Context, prID string, id2 string) (*domain.ReassignPRResponse, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, teamName string, members []domain.TeamMember) (*domain.Team, error)
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}

type UserService interface {
	GetPrUserReviewer(ctx context.Context, userID string) (*domain.UserReviewsResponse, error)
	UpdateUserActive(ctx context.Context, userID string, active bool) (*domain.UserUpdActiveResponse, error)
}

func NewHTTPHandler(service Service,
	l logger.Interface, v *validator.Validate) *Handler {
	return &Handler{
		service: service,
		l:       l,
		v:       v,
	}
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		h.l.Error("failed to encode response", "error", err)
		return
	}
}
