package domain

import "errors"

// Defines values for ErrorResponseErrorCode.
const (
	NOCANDIDATE ErrorResponseErrorCode = "NO_CANDIDATE"
	NOTASSIGNED ErrorResponseErrorCode = "NOT_ASSIGNED"
	NOTFOUND    ErrorResponseErrorCode = "NOT_FOUND"
	PREXISTS    ErrorResponseErrorCode = "PR_EXISTS"
	PRMERGED    ErrorResponseErrorCode = "PR_MERGED"
	TEAMEXISTS  ErrorResponseErrorCode = "TEAM_EXISTS"

	// add new statuses
	INTERNAL ErrorResponseErrorCode = "INTERNAL_ERROR"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

// ErrorDetails defines error model for ErrorResponse
type ErrorDetails struct {
	Code    ErrorResponseErrorCode `json:"code"`
	Message string                 `json:"message"`
}

// ErrorResponseErrorCode defines model for ErrorResponse.Error.Code.
type ErrorResponseErrorCode string

var (
	ErrPullRequestNotFound = errors.New("pull request not found")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrTeamNotFound        = errors.New("team not found")
	ErrPRAlreadyExists     = errors.New("pull request already exists")
	ErrTeamAlreadyExists   = errors.New("team already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrNoCandidatesFound   = errors.New("no candidates found")
	ErrUserNotReviewer     = errors.New("user not reviewer")
	ErrChangeAfterMerge    = errors.New("cannot change after merge PR")
)
