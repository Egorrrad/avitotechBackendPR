package http

import (
	"encoding/json"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

// Получить PR'ы, где пользователь назначен ревьювером
// (GET /users/getReview)
func (h *HTTPHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	ctx := r.Context()
	prs, err := h.service.GetPrUserReviewer(ctx, userId)
	if err != nil {
		handleError(ctx, w, err)
	}

	respondJSON(w, http.StatusOK, prs)
}

// Установить флаг активности пользователя
// (POST /users/setIsActive)
func (h *HTTPHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var req domain.PostUsersSetIsActiveJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, domain.NOTFOUND, "Invalid request format")
	}

	ctx := r.Context()
	updUser, err := h.service.UpdateUserActive(ctx, req.UserId, req.IsActive)
	if err != nil {
		handleError(ctx, w, err)
	}

	respondJSON(w, http.StatusOK, updUser)
}
