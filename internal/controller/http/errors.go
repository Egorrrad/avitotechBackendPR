package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
	"github.com/Egorrrad/avitotechBackendPR/internal/middleware"
)

func sendError(w http.ResponseWriter, status int, errCode domain.ErrorResponseErrorCode, message string) {
	respondJSON(w, status, domain.ErrorResponse{
		Error: domain.ErrorDetails{
			errCode,
			message,
		},
	})
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	slog.ErrorContext(middleware.ErrorCtx(ctx, err), "Error: "+err.Error())

	switch {
	case errors.Is(err, domain.ErrAuthorNotFound):
		sendError(w, http.StatusNotFound, domain.NOTFOUND, "pull request author not found")
	case errors.Is(err, domain.ErrPullRequestNotFound):
		sendError(w, http.StatusNotFound, domain.NOTFOUND, "pull request not found")
	case errors.Is(err, domain.ErrTeamNotFound):
		sendError(w, http.StatusNotFound, domain.NOTFOUND, "team not found")
	case errors.Is(err, domain.ErrPRAlreadyExists):
		sendError(w, http.StatusConflict, domain.PREXISTS, "pull request already exists")

	default:
		sendError(w, http.StatusInternalServerError, domain.INTERNAL, "internal server error")
	}
}
