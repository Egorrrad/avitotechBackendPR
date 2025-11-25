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

func (s *Service) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.ReassignPRResponse, error) {
	pr, err := s.validateReassignRequest(ctx, prID, oldReviewerID)
	if err != nil {
		return nil, err
	}

	newReviewer, err := s.findReplacementReviewer(ctx, pr, oldReviewerID)
	if err != nil {
		return nil, err
	}

	s.replaceReviewer(pr, oldReviewerID, newReviewer)

	if err := s.pr.Update(ctx, pr); err != nil {
		return nil, err
	}

	return &domain.ReassignPRResponse{
		PR:         *pr,
		ReplacedBy: newReviewer,
	}, nil
}

func (s *Service) validateReassignRequest(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, error) {
	pr, err := s.pr.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, domain.ErrPullRequestNotFound
	}

	if pr.Status == domain.PullRequestStatusMERGED {
		return nil, domain.ErrChangeAfterMerge
	}

	if !s.isUserAssigned(pr, oldReviewerID) {
		return nil, domain.ErrUserNotReviewer
	}

	return pr, nil
}

func (s *Service) isUserAssigned(pr *domain.PullRequest, userID string) bool {
	for _, id := range pr.AssignedReviewers {
		if id == userID {
			return true
		}
	}
	return false
}

func (s *Service) findReplacementReviewer(ctx context.Context, pr *domain.PullRequest, oldReviewerID string) (string, error) {
	oldReviewerUser, err := s.users.GetByID(ctx, oldReviewerID)
	if err != nil {
		return "", err
	}
	if oldReviewerUser == nil {
		return "", domain.ErrUserNotFound
	}

	candidates, err := s.getEligibleCandidates(ctx, pr, oldReviewerUser.TeamName, oldReviewerID)
	if err != nil {
		return "", err
	}

	if len(candidates) == 0 {
		return "", domain.ErrNoCandidatesFound
	}

	return selectRandomReviewers(candidates, 1)[0], nil
}

func (s *Service) getEligibleCandidates(ctx context.Context, pr *domain.PullRequest, teamName, oldReviewerID string) ([]string, error) {
	teamMembers, err := s.users.GetByTeamActive(ctx, teamName)
	if err != nil {
		return nil, err
	}

	currentReviewers := s.getCurrentReviewersSet(pr)

	var candidates []string
	for _, m := range teamMembers {
		if m.UserID != oldReviewerID && m.UserID != pr.AuthorID && !currentReviewers[m.UserID] {
			candidates = append(candidates, m.UserID)
		}
	}

	return candidates, nil
}

func (s *Service) getCurrentReviewersSet(pr *domain.PullRequest) map[string]bool {
	currentReviewersSet := make(map[string]bool)
	for _, id := range pr.AssignedReviewers {
		currentReviewersSet[id] = true
	}
	return currentReviewersSet
}

func (s *Service) replaceReviewer(pr *domain.PullRequest, oldReviewerID, newReviewer string) {
	for i, id := range pr.AssignedReviewers {
		if id == oldReviewerID {
			pr.AssignedReviewers[i] = newReviewer
			break
		}
	}
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
