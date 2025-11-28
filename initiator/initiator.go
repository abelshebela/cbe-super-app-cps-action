package initiator

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"cbe-super-app-cps-action/cmd/client"
	"cbe-super-app-cps-action/cmd/server"
	local "cbe-super-app-cps-action/config"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/storage/api"

	"cbe-super-app-cps-action/platform/telemetry"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"log"

	"gitlab.com/yohannesteshome/coreio/core"

	"github.com/go-chi/chi/v5"
)

func Init(ctx context.Context) {
	done := make(chan struct{})
	coreConfig := core.CBECoreCredential{
		Username: "SUPERAPP",
		Password: "123456",
		Url:      "http://10.1.15.195:8080/CBESUPERAPPV2/services?wsdl=null",
	}
	logger := utils.NewLogger()
	logger.Infof("Initializing configuration...")
	cfg := InitConfig(logger)
	logger.Infof("Configuration initialized")

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
	}()

	logger.Infof("Initializing MongoDB client...")
	mongoClient := InitMongo(cfg.MongoDBURI, logger)
	logger.Infof("MongoDB client initialized")

	logger.Infof("Initializing Minio client...")
	minioClient := InitMinio(*cfg, logger)
	logger.Infof("Minio client initialized")

	logger.Infof("initializing kafka")
	kafkaInit := InitKafkaService(cfg, logger)
	logger.Infof("kafka initialized")

	logger.Infof("Initializing persistence...")
	notificationApi := "https://devcbe.eaglelionsystems.com/api/v1.0/chatbirrapi/ldapnotif/sms/send"
	merchantApi := "https://devcbe.eaglelionsystems.com/api/v1.0/chatbirrapi/ldapnotif/sms/send"
	persitence := InitPersistanceLayer(mongoClient, cfg.MongoDBDatabase, coreConfig, merchantApi, notificationApi, *kafkaInit, cfg, logger)
	logger.Infof("Persistence initialized")

	redis := InitRedis(cfg, logger)
	logger.Infof("Initializing redis...")
	redisStorage := InitRedisStorageLayer(redis, logger)
	logger.Infof("redis initialized")

	redisRepository := redisStorage.GetRedisRepository()

	oracleDB := InitOracle(cfg.OracleConnectionString, logger)
	logger.Infof("Oracle database initialized")

	logger.Infof("Initializing Oracle DB client...")
	OraclePersistence := InitOraclePersistence(oracleDB, logger)
	logger.Infof("Oracle DB client initialized")

	logger.Infof("Initializing SMS service...")
	smsService := lib.InitNotificationStore(logger, cfg, kafkaInit)
	logger.Infof("SMS service initialized")

	sessionGRPCClient, clientStore, err := api.NewSessionGRPCClient(cfg.CommonSvcGrpcAddress, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC session client: %v", err)
	}
	defer func() {
		if cerr := clientStore.Close(); cerr != nil {
			logger.Errorf("error closing client store: %v", cerr)
		}
	}()

	auth_client, err := client.NewAuthGRPCClient(cfg.CPSAuthSvcGrpcAddress, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC client for auth: %v", err)
	} else {
		logger.Infof("grpc is live and running at:%s", cfg.CPSAuthSvcGrpcAddress)
	}
	defer auth_client.Close()

	sitotagRPCClient, err := api.NewSitotagRPCClient(ctx, logger, cfg.CbeToCbeGrpcAddress)
	if err != nil {
		logger.Fatalf("Failed to initialize gRPC client for sitota: %v", err)
	}
	if cerr := sitotagRPCClient.Close(); cerr != nil {
		logger.Errorf("Failed to close sitota RPC client: %v", cerr)
	}

	defer local.DisconnectMongo(ctx, mongoClient, logger)

	logger.Infof("initialize service layer")
	serviceLayer := InitServiceLayer(mongoClient, persitence, OraclePersistence, logger, sessionGRPCClient, sitotagRPCClient, cfg, minioClient, redisRepository, smsService)

	go func() {
		if err := InitFeedbackConsumer(serviceLayer.Feedback, cfg, logger); err != nil {
			logger.Errorf("Failed to start feedback consumer: %v", err)
		}
	}()

	logger.Infof("initialize handler layer")
	handlerLayer := InitHandler(serviceLayer, logger)

	r := chi.NewRouter()
	InitRoute(ctx, r, handlerLayer, auth_client.Client, logger, cfg)

	// wrap the router with OpenTelemetry instrumentation handler
	otlr := telemetry.WrapHandler(r, "http-server")

	go func() {
		fmt.Println("Goroutines: ", runtime.NumGoroutine())
	}()

	grpcHandlers := server.NewGrpcServer(serviceLayer.Bank, serviceLayer.Wallet, serviceLayer.ServiceDetails, serviceLayer.Topup, logger)
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
