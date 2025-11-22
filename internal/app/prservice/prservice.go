package prservice

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Egorrrad/avitotechBackendPR/internal/repository"
	"github.com/gorilla/mux"
)

type APIServer struct {
	cfg    *Config
	router *mux.Router
	store  *repository.DataStorage
}

func NewServer(cfg *Config) (*APIServer, error) {
	srv := &APIServer{
		cfg: cfg,
		// store:  store,
	}

	srv.router = InitRouter()

	return srv, nil
}

func (s *APIServer) Start(ctx context.Context) {
	server := http.Server{
		Addr:    ":" + s.cfg.Service.Port,
		Handler: s.router,
	}

	shutdownComplete := handleShutdown(func() {
		// ждем пока все выполнится
		ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctxCancel); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
	})

	slog.Info("Service started", "adr", server.Addr)
	if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		<-shutdownComplete
	} else {
		slog.Error("http.ListenAndServe failed", "error", err)
	}

	slog.Info("Shutdown gracefully")
}

func handleShutdown(onShutdownSignal func()) <-chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		sig := <-shutdownSignal
		slog.Info("Received shutdown signal", "signal", sig.String())

		onShutdownSignal()
		close(shutdown)
	}()

	return shutdown
}
