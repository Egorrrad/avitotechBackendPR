package usecase

import (
	"context"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

func (s *Service) UpdateUserActive(ctx context.Context, userID string, active bool) (*domain.UserUpdActiveResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	user.IsActive = active

	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}

	return &domain.UserUpdActiveResponse{User: *user}, nil
}

func (s *Service) GetPrUserReviewer(ctx context.Context, userID string) (*domain.UserReviewsResponse, error) {
	prs, err := s.pr.GetByReviewerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserReviewsResponse{
		UserID:       userID,
		PullRequests: prs,
	}, nil
}
