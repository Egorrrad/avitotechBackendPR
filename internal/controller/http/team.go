package http

import (
	"encoding/json"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
)

// Создать команду с участниками (создаёт/обновляет пользователей)
// (POST /team/add)
func (h *Handler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var req domain.PostTeamAddJSONRequestBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, domain.NOTFOUND, "Invalid request format")
	}

	ctx := r.Context()
	team, err := h.service.CreateTeam(ctx, req.TeamName, req.Members)
	if err != nil {
		h.handleError(ctx, w, err)
	}

	h.respondJSON(w, http.StatusCreated, team)
}

// Получить команду с участниками
// (GET /team/get)
func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	ctx := r.Context()
	team, err := h.service.GetTeam(ctx, teamName)
	if err != nil {
		h.handleError(ctx, w, err)
	}

	h.respondJSON(w, http.StatusOK, team)
}
