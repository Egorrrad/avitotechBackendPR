package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/internal/models"
)

// Создать команду с участниками (создаёт/обновляет пользователей)
// (POST /team/add)
func (h *HTTPHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var req models.PostTeamAddJSONRequestBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, models.NOTFOUND, "Invalid request format")
	}

	ctx := r.Context()
	team, err := h.service.CreateTeam(ctx, req.TeamName, req.Members)
	if err != nil {
		handleError(ctx, w, err)
	}

	respondJSON(w, http.StatusCreated, team)
}

// Получить команду с участниками
// (GET /team/get)
func (h *HTTPHandler) GetTeamGet(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	ctx := r.Context()
	team, err := h.service.GetTeam(ctx, teamName)
	if err != nil {
		handleError(ctx, w, err)
	}

	respondJSON(w, http.StatusOK, team)
}
