package logger

import (
	"log/slog"
	"os"
	"strings"
)

func InitSlogHandler(level, format, output string) slog.Handler {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var logOutput *os.File
	if output == "stdout" || output == "" {
		logOutput = os.Stdout
	} else {
		var err error
		logOutput, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			logOutput = os.Stdout
		}
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(logOutput, &slog.HandlerOptions{Level: logLevel})
	default:
		handler = slog.NewTextHandler(logOutput, &slog.HandlerOptions{Level: logLevel})
	}

	return handler
}
