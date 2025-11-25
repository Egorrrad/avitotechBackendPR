package domain

import (
	"time"
)

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

// PostPullRequestCreateJSONBody defines parameters for PostPullRequestCreate.
type PostPullRequestCreateJSONBody struct {
	AuthorID        string `json:"author_id"`
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
}

// PostPullRequestMergeJSONBody defines parameters for PostPullRequestMerge.
type PostPullRequestMergeJSONBody struct {
	PullRequestID string `json:"pull_request_id"`
}

// PostPullRequestReassignJSONBody defines parameters for PostPullRequestReassign.
type PostPullRequestReassignJSONBody struct {
	OldUserID     string `json:"old_user_id"`
	PullRequestID string `json:"pull_request_id"`
}

type PullRequestStatus string

// PullRequest defines model for PullRequest.
type PullRequest struct {
	// AssignedReviewers user_id назначенных ревьюверов (0..2)
	AssignedReviewers []string          `json:"assigned_reviewers"`
	AuthorID          string            `json:"author_id"`
	CreatedAt         *time.Time        `json:"createdAt"`
	MergedAt          *time.Time        `json:"mergedAt"`
	PullRequestID     string            `json:"pull_request_id"`
	PullRequestName   string            `json:"pull_request_name"`
	Status            PullRequestStatus `json:"status"`
}

// PullRequestShort defines model for PullRequestShort.
type PullRequestShort struct {
	AuthorID        string            `json:"author_id"`
	PullRequestID   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	Status          PullRequestStatus `json:"status"`
}

// for responses
type PullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

type ReassignPRResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by,omitempty"`
}
