package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

func (h *Handler) sendError(w http.ResponseWriter, status int, errCode domain.ErrorResponseErrorCode, message string) {
	h.respondJSON(w, status, domain.ErrorResponse{
		Error: domain.ErrorDetails{
			Code:    errCode,
			Message: message,
		},
	})
}

func (h *Handler) handleError(ctx context.Context, w http.ResponseWriter, err error) {
	h.l.Error("Handler error", "error", err)

	switch {
	case errors.Is(err, domain.ErrAuthorNotFound):
		h.sendError(w, http.StatusNotFound, domain.NOTFOUND, "pull request author not found")
	case errors.Is(err, domain.ErrPullRequestNotFound):
		h.sendError(w, http.StatusNotFound, domain.NOTFOUND, "pull request not found")
	case errors.Is(err, domain.ErrTeamNotFound):
		h.sendError(w, http.StatusNotFound, domain.NOTFOUND, "team not found")
	case errors.Is(err, domain.ErrPRAlreadyExists):
		h.sendError(w, http.StatusConflict, domain.PREXISTS, "pull request already exists")
	case errors.Is(err, domain.ErrTeamAlreadyExists):
		h.sendError(w, http.StatusConflict, domain.TEAMEXISTS, "team already exists")

	default:
		h.sendError(w, http.StatusInternalServerError, domain.INTERNAL, "internal server error")
	}
}
