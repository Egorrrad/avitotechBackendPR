package logger

import (
	"fmt"
	"os"
	"strings"

	"log/slog"
)

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *slog.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level, format, output string) *Logger {
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

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(logOutput, opts)
	default:
		handler = slog.NewTextHandler(logOutput, opts)
	}

	logger := slog.New(handler)

	return &Logger{
		logger: logger,
	}
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg(slog.LevelDebug, message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(slog.LevelInfo, message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(slog.LevelWarn, message, args...)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.msg(slog.LevelError, message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg(slog.LevelError, message, args...) // Fatal логируем как Error
	os.Exit(1)
}

func (l *Logger) log(level slog.Level, message string, args ...interface{}) {
	l.logger.Log(nil, level, message, args...)
}

func (l *Logger) msg(level slog.Level, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(level, msg.Error(), args...)
	case string:
		l.log(level, msg, args...)
	default:
		l.log(level, fmt.Sprintf("message %v has unknown type %v", message, msg), args...)
	}
}
