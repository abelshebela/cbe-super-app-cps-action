package initiator

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/responseutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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
	persitence := InitPersistence(mongoClient, cfg.MongoDBDatabase, logger)
	logger.Infof("Persistence initialized")

	logger.Infof("Initializing domain services...")
	domain := InitDomain(minioClient, persitence, logger)
	logger.Infof("Domain services initialized")

	logger.Infof("Initializing application services...")
	application := InitApplication(domain, minioClient, logger)
	logger.Infof("Application services initialized")

	logger.Infof("Initializing adapter services...")
	adapter := InitAdapter(application, logger)
	logger.Infof("Adapter services initialized")

	logger.Infof("Initializing Chi router.....")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	r.NotFound(responseutil.NotFoundHandler)
	logger.Infof("Chi router initialized")

	logger.Infof("Initializing routes...")
	InitRoutes(r, adapter, cfg.JwtSecretKey, cfg.Key, cfg.IV, logger)
	logger.Infof("Routes initialized")

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		logger.Infof("🚀 Server started")
		logger.Infof("Server stopped with error: %v\n", server.ListenAndServe())
	}()

	sig := <-quit

	logger.Infof("server shutting down with signal: %v\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("failed to shutdown gracefully with error %v", err)
	}

	logger.Infof("Server shutdown successfully")
}
