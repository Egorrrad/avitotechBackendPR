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
	AuthorId        string `json:"author_id"`
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
}

// PostPullRequestMergeJSONBody defines parameters for PostPullRequestMerge.
type PostPullRequestMergeJSONBody struct {
	PullRequestId string `json:"pull_request_id"`
}

// PostPullRequestReassignJSONBody defines parameters for PostPullRequestReassign.
type PostPullRequestReassignJSONBody struct {
	OldUserId     string `json:"old_user_id"`
	PullRequestId string `json:"pull_request_id"`
}

type PullRequestStatus string

// PullRequest defines model for PullRequest.
type PullRequest struct {
	// AssignedReviewers user_id назначенных ревьюверов (0..2)
	AssignedReviewers []string          `json:"assigned_reviewers"`
	AuthorId          string            `json:"author_id"`
	CreatedAt         *time.Time        `json:"createdAt"`
	MergedAt          *time.Time        `json:"mergedAt"`
	PullRequestId     string            `json:"pull_request_id"`
	PullRequestName   string            `json:"pull_request_name"`
	Status            PullRequestStatus `json:"status"`
}

// PullRequestShort defines model for PullRequestShort.
type PullRequestShort struct {
	AuthorId        string            `json:"author_id"`
	PullRequestId   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	Status          PullRequestStatus `json:"status"`
}
