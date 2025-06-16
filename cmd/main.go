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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/branch_handler"
	bulkservices_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/bulk_service"
	customerhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/customer_handler"
	faydaRoutes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/fayda_account"
	adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound"
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
	departmenthandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/department_handler"
	departmentPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/department"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/department"
	departmentService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/department"
	permissionhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/permission_handler"
	permissionPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/permission"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/permission"
	permissionService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/permission"
	feedbackPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/feedback"
	feedbackService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/feedback"
	feedback "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/feedback"
	feedbackhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/feedback_handler"
    unlinkDeviceHandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink_device_handler"
    unlink_outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/unlink"
	unlinkApp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/unlink"
	unlinkDomain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/unlink"
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

	server := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerPort),
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
