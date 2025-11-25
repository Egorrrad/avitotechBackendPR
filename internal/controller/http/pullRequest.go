package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

// Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// (POST /pullRequest/create)
func (h *Handler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var req domain.PostPullRequestCreateJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, domain.NOTFOUND, "Invalid request format")
		return
	}

	ctx := r.Context()
	pr, err := h.service.CreatePullRequest(ctx, req.PullRequestID, req.AuthorID, req.PullRequestName)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, pr)
}

// Пометить PR как MERGED (идемпотентная операция)
// (POST /pullRequest/merge)
func (h *Handler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var req domain.PostPullRequestMergeJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, domain.NOTFOUND, "Invalid request format")
		return
	}

	ctx := r.Context()
	mergedAt := time.Now()
	mergedPr, err := h.service.MergePullRequest(ctx, req.PullRequestID, mergedAt)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, mergedPr)
}

// Переназначить конкретного ревьювера на другого из его команды
// (POST /pullRequest/reassign)
func (h *Handler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var req domain.PostPullRequestReassignJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, domain.NOTFOUND, "Invalid request format")
		return
	}

	ctx := r.Context()
	reasigned, err := h.service.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, reasigned)
}
