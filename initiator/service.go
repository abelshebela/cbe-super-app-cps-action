package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/service"
	accountblock "cbe-super-app-cps-action/internal/service/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/service/account_validation"
	advert "cbe-super-app-cps-action/internal/service/ad"
	amount_based_auth "cbe-super-app-cps-action/internal/service/amount_based_auth"
	bankService "cbe-super-app-cps-action/internal/service/bank"
	bpsService "cbe-super-app-cps-action/internal/service/bps_user"
	"cbe-super-app-cps-action/internal/service/budget"
	bulk_service "cbe-super-app-cps-action/internal/service/bulk"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	cpsusersvc "cbe-super-app-cps-action/internal/service/cps_user"
	customer "cbe-super-app-cps-action/internal/service/customer"
	"cbe-super-app-cps-action/internal/service/department"
	"cbe-super-app-cps-action/internal/service/event"
	miniapp "cbe-super-app-cps-action/internal/service/mini_app"
	"cbe-super-app-cps-action/pkgs/keygen"

	avatar "cbe-super-app-cps-action/internal/service/avatar"
	"cbe-super-app-cps-action/internal/service/fayda"
	feedback "cbe-super-app-cps-action/internal/service/feedback"
	"cbe-super-app-cps-action/internal/service/hq"
	mini_app_merchant "cbe-super-app-cps-action/internal/service/mini_app_merchant"
	password "cbe-super-app-cps-action/internal/service/password_rule"
	permission "cbe-super-app-cps-action/internal/service/permission"
	portalcard "cbe-super-app-cps-action/internal/service/portal_card"

	service_details "cbe-super-app-cps-action/internal/service/service_details"

	"cbe-super-app-cps-action/internal/service/productcode"

	"cbe-super-app-cps-action/internal/service/unlink"
	"cbe-super-app-cps-action/internal/service/wallet"
	"cbe-super-app-cps-action/internal/storage/persistance"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceLayer struct {
	EventService    service.EventService
	BulkService     service.BulkService
	CustomerService service.CustomerService
	CPSAction       service.CPSActionService
	Feedback        service.FeedbackService
	// Services  service.ServiceContainer
	Unlink            service.UnlinkService
	BpsUser           service.BPSUserService
	Budget            service.BudgetService
	Bank              service.BankService
	PortalCard        service.PortalCardService
	Advert            service.AdvertService
	ValidationService service.AccountValidationService
	Wallet            service.WalletService
	MiniAppMerchant   service.MiniAppMerchantService
	AccountBlock      service.AccountBlockService
	Department        service.DepartmentService
	PasswordRule      service.PasswordRuleService
	HQService         service.HQService
	MiniAppService    service.MiniAppService
	Fayda             service.FaydaAccountService
	Avatar            service.AvatarService

	Permission     service.PermissionService
	CPSUser        service.CPSUserService
	ServiceDetails service.ServiceService

	ProductCode service.ProductCodeService

	AmountBasedAuth service.AmountBasedAuthService
}

var advertBucketName = "advert-bucket" // TODO: Add to config

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {
	const minioPubUrl = "https://assetscbedev.eaglelionsystems.com"
	// Create CPS action service with the dispatcher
	// cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)

	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, logger)
	productService := productcode.NewProductCodeService(persistence.ProductCodePersistence, nil, logger) // cpsActionService updefine
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	// miniAppMerchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, cpsActionService, logger)
	avatarService := avatar.NewAvatarService(persistence.AvatarPersistence, nil, logger, minioClient, "avatar", minioPubUrl)

	// customerSerice := customer.NewCustomerService(persistence.CustomerService, logger)
	// bank_service := bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, minioPubUrl, cfg, "banks")
	// walletService := wallet.NewWalletService(persistence.WalletPersistence, cpsActionService, minioClient, "wallets", cfg, logger)
	// accountBlockService := accountblock.NewAccountService(persistence.AccountBlockPersistence, cpsActionService)
	// departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	// passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	// hqService := hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	// keygenService := keygen.NewKeyGenerator(logger, cfg)

	// miniAppService := miniapp.NewMiniAppService(persistence.MiniAppPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, keygenService, minioClient, minioPubUrl, "miniapps", cfg, logger)
	// fayda := fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)

	// serviceDetails := service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, cpsActionService, logger)
	accountValidation := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, nil, logger)
	eventService := event.NewEventService(persistence.EventPersistence, nil, nil, persistence.UserPersistence, minioClient, minioPubUrl, "events", cfg, logger) // Will be updated after CPS action service is created
	bulkService := bulk_service.NewBulkService(persistence.BulkService, nil, logger)                                                                            // Will be updated after CPS action service is created
	customerSerice := customer.NewCustomerService(persistence.CustomerService, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, nil, minioClient, minioPubUrl, cfg, "banks")                                               // Will be updated after CPS action service is created
	walletService := wallet.NewWalletService(persistence.WalletPersistence, nil, minioClient, minioPubUrl, "wallets", cfg, logger)                                             // Will be updated after CPS action service is created
	accountBlockService := accountblock.NewAccountService(persistence.AccountBlockPersistence, nil)                                                                            // Will be updated after CPS action service is created
	departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, nil, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger) // Will be updated after CPS action service is created
	passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, nil, logger)                                                                           // Will be updated after CPS action service is created
	hqService := hq.NewHQService(persistence.HQPersistence, nil, logger)                                                                                                       // Will be updated after CPS action service is created
	keygenService := keygen.NewKeyGenerator(logger, cfg)

	miniAppMerchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, nil, logger)                                                                              // Will be updated after CPS action service is created
	miniAppService := miniapp.NewMiniAppService(persistence.MiniAppPersistence, nil, miniAppMerchantService, persistence.UserPersistence, keygenService, minioClient, minioPubUrl, "miniapps", cfg, logger) // Will be updated after CPS action service is created
	faydaService := fayda.NewFaydaService(persistence.FaydaPersistence, nil, logger)                                                                                                                        // Will be updated after CPS action service is created

	adService := advert.NewAdvertService(persistence.AdvertRepositoryPersistence, nil, minioClient, "", advertBucketName, cfg, logger)

	serviceDetails := service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, nil, logger) // Will be updated after CPS action service is created
	// productService := productcode.NewProductCodeService(persistence.ProductCodePersistence, nil, logger)

	permissionService := permission.InitPermissionService(
		persistence.PermissionPersistence,
		nil, // Will be updated after CPS action service is created
		logger,
	)

	cpsUserService := cpsusersvc.NewCPSUserService(
		persistence.CpsUserPersistence,
		persistence.DepartmentPersistence, // temporary it will replaced by department repo
		permissionService,
		nil, // Will be updated after CPS action service is created
		logger,
	)

	unlinkService := unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, nil, logger)
	// Create the service container with all services
	serviceContainer := service.ServiceContainer{
		EventContainer:           eventService,
		FeedbackContainer:        feedbackService,
		UnlinkContainer:          nil, // Will be updated after CPS action service is created
		BPSUserContainer:         nil, // Will be updated after CPS action service is created
		AdContainer:              adService,
		PortalCardContainer:      portalCardService,
		ServiceCheckContainer:    serviceDetails,
		BankContainer:            bank_service,
		WalletContainer:          walletService,
		PasswordRuleContainer:    passwordRule,
		AccountBlockContainer:    accountBlockService,
		DepartmentContainer:      departmentService,
		HQContainer:              hqService,
		MiniAppContainer:         miniAppService,
		MiniAppMerchantContainer: miniAppMerchantService,
		FaydaContainer:           faydaService,
		BulkServiceContainer:     bulkService,
		CustomerContainer:        customerSerice,
		PermissionContainer:      permissionService,
		CPSUserContainer:         cpsUserService,
		BudgetContainer:          nil, // Will be updated after CPS action service is created
		AccountContainer:         accountValidation,
		AmountBasedAuthContainer: nil,            // Not implemented yet
		AvatarDomian:             avatarService,  // Not implemented yet
		BudgetCategoryContainer:  nil,            // Not implemented yet
		NotificationService:      nil,            // Not implemented yet
		ProductCodeService:       productService, // Not implemented yet
		DonationContainer:        nil,            // Not implemented yet
		Unlink:                   unlinkService,
	}

	// Create the dispatcher with the service container
	dispatcher := cpsaction.NewDispatcher(serviceContainer)

	// Create CPS action service with the dispatcher
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger, *dispatcher)
	productService = productcode.NewProductCodeService(persistence.ProductCodePersistence, cpsActionService, logger)

	// Now update all services that need the CPS action service
	eventService = event.NewEventService(persistence.EventPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, minioClient, minioPubUrl, "events", cfg, logger)
	serviceContainer.EventContainer = eventService
	avatarService = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, "avatar", minioPubUrl)
	bulkService = bulk_service.NewBulkService(persistence.BulkService, cpsActionService, logger)
	serviceContainer.BulkServiceContainer = bulkService

	bank_service = bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, minioPubUrl, cfg, "banks")
	serviceContainer.BankContainer = bank_service

	walletService = wallet.NewWalletService(persistence.WalletPersistence, cpsActionService, minioClient, minioPubUrl, "wallets", cfg, logger)
	serviceContainer.WalletContainer = walletService

	accountBlockService = accountblock.NewAccountService(persistence.AccountBlockPersistence, cpsActionService)
	serviceContainer.AccountBlockContainer = accountBlockService

	departmentService = department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	serviceContainer.DepartmentContainer = departmentService

	passwordRule = password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	serviceContainer.PasswordRuleContainer = passwordRule

	hqService = hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	serviceContainer.HQContainer = hqService

	miniAppService = miniapp.NewMiniAppService(persistence.MiniAppPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, keygenService, minioClient, minioPubUrl, "miniapps", cfg, logger)
	serviceContainer.MiniAppContainer = miniAppService

	faydaService = fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)
	serviceContainer.FaydaContainer = faydaService

	serviceDetails = service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, cpsActionService, logger)
	serviceContainer.ServiceCheckContainer = serviceDetails

	permissionService = permission.InitPermissionService(
		persistence.PermissionPersistence,
		cpsActionService,
		logger,
	)
	serviceContainer.PermissionContainer = permissionService

	cpsUserService = cpsusersvc.NewCPSUserService(
		persistence.CpsUserPersistence,
		persistence.DepartmentPersistence, // Use actual department repository
		permissionService,
		cpsActionService,
		logger,
	)
	serviceContainer.CPSUserContainer = cpsUserService

	serviceContainer.BudgetContainer = budget.NewBudgetService(persistence.IconPersistence, persistence.ColorPersistence, cpsActionService, "budget", minioClient, minioPubUrl, cfg, logger)

	serviceContainer.UnlinkContainer = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger)

	serviceContainer.BPSUserContainer = bpsService.NewBPSUserService(persistence.BPSUserPersistence, cpsActionService, logger)

	serviceContainer.AdContainer = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, minioPubUrl, advertBucketName, cfg, logger)
	serviceContainer.AvatarDomian = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, "avatar", minioPubUrl)
	// Update miniAppMerchantService with the CPS action service
	miniAppMerchantService = mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, cpsActionService, logger)
	serviceContainer.MiniAppMerchantContainer = miniAppMerchantService

	serviceContainer.ProductCodeService = productcode.NewProductCodeService(persistence.ProductCodePersistence, cpsActionService, logger)
	serviceContainer.Unlink = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger)
	return ServiceLayer{
		CPSAction: cpsActionService,

		Feedback:          feedbackService,
		EventService:      eventService,
		Avatar:            avatarService,
		Advert:            advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, minioPubUrl, advertBucketName, cfg, logger),
		BpsUser:           bpsService.NewBPSUserService(persistence.BPSUserPersistence, cpsActionService, logger),
		Bank:              bank_service,
		Unlink:            unlinkService,
		Budget:            budget.NewBudgetService(persistence.IconPersistence, persistence.ColorPersistence, cpsActionService, "budget", minioClient, minioPubUrl, cfg, logger),
		PortalCard:        portalCardService,
		AmountBasedAuth:   amount_based_auth.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, cpsActionService, minioClient, "amount_based_auth", cfg, logger),
		ValidationService: accountValidation,
		Wallet:            walletService,
		PasswordRule:      passwordRule,
		HQService:         hqService,
		AccountBlock:      accountBlockService,
		MiniAppService:    miniAppService,
		MiniAppMerchant:   miniAppMerchantService,
		Department:        departmentService,
		BulkService:       bulkService,
		CustomerService:   customerSerice,
		Fayda:             faydaService,
		Permission:        permissionService,
		CPSUser:           cpsUserService,

		ServiceDetails: serviceDetails,

		ProductCode: productService,
	}
}
