package app

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Egorrrad/avitotechBackendPR/config"
	repo "github.com/Egorrrad/avitotechBackendPR/internal/adapter/postgres"
	"github.com/Egorrrad/avitotechBackendPR/internal/controller/http"
	"github.com/Egorrrad/avitotechBackendPR/internal/usecase"
	"github.com/Egorrrad/avitotechBackendPR/pkg/httpserver"
	"github.com/Egorrrad/avitotechBackendPR/pkg/logger"
	"github.com/Egorrrad/avitotechBackendPR/pkg/postgres"

	"github.com/grafana/pyroscope-go"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "avito.app",
		ServerAddress:   "http://pyroscope:4040",
		Logger:          pyroscope.StandardLogger,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

	if err != nil {
		panic(err)
	}

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
		l.Fatal("app - Run - postgres.New", "error", err)
	}
	defer pg.Close()

	pr, err := repo.NewPullRequestRepo(pg)
	if err != nil {
		l.Fatal("app - Run - repo.NewPullRequestRepo", "error", err)
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
		l.Info("app - Run - signal", "signal", s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify", "error", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown", "error", err)
	}
}
