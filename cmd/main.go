package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// "testing/quick"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/amount_based_auth_handler"
	persistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/amount_based_auth"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/avatar"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/amount_based_auth_app"
	amount_based_auth_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/amount_based_auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	accountvalidation_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/account_validation"
	ad_adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/ad"
	bank_adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/bank"
	branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/branch_handler"
	budgethandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/budget_handler"
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
	wallet_adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/wallet"
	adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound"
	cpsusermaker_persistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound"
	accountvalidation_persistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_validation"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/ad"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/bank"
	branch_repo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/branch"
	budgetPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/budget"
	customerPersistance "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/customer"
	departmentPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/department"
	faydaaccount "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/fayda_account"
	feedbackPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/feedback"
	permissionPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/permission"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/wallet"
	unlink_outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/unlink"
	accountvalidation_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/account_validation"
	ad_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/ad"
	bank_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/bank"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/budget"
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
	wallet_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/wallet"
	accountvalidation_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_validation"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	ad_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/service"
	bank_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/service"
	budgetService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/budget"
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
	wallet_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/wallet/service"

	//branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/branch_handler"
	//customerhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/customer_handler"
	//branch_repo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/branch"
	//customerPersistance "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/customer"
	//"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/customer"
	//branch_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/services"
	//"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/customer/service"
	//faydaRoutes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/fayda_account"
	//faydaaccount "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/fayda_account"
	//faydaHandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/fayda_account"
	//faydaService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/fayda_account/service"

	hq_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/hq"
	password_rule_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/password_rule"
	password_rule_routes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/password_rule"
	hq_persistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/persistence/hq"
	application_hq "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/hq"
	hq_service "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/hq"
	password_rule_services "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/password_rule/services"

	avatar_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/avatar"

	avatar_adapter "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http/avatar"
	avatar_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/avatar"
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

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

	dbname := cfg.MongoDBDatabase
	// dbname := "ldap_cbs"
	collectionNames := []string{
		"bps_user",
		"cps_action",
		"cps_users",
		"service",
		"member",
		"linked_accounts",
		"mini_app",

		"portal_card",
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

	adapter_port := adapter.NewOutBoundStore(mongoClient, dbname, collectionNames, logger)
	domain_services := domain.NewService(adapter_port, logger)
	application := bulkservices_application.NewAttachDetachChecker(domain_services, logger)
	handlers := bulkservices_inbound.NewHttpBulkService(application, logger)
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
	ad_adapter.InitADRoutes(r, adAdapter, authMddleware)

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

	budgetPersistence := budgetPersistence.InitBudget(mongoClient, cfg.MongoDBDatabase, []string{"icons", "colors", "cps_actions"}, logger)
	budgetDomain := budgetService.InitBudgetDomain(budgetPersistence, logger)
	budgetApp := budget.InitBudgetHandler(budgetDomain, minioClient, "icons", logger)
	budgetRoutes := budgethandler.NewBudgetHTTPHandler(budgetApp, logger)
	budgethandler.InitBudgetRoutes(r, budgetRoutes, authMddleware)

	accountvalidation_persistence := accountvalidation_persistence.InitAccountValidationPersistence(
		mongoClient,
		cfg.MongoDBDatabase,
		10*time.Second,
		logger,
	)
	domainAccountValidationService := accountvalidation_domain.NewService(accountvalidation_persistence, adapter_port, logger)

	accountValidationApp := accountvalidation_app.NewApplication(domainAccountValidationService, logger)
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
		[]string{
			"cps_users",
			"cps_actions",
			"BPSUsers",
			"CPSServices",
			"portal_cards",
			"validation_rules",
		},
	)
	cpsUserService := cpsusermaker_service.NewCPSUserService(cpsUserPersistence)
	cpsUserHandler := cpsusermaker_handler.InitCPSUserMakerHandler(cpsUserService, logger)
	cpsusermaker_handler.RegisterCPSUserMakerRoutes(r, cpsUserHandler, authMddleware)

	bankPersistence := bank.InitBank(mongoClient, cfg.MongoDBDatabase, []string{"banks", "cps_actions"}, logger)
	bankDomain := bank_domain.InitBankDomain(bankPersistence, minioClient, "banks", logger)
	bankHandler := bank_handler.InitBankHanlder(bankDomain, logger)
	bankAdapter := bank_adapter.InitBankAdapter(bankHandler, logger)
	bank_adapter.InitBankRoutes(r, bankAdapter, authMddleware)

	walletRepo := wallet.InitWallet(mongoClient, cfg.MongoDBDatabase, []string{"wallets", "cps_actions"}, logger)
	walletDomain := wallet_domain.InitWalletDomain(walletRepo, minioClient, "wallets", logger)
	walletHandler := wallet_handler.InitWalletHanlder(walletDomain, logger)
	walletAdapter := wallet_adapter.InitWalletAdapter(walletHandler, logger)
	wallet_adapter.InitWalletRoutes(r, walletAdapter, authMddleware)

	amountBasedAuthRepo := persistence.InitAmountBasedAuth(mongoClient, cfg.MongoDBDatabase, []string{"auth_tier", "cps_action"}, logger)
	amountBasedAuthService := amount_based_auth_domain.NewAmountBasedAuthService(amountBasedAuthRepo, logger)
	amountBasedAuthApplication := amount_based_auth_app.AmountBasedAuthHandler(amountBasedAuthService)
	amountBasedAuthHandler := amount_based_auth_handler.NewAmountBasedAuthHandler(amountBasedAuthApplication, logger)
	amount_based_auth_handler.InitAmountBasedAuthHandler(r, amountBasedAuthHandler, authMddleware)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	passwordRuleOutbound := adapter.NewOutboundPasswordRuleInfra(
		mongoClient,
		cfg.MongoDBDatabase,
		[]string{
			"password_rules",
			"cps_actions",
		},
	)
	passwordRuleService := password_rule_services.NewPasswordRuleService(passwordRuleOutbound)
	passwordRuleHandler := password_rule_handler.NewPasswordRuleHTTPHandler(passwordRuleService, logger)
	password_rule_routes.RegisterPasswordRuleRoutes(r, passwordRuleHandler, authMddleware)
	hq_persistence := hq_persistence.NewHQPersistence(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
	hqService := hq_service.NewService(hq_persistence, adapter_port, logger)
	hqApp := application_hq.NewApplication(hqService, logger)
	hqHandler := hq_handler.NewHQHTTPHandler(hqApp, logger)
	hq_handler.InitHQRoutes(r, hqHandler, authMddleware)

	avatarPersitence := avatar.InitAvatarPersistence(mongoClient, cfg.MongoDBDatabase, []string{"cps_actions", "avatares"}, logger)
	avatarDomain := avatar_domain.InitAvatarDomain(avatarPersitence, minioClient, "avatares", logger)
	avatarApp := avatar_app.InitAvatarAPP(avatarDomain, logger)
	avatarHanler := avatar_adapter.InitAvatarHTTPHandler(avatarApp, logger)
	avatar_adapter.InitAvatarRoutes(r, avatarHanler, authMddleware)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		log.Println("🚀 Server started on", viper.GetString("Port"))
		log.Printf("Server stopped with error: %v\n", server.ListenAndServe())
	}()

	sig := <-quit

	log.Printf("server shutting down with signal: %v\n", sig)

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown gracefully with error %v", err)
	}
	//
	log.Println("Server shutdown successfully")
}
