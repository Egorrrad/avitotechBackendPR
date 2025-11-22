package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Repository interface {
}

type Server struct {
	repo Repository
}

func NewServer(repository Repository) Server {
	return Server{
		repo: repository,
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
