package initiator

import (
	"context"
	"fmt"
	"runtime"

	"cbe-super-app-cps-action/cmd/server"
	local "cbe-super-app-cps-action/config"
	"cbe-super-app-cps-action/internal/storage/api"
	"cbe-super-app-cps-action/internal/storage/external_call"

	// "cbe-super-app-cps-action/platform/logger"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"log"

	// local_logger "cbe-super-app-cps-action/platform/logger"
	"github.com/go-chi/chi/v5"
)

func Init(ctx context.Context) {
	done := make(chan struct{})

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

	redis := InitRedis(cfg, logger)
	logger.Infof("Initializing redis...")
	redisStorage := InitRedisStorageLayer(redis, logger)
	logger.Infof("redis initialized")

	redisRepository := redisStorage.GetRedisRepository()

	// oracleDB := InitOracle(cfg.OracleConnectionString, logger)
	// logger.Infof("Oracle database initialized")

	// logger.Infof("Initializing Oracle DB client...")
	// OraclePersistence := InitOraclePersistence(oracleDB, logger)
	// logger.Infof("Oracle DB client initialized")

	logger.Infof("Initializing SMS service...")
	smsService := external_call.NewSMSPersistence(cfg.SMSBaseURL, logger)
	logger.Infof("SMS service initialized")

	logger.Infof("Initializing account lookup service...")
	accountLookupService := InitAccountLookupService(cfg.CBEBaseURL, logger)
	logger.Infof("Account lookup service initialized")

	sessionGRPCClient, clientStore, err := api.NewSessionGRPCClient("cfg.CommonSvcGrpcAddress", logger) // TODO: Add to config
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC session client: %v", err)
	}
	defer func() {
		clientStore.Close()
	}()

	sitotagRPCClient, err := api.NewSitotagRPCClient(ctx, logger, cfg.CbeToCbeGrpcAddress)
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC client for sitota")
	}
	defer sitotagRPCClient.Close()

	defer local.DisconnectMongo(ctx, mongoClient, logger)

	logger.Infof("initialize service layer")
	serviceLayer := InitServiceLayer(mongoClient, persitence, logger, sessionGRPCClient, cfg, minioClient, accountLookupService, redisRepository, *smsService)
	// serviceLayer := InitServiceLayer(mongoClient, persitence, logger, sessionGRPCClient, cfg, minioClient, accountLookupService, OraclePersistence)

	go func() {
		if err := InitFeedbackConsumer(serviceLayer.Feedback, cfg, logger); err != nil {
			logger.Errorf("Failed to start feedback consumer: %v", err)
		}
	}()

	logger.Infof("initialize handler layer")
	handlerLayer := InitHandler(serviceLayer, logger)

	r := chi.NewRouter()
	InitRoute(ctx, r, handlerLayer, logger, cfg)

	go func() {
		fmt.Println("Goroutines: ", runtime.NumGoroutine())
	}()

	grpcHandlers := server.NewGrpcServer(serviceLayer.Bank, serviceLayer.Wallet, serviceLayer.ServiceDetails, serviceLayer.Topup, logger)
	srv := server.NewHTTPServer(cfg, r)

	grpcServer, lis := server.StartGrpcServer(grpcHandlers, logger)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
		done <- struct{}{}
	}()

	go func() {
		srv.HTTPServerStart(ctx, logger)
		done <- struct{}{}
	}()
	<-done
	logger.Infof("Shutdown signal received. Stopping servers...")
	srv.HTTPServerStop(ctx, logger)
	server.StopGrpcServer(grpcServer, logger)
}
