package usecase

import (
	"context"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

func (s *Service) UpdateUserActive(ctx context.Context, userID string, active bool) (*domain.User, error) {
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

	return user, nil
}

func (s *Service) GetPrUserReviewer(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	return s.pr.GetByReviewerID(ctx, userID)
}
