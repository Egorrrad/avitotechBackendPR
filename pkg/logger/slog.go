package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
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

type sourceHandler struct {
	handler slog.Handler
	skip    int
}

func (h *sourceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *sourceHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= slog.LevelDebug {
		pc, file, line, ok := runtime.Caller(h.skip)
		if ok {
			funcName := runtime.FuncForPC(pc)
			if funcName != nil {
				record.AddAttrs(
					slog.Group("source",
						slog.String("function", funcName.Name()),
						slog.String("file", file),
						slog.Int("line", line),
					),
				)
			}
		}
	}
	return h.handler.Handle(ctx, record)
}

func (h *sourceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &sourceHandler{handler: h.handler.WithAttrs(attrs), skip: h.skip}
}

func (h *sourceHandler) WithGroup(name string) slog.Handler {
	return &sourceHandler{handler: h.handler.WithGroup(name), skip: h.skip}
}

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
		Level: logLevel,
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(logOutput, opts)
	default:
		handler = slog.NewTextHandler(logOutput, opts)
	}

	handler = &sourceHandler{handler: handler, skip: 5}
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
	l.msg(slog.LevelError, message, args...)
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
