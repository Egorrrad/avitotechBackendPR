package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

func (s *Service) CreatePullRequest(ctx context.Context, prID, authorId, name string) (*domain.PullRequest, error) {
	exists, err := s.pr.Exists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPRAlreadyExists
	}

	author, err := s.users.GetByID(ctx, authorId)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, domain.ErrUserNotFound
	}

	teamMembers, err := s.users.GetByTeamActive(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	candidates := make([]string, 0)
	for _, m := range teamMembers {
		if m.UserID != authorId {
			candidates = append(candidates, m.UserID)
		}
	}

	reviewers := selectRandomReviewers(candidates, 2)

	now := time.Now()
	newPR := &domain.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   name,
		AuthorID:          authorId,
		Status:            domain.PullRequestStatusOPEN,
		AssignedReviewers: reviewers,
		CreatedAt:         &now,
	}

	if err := s.pr.Create(ctx, newPR); err != nil {
		return nil, err
	}

	return newPR, nil
}

func (s *Service) MergePullRequest(ctx context.Context, prID string, mergedAt time.Time) (*domain.PullRequest, error) {
	pr, err := s.pr.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrPullRequestNotFound
	}

	if pr.Status == domain.PullRequestStatusMERGED {
		return pr, nil
	}

	pr.Status = domain.PullRequestStatusMERGED
	pr.MergedAt = &mergedAt

	if err := s.pr.Update(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, prId string, oldReviewerId string) (*domain.PullRequest, error) {
	pr, err := s.pr.GetByID(ctx, prId)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrPullRequestNotFound
	}

	if pr.Status == domain.PullRequestStatusMERGED {
		return nil, domain.ErrChangeAfterMerge
	}

	isAssigned := false
	currentReviewersSet := make(map[string]bool)

	for _, id := range pr.AssignedReviewers {
		currentReviewersSet[id] = true
		if id == oldReviewerId {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, domain.ErrUserNotReviewer
	}

	oldReviewerUser, err := s.users.GetByID(ctx, oldReviewerId)
	if err != nil {
		return nil, err
	}
	if oldReviewerUser == nil {
		return nil, domain.ErrUserNotFound
	}

	teamMembers, err := s.users.GetByTeamActive(ctx, oldReviewerUser.TeamName)
	if err != nil {
		return nil, err
	}

	var candidates []string
	for _, m := range teamMembers {
		if m.UserID != oldReviewerId && m.UserID != pr.AuthorID && !currentReviewersSet[m.UserID] {
			candidates = append(candidates, m.UserID)
		}
	}

	if len(candidates) == 0 {
		return nil, domain.ErrNoCandidatesFound
	}

	newReviewer := selectRandomReviewers(candidates, 1)[0]

	for i, id := range pr.AssignedReviewers {
		if id == oldReviewerId {
			pr.AssignedReviewers[i] = newReviewer
			break
		}
	}

	if err := s.pr.Update(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func selectRandomReviewers(candidates []string, maxCount int) []string {
	if len(candidates) == 0 {
		return []string{}
	}

	shuffled := make([]string, len(candidates))
	copy(shuffled, candidates)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	if len(shuffled) < maxCount {
		return shuffled
	}
	return shuffled[:maxCount]
}
