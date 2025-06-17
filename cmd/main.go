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
	"github.com/spf13/viper"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	accountvalidation_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/account_validation"
	ad_adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/ad"
	branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/branch_handler"
	bulkservices_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/bulk_service"
	cpsusermaker_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/cps-maker_handler"
	customerhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/customer_handler"
	departmenthandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/department_handler"
	faydaRoutes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/fayda_account"
	feedbackhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/feedback_handler"
	permissionhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/permission_handler"
	portal_card_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/portal_card"
	service_details_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/service_details"
	unlinkDeviceHandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink_device_handler"
	adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound"
	cpsusermaker_persistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound"
	accountvalidation_persistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_validation"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/ad"
	branch_repo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/branch"
	customerPersistance "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/customer"
	departmentPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/department"
	faydaaccount "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/fayda_account"
	feedbackPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/feedback"
	permissionPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/permission"
	unlink_outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/unlink"
	accountvalidation_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/account_validation"
	ad_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/ad"
	bulkservices_application "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/bulk_services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/customer"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/department"
	faydaHandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/fayda_account"
	feedback "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/feedback"
	authMiddleware "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/permission"
	portal_card_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/portal_card"
	service_details_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/service_details"
	unlinkApp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/unlink"
	accountvalidation_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_validation"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	ad_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/service"
	branch_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bulkcustomer/services"
	cpsusermaker_service "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/cps_user_maker/services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/customer/service"
	departmentService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/department"
	faydaService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	feedbackService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/feedback"
	permissionService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/permission"
	portal_card_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/portal_card"
	service_details_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
	unlinkDomain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/unlink"
)

func main() {
	logger := utils.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config %v", err)
	}
	//
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
		logger.Fatalf("failed to initialize minio clinet", err)
	}

	dbname := cfg.MongoDBDatabase
	//dbname := "ldap_cbs"
	collectionNames := []string{
		"BPSUsers",
		"BPSActions",
		"cps_users",
		"CPSServices",
		"Member",
		"linked_accounts",
		"mini_app",
		"hq_services",
	}

	// r := chi.NewRouter()
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

	branchRepo := branch_repo.NewBranchPersistence(
		mongoClient,
		cfg.MongoDBDatabase,
		[]string{
			"branches",
			"cps_actions",
		},
		viper.GetDuration("timeout"),
		logger,
	)
	branchService := branch_domain.NewBranchService(branchRepo)
	branchHandler := branch_handler.NewBranchHandler(branchService, logger)
	branch_handler.RegisterBranchRoutes(r, branchHandler, authMddleware)

	customerPersitance := customerPersistance.InitCustomerDetail(mongoClient, cfg.MongoDBDatabase, "customers", logger)
	customerDomain := service.IntiCustomerDomain(customerPersitance, logger)
	customerApp := customer.InitCustomerHandler(customerDomain, logger)
	customerRoutes := customerhandler.NewCustomerHTTPHandler(customerApp, logger)
	customerhandler.InitCustomerRoutes(r, customerRoutes, authMddleware)

	feedbackPersitance := feedbackPersistence.InitFeedback(mongoClient, cfg.MongoDBDatabase, "feedbacks", logger)
	feedbackDomain := feedbackService.InitFeedbackDomain(feedbackPersitance, logger)
	feedbackApp := feedback.InitFeedbackHandler(feedbackDomain, logger)
	feedbackRoutes := feedbackhandler.NewFeedbackHTTPHandler(feedbackApp, logger)
	feedbackhandler.InitFeedbackRoutes(r, feedbackRoutes)

	faydaPersistence := faydaaccount.InitFaydaAccountPersistence(mongoClient, cfg.MongoDBDatabase,
		[]string{"cps_actions", "customers"}, logger)
	faydaDomin := faydaService.InitFaydaAccountDomain(faydaPersistence, logger)
	faydaApp := faydaHandler.InitFaydaHandler(faydaDomin, logger)
	faydaHandlers := faydaRoutes.InitFaydaAdapter(faydaApp, logger)
	faydaRoutes.InitFaydaRoutes(r, faydaHandlers, authMddleware)

	adPersistence := ad.InitAD(mongoClient, cfg.MongoDBDatabase, []string{"adverts", "cps_actions"}, logger)
	adDomain := ad_domain.InitADDomian(adPersistence, logger)
	adHandler := ad_handler.InitADHandler(adDomain, minioClient, "adverts", logger)
	adAdapter := ad_adapter.InitADAdapter(adHandler, logger)
	// ad_adapter.InitADRoutes(r, adAdapter, authMddleware)

	departmentPersistence := departmentPersistence.InitDepartment(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
	departmentDomain := departmentService.InitDepartmentDomain(departmentPersistence, departmentPersistence, logger)
	departmentApp := department.InitDepartmentHandler(departmentDomain, logger)
	departmentRoutes := departmenthandler.NewDepartmentHTTPHandler(departmentApp, logger)
	departmenthandler.InitDepartmentRoutes(r, departmentRoutes, authMddleware)

	permissionPersistence := permissionPersistence.InitPermission(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
	permissionDomain := permissionService.InitPermissionDomain(permissionPersistence, permissionPersistence, permissionPersistence, logger)
	permissionApp := permission.InitPermissionHandler(permissionDomain, logger)
	permissionRoutes := permissionhandler.NewPermissionHTTPHandler(permissionApp, logger)
	permissionhandler.InitPermissionRoutes(r, permissionRoutes, authMddleware)

	unlinkRepo := unlink_outbound.NewUnlinkInfrastructure(mongoClient, cfg.MongoDBDatabase, []string{"user", "otp", "cps_action"}, logger)
	unlinkDeviceService := unlinkDomain.NewUnlinkService(unlinkRepo)
	unlinkDeviceApplication := unlinkApp.NewUnlinkHandler(unlinkDeviceService)
	unlink_handler := unlinkDeviceHandler.NewHTTPUnlinkHandler(unlinkDeviceApplication, logger)
	unlinkDeviceHandler.RegisterHTTPUnlinkRoutes(r, unlink_handler, authMddleware)

	accountvalidation_persistence := accountvalidation_persistence.InitAccountValidationPersistence(
		mongoClient,
		cfg.MongoDBDatabase,
		10*time.Second,
		logger,
	)
	domainAccountValidationService := accountvalidation_domain.NewService(accountvalidation_persistence, adapter_port, logger)

	accountValidationApp := accountvalidation_app.NewApplication(domainAccountValidationService)
	accountValidationHandler := accountvalidation_inbound.NewHttpAccountValidation(accountValidationApp, logger)
	accountvalidation_inbound.InitAccountValidationHandlerMaker(r, accountValidationHandler, authMddleware)

	// Initialize service details using the new implementation
	serviceDetailsStore := adapter.NewServiceDetailsPersistence(mongoClient, cfg.MongoDBDatabase, logger)
	domainServiceDetailsService := service_details_domain.NewServiceStore(serviceDetailsStore, serviceDetailsStore, logger)
	serviceDetailsApp := service_details_app.NewApplication(domainServiceDetailsService, logger)
	serviceDetailsHandler := service_details_inbound.NewHttpServiceDetails(serviceDetailsApp, logger)
	service_details_inbound.InitServiceDetailsRoutes(r, serviceDetailsHandler, authMddleware)

	portalCardStore := adapter.NewPortalCardPersistence(mongoClient, cfg.MongoDBDatabase, logger)
	domainPortalCardService := portal_card_domain.NewPortalCardDomain(portalCardStore, logger)
	portalCardApp := portal_card_app.NewPortalCardApp(domainPortalCardService, logger)
	portalCardHandler := portal_card_inbound.NewportalCardHandler(portalCardApp, logger)
	portal_card_inbound.InitPortalCardRoutes(r, portalCardHandler, authMddleware)

	cpsUserPersistence := cpsusermaker_persistence.NewCPSUserPersistence(
		mongoClient,
		cfg.MongoDBDatabase,
		collectionNames,
	)
	cpsUserService := cpsusermaker_service.NewCPSUserService(cpsUserPersistence)
	cpsUserHandler := cpsusermaker_handler.InitCPSUserMakerHandler(cpsUserService, logger)

	cpsusermaker_handler.RegisterCPSUserMakerRoutes(r, cpsUserHandler, authMddleware)

	ad_adapter.InitADRoutes(r, adAdapter, authMddleware)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	////
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		log.Println("🚀 Server started on", viper.GetString("Port"))
		log.Printf("Server stopped with error: %v\n", server.ListenAndServe())
	}()

	sig := <-quit

	log.Printf("server shutting down with signal: %v\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown gracefully with error %v", err)
	}
	//
	log.Println("Server shutdown successfully")
}
