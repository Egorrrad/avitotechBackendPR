package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/models"
)

// Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// (POST /pullRequest/create)
func (h *HTTPHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var req models.PostPullRequestCreateJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, models.NOTFOUND, "Invalid request format")
		return
	}

	ctx := r.Context()
	pr, err := h.service.CreatePullRequest(ctx, req.PullRequestId, req.AuthorId, req.PullRequestName)
	if err != nil {
		handleError(ctx, w, err)
	}

	respondJSON(w, http.StatusCreated, pr)
}

// Пометить PR как MERGED (идемпотентная операция)
// (POST /pullRequest/merge)
func (h *HTTPHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var req models.PostPullRequestMergeJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, models.NOTFOUND, "Invalid request format")
	}

	ctx := r.Context()
	mergedAt := time.Now()
	mergedPr, err := h.service.MergePullRequest(ctx, req.PullRequestId, mergedAt)
	if err != nil {
		handleError(ctx, w, err)
	}

	respondJSON(w, http.StatusOK, mergedPr)
}

// Переназначить конкретного ревьювера на другого из его команды
// (POST /pullRequest/reassign)
func (h *HTTPHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var req models.PostPullRequestReassignJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, models.NOTFOUND, "Invalid request format")
	}

	ctx := r.Context()
	reasigned, err := h.service.ReassignReviewer(ctx, req.PullRequestId, req.OldUserId)
	if err != nil {
		handleError(ctx, w, err)
	}

	respondJSON(w, http.StatusOK, reasigned)
}
