package http

import (
	"net/http"

	"github.com/Egorrrad/avitotechBackendPR/config"
	"github.com/Egorrrad/avitotechBackendPR/internal/controller/http/middleware"
	"github.com/Egorrrad/avitotechBackendPR/internal/usecase"
	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(cfg *config.Config, t *usecase.Service, l logger.Interface) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger(l))
	r.Use(middleware.Recovery(l))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		// В Chi метрики подключаются как handler
		r.Handle("/metrics", promhttp.Handler())
	}

	// K8s probe
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	v := validator.New()
	h := NewHTTPHandler(t, l, v)

	// Routers
	// pullRequest routes
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.PostPullRequestCreate)
		r.Post("/merge", h.PostPullRequestMerge)
		r.Post("/reassign", h.PostPullRequestReassign)
	})

	// team routes
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.PostTeamAdd)
		r.Get("/get", h.GetTeamGet)
	})

	// users routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/getReview", h.GetUsersGetReview)
		r.Post("/setIsActive", h.PostUsersSetIsActive)
	})

	return r
}
