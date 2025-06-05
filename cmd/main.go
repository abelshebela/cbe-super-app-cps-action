package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "cbe-super-app-member-users/internal/adapter/inbound/http"
	persistence "cbe-super-app-member-users/internal/adapter/outbound/persistence"
	app "cbe-super-app-member-users/internal/application/users"
	domain "cbe-super-app-member-users/internal/domain/users"

	"github.com/rs/zerolog/log"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func main() {
	ctx := context.Background()
	logger := utils.NewLogger()
	defer logger.Sync()

	client, err := config.ConnectToMongoDB()
	if err != nil {
		log.Fatal().Msgf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal().Msgf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	userRepo := persistence.NewMongoRepository(client, "cbe")
	userDomainService := domain.NewUserService(userRepo, logger)
	userAppService := app.NewApplicationHandler(*userDomainService)
	userHandler := handler.NewHTTPHandler(userAppService, logger)

	srv := handler.NewHTTPServer(userHandler)
	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info().Msgf("Received signal: %v", sig.String())
		shutdown <- srv.Shutdown(ctx)
	}()

	log.Info().Msg("Server starting on :8080")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("Server failed to start")
		os.Exit(1)
	}

	if err := <-shutdown; err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
		os.Exit(1)
	}
	log.Info().Msg("Server stopped gracefully")
}