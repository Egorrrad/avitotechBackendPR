package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Service interface {
	PRService
	TeamService
	UserService
}

type PRService interface {
}

type TeamService interface {
}

type UserService interface {
}

type HTTPHandler struct {
	service *Service
}

func NewHTTPHandler(repository *Service) *HTTPHandler {
	return &HTTPHandler{
		service: repository,
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		return
	}
}
