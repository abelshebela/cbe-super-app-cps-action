package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "cbe-super-app-member-users/internal/adapter/inbound/http"
	httpAccount "cbe-super-app-member-users/internal/adapter/inbound/http/account"
	httpUser "cbe-super-app-member-users/internal/adapter/inbound/http/users"
	AccountApi "cbe-super-app-member-users/internal/adapter/outbound/api"
	persistence "cbe-super-app-member-users/internal/adapter/outbound/persistence"
	appAccount "cbe-super-app-member-users/internal/application/account"
	appUsers "cbe-super-app-member-users/internal/application/users"
	domainAccount "cbe-super-app-member-users/internal/domain/account"
	domainUsers "cbe-super-app-member-users/internal/domain/users"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func main() {
	ctx := context.Background()
	logger := utils.NewLogger()
	defer logger.Sync()

	uri := os.Getenv("MONGODB_URI")
client, err := config.ConnectToMongoDB(uri)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to MongoDB: %s", err.Error())
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal().Msgf("Failed to disconnect from MongoDB: %s", err.Error())
		}
	}()

	
	repo := persistence.NewMongoRepository(client, "cbe")
	apiClient := AccountApi.NewAccountAPIClient(logger)
	

	userDomainService := domainUsers.NewUserService(repo, logger)
	accountService := domainAccount.NewAccountService(repo, apiClient, logger)
	
	
	userAppService := appUsers.NewApplicationHandler(userDomainService)
	accountAppService := appAccount.NewApplicationHandler(accountService)
	
	
	userHandler := httpUser.NewHTTPHandler(userAppService, logger)
	accountHandler := httpAccount.NewHTTPHandler(accountAppService, logger)

	
	router := chi.NewRouter()
	httpAdapter.RegisterRoutes(router, userHandler, accountHandler)

	
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info().Msgf("Received signal: %s", sig.String())
		shutdown <- srv.Shutdown(ctx)
	}()

	log.Info().Msgf("Server starting on %s", port)
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