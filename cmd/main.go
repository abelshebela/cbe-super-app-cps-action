package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	AccountAdapter "cbe-super-app-member-users/internal/adapter/inbound/http/account"
	accountRoutes "cbe-super-app-member-users/internal/adapter/inbound/http/account"
	UsersAdapter "cbe-super-app-member-users/internal/adapter/inbound/http/users"
	userRoutes "cbe-super-app-member-users/internal/adapter/inbound/http/users"
	AccountApi "cbe-super-app-member-users/internal/adapter/outbound/api"
	persistence "cbe-super-app-member-users/internal/adapter/outbound/persistence"
	appAccount "cbe-super-app-member-users/internal/application/account"
	appUsers "cbe-super-app-member-users/internal/application/users"
	domainAccount "cbe-super-app-member-users/internal/domain/account"
	domainUsers "cbe-super-app-member-users/internal/domain/users"

	authMiddleware "cbe-super-app-member-users/internal/application/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func main() {
	logger := utils.NewLogger()
	defer logger.Sync()

	cfg, err := config.Load()
	os.Setenv("GO_ENV", "dev")
	os.Setenv("KEY", "234567890-=1234567890-=1234567890-=1234567890-=") // 32 bytes key
	os.Setenv("IV", "1234567890-=12")                                   // 16 bytes IV
	if err != nil {
		logger.Fatalf("failed to load config %v", err)
	}

	mongoClient, err := config.ConnectToMongoDB(cfg.MongoDBURI)
	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
	log.Println("Connected to MongoDB!")

	minioClient, err := config.NewMinioClient(&config.VaultConfig{
		MinioEndPoint:  cfg.MinioEndPoint,
		MinioAccessKey: cfg.MinioAccessKey,
		MinioSecretKey: cfg.MinioSecretKey,
	})
	if err != nil {
		logger.Fatalf("failed to initialize minio client: %v", err)
	}

	// logger.Infof( cfg.MongoDBDatabase,"mongoosjfierrjgtiek")
	repo := persistence.NewMongoRepository(mongoClient, cfg.MongoDBDatabase)

	apiClient := AccountApi.NewAccountAPIClient(logger)

	userDomainService := domainUsers.NewUserService(repo, logger, minioClient, cfg)
	accountDomainService := domainAccount.NewAccountService(repo, apiClient, logger, cfg)

	userAppService := appUsers.InitUsersHandler(userDomainService, logger, minioClient, cfg)
	accountAppService := appAccount.InitAccountHandler(accountDomainService, logger)

	userAdapter := UsersAdapter.InitUsersAdapter(userAppService, logger, cfg)
	accountAdapter := AccountAdapter.InitAccountAdapter(accountAppService, logger)

	// Initialize router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	authMiddleware := authMiddleware.InitAuthMiddleware(cfg.JwtSecretKey, cfg.Key, cfg.IV, logger)

	// Setup routes
	userRoutes.InitUserRoutes(r, userAdapter, authMiddleware)
	accountRoutes.InitAccountRoutes(r, accountAdapter, authMiddleware)

	// Determine port
	PORT, _ := strconv.Atoi(cfg.ServerPort)
	port := strconv.Itoa(PORT)
	if port == "" {
		port = "8080"
	}

	// Configure server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Infof("🚀 Server started on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server stopped with error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sig := <-quit

	logger.Infof("Received signal: %v", sig)
	// Create shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Failed to shutdown gracefully: %v", err)
	}

	logger.Infof("Server shutting down with signal: %v", sig)
}
