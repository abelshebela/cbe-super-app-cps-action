package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/branch_handler"
	bulkservices_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/bulk_service"
	customerhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/customer_handler"
	faydaRoutes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/fayda_account"
	adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/ad"
	branch_repo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/branch"
	customerPersistance "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/customer"
	faydaaccount "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/fayda_account"
	bulkservices_application "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/bulk_services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/customer"
	faydaHandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/fayda_account"
	authMiddleware "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	branch_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bulkcustomer/services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/customer/service"
	faydaService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	ad_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/service"
	ad_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/ad"
	ad_adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/ad"

)

func main() {
	logger := utils.NewLogger()
	defer logger.Sync()

	cfg, err := config.Load()
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

	minioClient,err := config.NewMinioClient(&config.VaultConfig{
		MinioEndPoint: cfg.MinioEndPoint,
		MinioAccessKey: cfg.MinioAccessKey,
		MinioSecretKey: cfg.MinioSecretKey,
	})
	if err != nil{
		logger.Fatalf("failed to initialize minio clinet",err)
	}

	dbname := cfg.MongoDBDatabase
	//dbname := "ldap_cbs"
	collectionNames := []string{
		"BPSActions",
		"BPSUsers",
		"CPSServices",
		"Member",
		"linked_accounts",
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	authMddleware := authMiddleware.InitAuthMiddleware(cfg.JwtSecretKey, cfg.Key, cfg.IV, logger)

	adapter_port := adapter.NewOutBoundStore(mongoClient, dbname, collectionNames)
	domain_services := domain.NewService(adapter_port)
	application := bulkservices_application.NewAttachDetachChecker(domain_services)
	handlers := bulkservices_inbound.NewHttpBulkService(application)
	bulkservices_inbound.InitServiceHandlerMaker(r, handlers, authMddleware)

	branchPersistence := branch_repo.NewBranchPersistence(mongoClient, cfg.MongoDBDatabase,
		"branches", "cps_actions", logger)
	branchService := branch_domain.NewBranchService(branchPersistence)
	branchHandler := branch_handler.NewBranchHandler(branchService, logger)
	branch_handler.RegisterBranchRoutes(r, branchHandler, authMddleware)

	customerPersitance := customerPersistance.InitCustomerDetail(mongoClient, cfg.MongoDBDatabase, "customers", logger)
	customerDomain := service.IntiCustomerDomain(customerPersitance, logger)
	customerApp := customer.InitCustomerHandler(customerDomain, logger)
	customerRoutes := customerhandler.NewCustomerHTTPHandler(customerApp, logger)
	customerhandler.InitCustomerRoutes(r, customerRoutes, authMddleware)

	faydaPersistence := faydaaccount.InitFaydaAccountPersistence(mongoClient, cfg.MongoDBDatabase,
		"cps_actions", "customers_new", logger)
	faydaDomin := faydaService.InitFaydaAccountDomain(faydaPersistence, logger)
	faydaApp := faydaHandler.InitFaydaHandler(faydaDomin, logger)
	faydaHandlers := faydaRoutes.InitFaydaAdapter(faydaApp, logger)
	faydaRoutes.InitFaydaRoutes(r, faydaHandlers, authMddleware)

	adPersistence := ad.InitAD(mongoClient,cfg.MongoDBDatabase,[]string{"adverts","cps_actions"},logger)
	adDomain := ad_domain.InitADDomian(adPersistence,logger)
	adHandler := ad_handler.InitADHandler(adDomain,minioClient,"adverts",logger)
	adAdapter := ad_adapter.InitADAdapter(adHandler,logger)
	ad_adapter.InitADRoutes(r,adAdapter,authMddleware)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("🚀 Server started on", server.Addr)
		log.Printf("Server stopped with error: %v\n", server.ListenAndServe())
	}()

	sig := <-quit
	log.Printf("server shutting down with signal: %v\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown gracefully with error %v", err)
	}

	log.Println("Server shutdown successfully")
}
