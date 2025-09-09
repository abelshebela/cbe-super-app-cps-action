package initiator

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/responseutil"
	app_middleware "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	grpcServer "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/grpc/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func Initiator() {
	logger := utils.NewLogger()

	logger.Infof("Initializing configuration...")
	cfg := InitConfig(logger)
	logger.Infof("Configuration initialized")

	logger.Infof("Initializing MongoDB client...")
	mongoClient := InitMongo(cfg.MongoDBURI, logger)
	logger.Infof("MongoDB client initialized")

	logger.Infof("Initializing Minio client...")
	minioClient := InitMinio(cfg.MinioEndPoint, cfg.MinioAccessKey, cfg.MinioSecretKey, logger)
	logger.Infof("Minio client initialized")

	logger.Infof("Initializing persistence...")
	persitence := InitPersistence(mongoClient, cfg.MongoDBDatabase, logger, cfg)
	logger.Infof("Persistence initialized")

	logger.Infof("Initializing domain services...")
	domain := InitDomain(minioClient, persitence, logger, cfg)
	logger.Infof("Domain services initialized")

	logger.Infof("Initializing application services...")
	application := InitApplication(domain, minioClient, logger, cfg)
	logger.Infof("Application services initialized")

	logger.Infof("Initializing adapter services...")
	adapter := InitAdapter(application, minioClient, logger, cfg)
	logger.Infof("Adapter services initialized")

	logger.Infof("Initializing feedback consumer...")
	if err := InitFeedbackConsumer(mongoClient, cfg, logger); err != nil {
		logger.Errorf("Failed to initialize feedback consumer: %v", err)
	} else {
		logger.Infof("Feedback consumer initialized")
	}

	logger.Infof("Initializing Chi router.....")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app_middleware.RecoveryMiddleware(logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	r.NotFound(responseutil.NotFoundHandler)
	r.MethodNotAllowed(responseutil.MethodNOtAllowedHandler)
	logger.Infof("Chi router initialized")

	logger.Infof("Initializing routes...")
	InitRoutes(r, adapter, cfg.JwtSecretKey, cfg.Key, cfg.IV, domain.CPSActionDomain, logger)
	logger.Infof("Routes initialized")

	// Initialize gRPC server
	logger.Infof("Initializing gRPC server...")
	grpcSrv := grpcServer.NewGrpcBankServer(application.BankApplication, logger)
	logger.Infof("gRPC server initialized")

	// HTTP server
	httpServer := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	// Start HTTP server
	go func() {
		logger.Infof("🚀 HTTP Server starting on port 8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP Server stopped with error: %v", err)
		}
	}()

	// grpcSrv := .NewGRPCServer(application.BankApplication, logger)
	// Start gRPC server
	go func() {
		grpcServer.StartGrpcServer(grpcSrv)
	}()

	sig := <-quit

	logger.Infof("Servers shutting down with signal: %v", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown gRPC server
	go grpcServer.StopGrpcServer(grpcSrv)

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatalf("Failed to shutdown HTTP server gracefully: %v", err)
	}

	logger.Infof("Servers shutdown successfully")
}
