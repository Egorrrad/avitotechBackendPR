package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

func buildRequestMessage(r *http.Request, status, length int, duration time.Duration) string {
	var result strings.Builder

	result.WriteString(r.RemoteAddr)
	result.WriteString(" - ")
	result.WriteString(r.Method)
	result.WriteString(" ")
	result.WriteString(r.URL.String())
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(status))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(length))
	result.WriteString(" - ")
	result.WriteString(duration.String())

	return result.String()
}

func Logger(l logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			message := buildRequestMessage(r, ww.Status(), ww.BytesWritten(), duration)
			l.Info(message)
		})
	}
}
