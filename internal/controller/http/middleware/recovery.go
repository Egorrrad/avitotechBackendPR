package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/Egorrrad/avitotechBackendPR/internal/domain"
	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
)

func Recovery(l logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := string(debug.Stack())

					l.Error("HTTP request panic recovered",
						"panic", err,
						"stack_trace", stack,
						"request.method", r.Method,
						"request.url", r.URL.String(),
						"request.remote_addr", r.RemoteAddr,
						"request.user_agent", r.UserAgent(),
						"request.host", r.Host,
						"request.proto", r.Proto,
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					resp := domain.ErrorResponse{
						Error: domain.ErrorDetails{
							Code:    domain.INTERNAL,
							Message: "Internal server error",
						},
					}

					if err := json.NewEncoder(w).Encode(resp); err != nil {
						l.Error("Failed to encode error response", "error", err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
