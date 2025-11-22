package prservice

import (
	"log/slog"

	"github.com/Egorrrad/avitotechBackendPR/internal/middleware"
	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
)

func InitLogging() {
	handler := logger.InitSlogHandler("info", "", "json")
	handler = middleware.NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}
