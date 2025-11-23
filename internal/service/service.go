package service

import "errors"

type Repository interface {
	PullRequestRepo
	TeamRepo
	UserRepo
}

type PullRequestRepo interface {
}

type TeamRepo interface {
}

type UserRepo interface {
}

type Service struct {
	repo Repository
}

var (
	ErrPullRequestNotFound = errors.New("pull request not found")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrTeamNotFound        = errors.New("team not found")
	ErrPRAlreadyExists     = errors.New("pull request already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrNoCandidatesFound   = errors.New("no candidates found")
	ErrUserNotReviewer     = errors.New("user not reviewer")
	ErrChangeAfterMerge    = errors.New("cannot change after merge PR")
)
