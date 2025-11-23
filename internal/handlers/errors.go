package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/middleware"
	"github.com/Egorrrad/avitotechBackendPR/internal/models"
	"github.com/Egorrrad/avitotechBackendPR/internal/service"
)

func sendError(w http.ResponseWriter, status int, errCode models.ErrorResponseErrorCode, message string) {
	respondJSON(w, status, models.ErrorResponse{
		Error: models.ErrorDetails{
			errCode,
			message,
		},
	})
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	slog.ErrorContext(middleware.ErrorCtx(ctx, err), "Error: "+err.Error())

	switch {
	case errors.Is(err, service.ErrAuthorNotFound):
		sendError(w, http.StatusNotFound, models.NOTFOUND, "pull request author not found")
	case errors.Is(err, service.ErrPullRequestNotFound):
		sendError(w, http.StatusNotFound, models.NOTFOUND, "pull request not found")
	case errors.Is(err, service.ErrTeamNotFound):
		sendError(w, http.StatusNotFound, models.NOTFOUND, "team not found")
	case errors.Is(err, service.ErrPRAlreadyExists):
		sendError(w, http.StatusConflict, models.PREXISTS, "pull request already exists")

	default:
		sendError(w, http.StatusInternalServerError, models.INTERNAL, "internal server error")
	}
}
