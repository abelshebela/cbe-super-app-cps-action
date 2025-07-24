package server

import (
	"cbe-super-app-budget/config"
	"cbe-super-app-budget/platform/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) HTTPServerStart(ctx context.Context, log logger.Logger) {
	go func() {
		log.Info(ctx, "server is on", zap.String("port", s.server.Addr))
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(ctx, "HTTPServer ListenAndServer", zap.Error(err))
		}
	}()
}

func (s *HTTPServer) HTTPServerStop(ctx context.Context, log logger.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info(ctx, "Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(ctx, "Server forced to shutdown", zap.Error(err))
	} else {
		log.Info(ctx, "Server gracefully stopped")
	}
}

func NewHTTPServer(config *config.VaultConfig, handler http.Handler) *HTTPServer {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
	}

	return &HTTPServer{server: srv}
}
