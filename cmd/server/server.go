package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.uber.org/zap"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) HTTPServerStart(ctx context.Context, log utils.Logger) {
	go func() {
		log.Infof("server is on", zap.String("port", s.server.Addr))
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTPServer ListenAndServer", zap.Error(err))
		}
	}()
}

func (s *HTTPServer) HTTPServerStop(ctx context.Context, log utils.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Infof("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown", zap.Error(err))
	} else {
		log.Infof("Server gracefully stopped")
	}
}

func NewHTTPServer(config *config.VaultConfig, handler http.Handler) *HTTPServer {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  time.Duration(config.ServerTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ServerTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.ServerTimeout) * time.Second,
	}

	return &HTTPServer{server: srv}
}
