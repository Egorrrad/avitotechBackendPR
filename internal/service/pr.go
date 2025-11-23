package service

import (
	"context"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/models"
	"github.com/google/uuid"
)

func (s *Service) CreatePullRequest(ctx context.Context, id, author, name string) (*models.PullRequest, error) {
	authorID := uuid.New().String() // здесь должен быть id автора
	createdAt := time.Now()

	pr := PullRequest{
		AssignedReviewers: nil,
		AuthorId:          req.AuthorId,
		CreatedAt:         &createdAt,
		MergedAt:          nil,
		PullRequestId:     req.PullRequestId,
		PullRequestName:   req.PullRequestName,
		Status:            PullRequestStatusOPEN,
	}
	return nil, nil
}
