package middleware

import (
	"net/http"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

func Logger(l logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			l.Info("HTTP request",
				"method", r.Method,
				"url", r.URL.String(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"status", ww.Status(),
				"bytes_written", ww.BytesWritten(),
				"duration", duration.String(),
				"duration_ms", duration.Milliseconds(),
				"protocol", r.Proto,
			)
		})
	}
}
