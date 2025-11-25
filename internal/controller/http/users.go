package http

import (
	"encoding/json"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

// Получить PR'ы, где пользователь назначен ревьювером
// (GET /users/getReview)
func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	ctx := r.Context()
	prs, err := h.service.GetPrUserReviewer(ctx, userID)
	if err != nil {
		h.handleError(ctx, w, err)
	}

	h.respondJSON(w, http.StatusOK, prs)
}

// Установить флаг активности пользователя
// (POST /users/setIsActive)
func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var req domain.PostUsersSetIsActiveJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, domain.NOTFOUND, "Invalid request format")
	}

	ctx := r.Context()
	updUser, err := h.service.UpdateUserActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		h.handleError(ctx, w, err)
	}

	h.respondJSON(w, http.StatusOK, updUser)
}
