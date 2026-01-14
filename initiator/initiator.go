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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/hugokessem/coreio/core"
)

func Init(ctx context.Context) {
	done := make(chan struct{})
	logger := utils.NewLogger()

	logger.Infof("Initializing configuration...")
	cfg := InitConfig(logger)
	logger.Infof("Configuration initialized")

	coreConfig := core.CBECoreCredential{
		Username: cfg.CbeCoreUsername,
		Password: cfg.CbeCorePassword,
		Url:      cfg.CbeCoreUrl,
	}

	// Initialize OpenTelemetry Tracing using platform/telemetry package
	logger.Infof("Initializing OpenTelemetry Tracing...")

	// Build OTEL config from a dedicated helper that reads env vars and provides defaults.
	otelCfg := NewOtelConfig()

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

	logger.Infof("initializing kafka")
	notificationProducer, clientOrchestrationProducer := InitKafkaService(cfg, logger) //27G - 24G= 3G
	logger.Infof("kafka initialized")

	logger.Infof("Initializing persistence...")
	notificationApi := "https://devcbe.eaglelionsystems.com/api/v1.0/chatbirrapi/ldapnotif/sms/send"
	// merchantApi := "https://qaapisuperapp.cbe.com.et/api/v1/cbesuperapp/ecommerce/cps/merchant/"
	// merchantXAPIKey := "0e404061ea76caf9536bc7a38369ca38520aac3c"

	redis := InitRedis(cfg, logger)
	logger.Infof("Initializing redis...")
	redisStorage := InitRedisStorageLayer(redis, logger)
	logger.Infof("redis initialized")

	redisRepository := redisStorage.GetRedisRepository()

	persitence := InitPersistanceLayer(mongoClient, cfg.MongoDBDatabase, coreConfig, notificationApi, *notificationProducer, *clientOrchestrationProducer, redisRepository, cfg, logger)
	logger.Infof("Persistence initialized")

	// Initialize CPS Action Guard (role_id + action_name authorization with TTL cache)
	mid.InitCPSActionGuard(persitence.CPSActionApproveIndexPersistence, 5*time.Minute, logger)

	oracleDB := InitOracle(cfg.OracleConnectionString, logger)
	logger.Infof("Oracle database initialized")

	logger.Infof("Initializing Oracle DB client...")
	OraclePersistence := InitOraclePersistence(oracleDB, logger)
	logger.Infof("Oracle DB client initialized")

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
	serviceLayer := InitServiceLayer(mongoClient, persistence, OraclePersistence, logger, sitotagRPCClient, cfg, minioClient, redisRepository, smsService)

	go func() {
		if err := InitFeedbackConsumer(serviceLayer.Feedback, cfg, logger); err != nil {
			logger.Errorf("Failed to start feedback consumer: %v", err)
		}
	}()

	logger.Infof("initialize handler layer")
	handlerLayer := InitHandler(serviceLayer, logger)

	r := chi.NewRouter()
	// InitRoute(ctx, r, handlerLayer, nil, logger, cfg)
	InitRoute(ctx, r, handlerLayer, auth_client.Client, redisRepository, logger, cfg)

	// wrap the router with OpenTelemetry instrumentation handler
	otlr := telemetry.WrapHandler(r, "cps-action")

	grpcHandlers := server.NewGrpcServer(serviceLayer.Bank, serviceLayer.Wallet, serviceLayer.Services, serviceLayer.Topup, logger)
	srv := server.NewHTTPServer(cfg, otlr)

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
