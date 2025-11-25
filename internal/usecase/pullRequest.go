package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

func (s *Service) CreatePullRequest(ctx context.Context, prID, authorID, name string) (*domain.PullRequestResponse, error) {
	exists, err := s.pr.Exists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPRAlreadyExists
	}

	author, err := s.users.GetByID(ctx, authorID)
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
		if m.UserID != authorID {
			candidates = append(candidates, m.UserID)
		}
	}

	reviewers := selectRandomReviewers(candidates, 2)

	now := time.Now()
	newPR := &domain.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   name,
		AuthorID:          authorID,
		Status:            domain.PullRequestStatusOPEN,
		AssignedReviewers: reviewers,
		CreatedAt:         &now,
	}

	if _, err := s.pr.Create(ctx, newPR); err != nil {
		return nil, err
	}

	respNewPr := &domain.PullRequestResponse{
		PR: *newPR,
	}

	return respNewPr, nil
}

func (s *Service) MergePullRequest(ctx context.Context, prID string) (*domain.PullRequestResponse, error) {
	pr, err := s.pr.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrPullRequestNotFound
	}

	if pr.Status == domain.PullRequestStatusMERGED {
		return &domain.PullRequestResponse{PR: *pr}, nil
	}

	mergedAt := time.Now()
	pr.Status = domain.PullRequestStatusMERGED
	pr.MergedAt = &mergedAt

	if err := s.pr.Update(ctx, pr); err != nil {
		return nil, err
	}

	return &domain.PullRequestResponse{PR: *pr}, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, prId string, oldReviewerId string) (*domain.ReassignPRResponse, error) {
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

	respPR := &domain.ReassignPRResponse{
		PR:         *pr,
		ReplacedBy: newReviewer,
	}
	return respPR, nil
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
