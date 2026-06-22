package initiator

import (
	"context"
	"net/http"
	"time"

	"cbe-super-app-cps-action/cmd/client"
	"cbe-super-app-cps-action/cmd/server"
	local "cbe-super-app-cps-action/config"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/storage/api"

	mid "cbe-super-app-cps-action/internal/handlers/middleware"

	"cbe-super-app-cps-action/platform/telemetry"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/middleware"
	shared_producer "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/producer"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils/encryption"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/hugokessem/coreio/core"
)

func Init(ctx context.Context) {
	done := make(chan struct{})
	logger := utils.NewLogger()
	zapLogger := InitZapLogger()

	logger.Infof("Initializing configuration...")
	cfg := InitConfig(logger)
	logger.Infof("Configuration initialized")

	coreConfig := core.CBECoreCredential{
		Username: cfg.CbeCoreUsername,
		Password: cfg.CbeCorePassword,
		Url:      cfg.CbeCoreUrl,
	}
	coreInterface := core.NewCBECoreAPI(&coreConfig)

	// Initialize OpenTelemetry Tracing using platform/telemetry package
	logger.Infof("Initializing OpenTelemetry Tracing...")

	// Build OTEL config from a dedicated helper that reads env vars and provides defaults.
	otelCfg := NewOtelConfig(*cfg)

	if otelCfg.Enabled {
		logger.Infof("Checking OTLP endpoint connectivity at %s...", otelCfg.OTLPEndpoint)
		if err := telemetry.CheckOTLPConnection(ctx, otelCfg.OTLPEndpoint); err != nil {
			logger.Warnf("OTLP endpoint connectivity check failed: %v. Continuing without tracing.", err)
			logger.Infof("To enable tracing, ensure Jaeger or OTLP collector is running at %s", otelCfg.OTLPEndpoint)
		} else {
			logger.Infof("OTLP endpoint is reachable. Initializing tracer provider...")
			tracerProvider, err := telemetry.NewTracerProvider(ctx, otelCfg)
			if err != nil {
				logger.Warnf("Failed to initialize OpenTelemetry Tracing: %v. Continuing without tracing.", err)
			} else {
				logger.Infof("✓ OpenTelemetry Tracing initialized successfully")
				defer func() {
					if err := tracerProvider.Shutdown(ctx); err != nil {
						logger.Errorf("Failed to shutdown tracer provider: %v", err)
					}
				}()
			}
		}
	} else {
		logger.Infof("OpenTelemetry is disabled in configuration; skipping tracer initialization")
	}

	// Start Prometheus metrics server on :9090
	metricsPort := "9090"
	go func() {
		mux := http.NewServeMux()
		telemetry.RegisterMetricsEndpoint(mux)
		logger.Infof("Starting metrics endpoint on port %s", metricsPort)
		if err := http.ListenAndServe(":"+metricsPort, mux); err != nil {
			logger.Errorf("metrics server failed: %v", err)
		}
	}() // 5G

	logger.Infof("Initializing MongoDB client...")
	mongoClient := InitMongo(cfg.MongoDBURI, logger) //23G - 5G = 18G
	logger.Infof("MongoDB client initialized")

	logger.Infof("Initializing Minio client...")
	minioClient := InitMinio(*cfg, logger) // 24G -23G= 1G
	logger.Infof("Minio client initialized")

	presignClient := s3.NewPresignClient(minioClient)

	logger.Infof("initializing kafka")
	notificationProducer, clientOrchestrationProducer, accessListSegmentationProducer := InitKafkaService(cfg, logger) //27G - 24G= 3G
	logger.Infof("kafka initialized")

	// Expose orchestration producer to handler/middleware layer for BPS action publishing.
	mid.InitClientOrchestrationProducer(clientOrchestrationProducer)

	// Init shared kafka notification producer
	sharedKafkaProducer, err := shared_producer.NewNotificationProducer(*cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize shared Kafka notification producer: %v", err)
	}

	logger.Infof("Initializing persistence...")
	notificationApi := cfg.SMSBaseURL

	sharedRedis, redis, err := InitRedis(ctx, cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize Redis: %v", err)
	}
	logger.Infof("Initializing redis...")
	redisStorage := InitRedisStorageLayer(redis, logger)
	logger.Infof("redis initialized")

	// encMiddleWare := encryption.NewEncryptionImpl(zapLogger)
	encryptionMiddleware := middleware.InitEncMiddleware(cfg, sharedRedis, zapLogger)

	redisRepository := redisStorage.GetRedisRepository()

	// Initialize queue system (memory + Redis backends, dedup, metrics)
	logger.Infof("Initializing queue system...")
	queueInfra := InitQueueSystem(redis, logger)
	queueInfra.Manager.Start(ctx, queueWorkers)
	defer queueInfra.Manager.Stop()
	logger.Infof("Queue system initialized and started")

	persistence := InitPersistanceLayer(mongoClient, cfg.MongoDBDatabase, coreInterface, notificationApi, *notificationProducer, sharedKafkaProducer, *clientOrchestrationProducer, *accessListSegmentationProducer, redisRepository, cfg, logger)
	logger.Infof("Persistence initialized")

	// Initialize CPS Action Guard (role_id + action_name authorization with TTL cache)
	mid.InitCPSActionGuard(persistence.CPSActionApproveIndexPersistence, persistence.BPSActionApproveIndexPersistence, persistence.RolePersistence, persistence.CPSActionRolePersistence, 5*time.Minute, logger)
	oracleDB := InitOracle(cfg.OracleConnectionString, logger)
	logger.Infof("Oracle database initialized")

	logger.Infof("Initializing Oracle DB client...")
	OraclePersistence := InitOraclePersistence(mongoClient, oracleDB, cfg, *clientOrchestrationProducer, accessListSegmentationProducer, redisRepository, logger)
	logger.Infof("Oracle DB client initialized")
	logger.Infof("Oracle DB client initialized", cfg.OracleConnectionString)

	// Account blocks live in Oracle; persistence was initialized with nil — wire before any service uses it (e.g. BPS user create).
	persistence.AccountBlockPersistence = OraclePersistence.AccountBlock

	logger.Infof("Initializing SMS service...")
	smsService := lib.InitNotificationStore(logger, cfg, notificationProducer)
	logger.Infof("SMS service initialized")

	auth_client, err := client.NewAuthGRPCClient(cfg.CPSAuthSvcGrpcAddress, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC client for auth: %v", err)
	} else {
		logger.Infof("grpc is live and running at:%s", cfg.CPSAuthSvcGrpcAddress)
	}
	defer auth_client.Close()

	sitotagRPCClient, err := api.NewSitotagRPCClient(ctx, logger, cfg.CommonSvcGrpcAddress)
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC client for sitota: %v", err)
	}

	defer local.DisconnectMongo(ctx, mongoClient, logger)

	logger.Infof("initialize service layer")

	serviceLayer := InitServiceLayer(mongoClient, persistence, OraclePersistence, coreInterface, logger, sitotagRPCClient, cfg, minioClient, redisRepository, smsService, notificationProducer, clientOrchestrationProducer, presignClient, queueInfra.Manager, sharedRedis, zapLogger)

	go func() {
		if err := InitFeedbackConsumer(serviceLayer.Feedback, cfg, logger); err != nil {
			logger.Errorf("Failed to start feedback consumer: %v", err)
		}
	}()

	logger.Infof("initialize handler layer")
	handlerLayer := InitHandler(serviceLayer, logger, queueInfra.Manager)

	r := chi.NewRouter()

	// InitRoute(ctx, r, handlerLayer, nil, logger, cfg)
	InitRoute(ctx, r, encryptionMiddleware, handlerLayer, auth_client.Client, redisRepository, logger, cfg)

	// wrap the router with OpenTelemetry instrumentation handler
	otlr := telemetry.WrapHandler(r, "cps-action")

	grpcHandlers := server.NewGrpcServer(serviceLayer.Bank, serviceLayer.Wallet, serviceLayer.Services, serviceLayer.Topup, logger)
	srv := server.NewHTTPServer(cfg, otlr)

	grpcServer, lis := server.StartGrpcServer(grpcHandlers, logger, cfg)

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
	// queue system stopped via deferred queueInfra.Manager.Stop()
}
