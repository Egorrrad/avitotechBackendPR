package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Egorrrad/avitotechBackendPR/config"
	repo "github.com/Egorrrad/avitotechBackendPR/internal/adapter/postgres"
	"github.com/Egorrrad/avitotechBackendPR/internal/controller/http"
	"github.com/Egorrrad/avitotechBackendPR/internal/usecase"
	"github.com/Egorrrad/avitotechBackendPR/pkg/httpserver"
	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
	"github.com/Egorrrad/avitotechBackendPR/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output)

	// postgres
	pg, err := postgres.New(
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.User,
		cfg.PG.Name,
		cfg.PG.Password,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	pr, err := repo.NewPullRequestRepo(pg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - repo.NewPullRequestRepo: %w", err))
	}

	// UseCase
	prsUseCase := usecase.NewService(
		repo.NewTeamRepo(pg),
		repo.NewUserRepo(pg),
		pr,
	)

	// HTTP Router (Chi)
	router := http.NewRouter(cfg, prsUseCase, l)

	// HTTP Server
	httpServer := httpserver.New(router, l, httpserver.Port(cfg.HTTP.Port))

	// Start server
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
