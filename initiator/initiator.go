package initiator

import (
	"context"
	"fmt"
	"runtime"

	"cbe-super-app-cps-action/cmd/server"
	local "cbe-super-app-cps-action/config"
	"cbe-super-app-cps-action/internal/storage/api"

	// "cbe-super-app-cps-action/platform/logger"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

func Init(ctx context.Context) {
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
	persitence := InitPersistanceLayer(mongoClient, cfg.MongoDBDatabase, logger)
	logger.Infof("Persistence initialized")

	logger.Infof("Initializing account lookup service...")
	accountLookupService := InitAccountLookupService(cfg, logger)
	logger.Infof("Account lookup service initialized")

	sessionGRPCClient, clientStore, err := api.NewSessionGRPCClient("cfg.CommonSvcGrpcAddress", logger) // TODO: Add to config
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC session client: %v", err)
	}
	defer func() {
		clientStore.Close()
	}()

	defer local.DisconnectMongo(ctx, mongoClient, logger)

	logger.Infof("initialize service layer")
	serviceLayer := InitServiceLayer(mongoClient, persitence, logger, sessionGRPCClient, cfg, minioClient, accountLookupService)

	go func() {
		if err := InitFeedbackConsumer(serviceLayer.Feedback, cfg, logger); err != nil {
			logger.Errorf("Failed to start feedback consumer: %v", err)
		}
	}()

	logger.Infof("initialize handler layer")
	handlerLayer := InitHandler(serviceLayer, logger)

	r := chi.NewRouter()
	InitRoute(ctx, r, handlerLayer, logger, cfg)

	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	go func() {
		fmt.Println("Goroutines: ", runtime.NumGoroutine())
	}()

	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	grpcServer := server.NewGrpcBankServer(serviceLayer.Bank, logger)
	srv := server.NewHTTPServer(cfg, r)

	// Start servers
	go func() {
		server.StartGrpcServer(grpcServer, cfg.GRPCPort)
	}()

	go func() {
		srv.HTTPServerStart(ctx, logger)
	}()
	srv.HTTPServerStop(ctx, logger)
	server.StopGrpcServer(grpcServer)
}
