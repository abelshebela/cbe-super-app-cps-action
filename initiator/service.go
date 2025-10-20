package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/service"
	accountblock "cbe-super-app-cps-action/internal/service/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/service/account_validation"
	advert "cbe-super-app-cps-action/internal/service/ad"
	amount_based_auth "cbe-super-app-cps-action/internal/service/amount_based_auth"
	bankService "cbe-super-app-cps-action/internal/service/bank"
	"cbe-super-app-cps-action/internal/service/media"
	"cbe-super-app-cps-action/internal/storage"

	// bankvault "cbe-super-app-cps-action/internal/service/bankvault"
	bpsService "cbe-super-app-cps-action/internal/service/bps_user"
	"cbe-super-app-cps-action/internal/service/budget"
	bulk_service "cbe-super-app-cps-action/internal/service/bulk"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	cpsusersvc "cbe-super-app-cps-action/internal/service/cps_user"
	customer "cbe-super-app-cps-action/internal/service/customer"
	"cbe-super-app-cps-action/internal/service/department"
	"cbe-super-app-cps-action/internal/service/event"
	miniapp "cbe-super-app-cps-action/internal/service/mini_app"

	// vaultGroupCategory "cbe-super-app-cps-action/internal/service/vaultgroup_category"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"cbe-super-app-cps-action/pkgs/keygen"

	avatar "cbe-super-app-cps-action/internal/service/avatar"
	"cbe-super-app-cps-action/internal/service/fayda"
	feedback "cbe-super-app-cps-action/internal/service/feedback"
	"cbe-super-app-cps-action/internal/service/hq"
	mini_app_merchant "cbe-super-app-cps-action/internal/service/mini_app_merchant"
	password "cbe-super-app-cps-action/internal/service/password_rule"
	permission "cbe-super-app-cps-action/internal/service/permission"
	portalcard "cbe-super-app-cps-action/internal/service/portal_card"

	donation "cbe-super-app-cps-action/internal/service/donation"
	donation_category "cbe-super-app-cps-action/internal/service/donation_category"
	donation_company "cbe-super-app-cps-action/internal/service/donation_company"
	service_details "cbe-super-app-cps-action/internal/service/service_details"

	"cbe-super-app-cps-action/internal/service/productcode"

	"cbe-super-app-cps-action/internal/service/notification"
	"cbe-super-app-cps-action/internal/service/topup"
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
	Topup             service.TopupService
	MiniAppMerchant   service.MiniAppMerchantService
	AccountBlock      service.AccountBlockService
	Department        service.DepartmentService
	PasswordRule      service.PasswordRuleService
	HQService         service.HQService
	MiniAppService    service.MiniAppService
	Fayda             service.FaydaAccountService
	Avatar            service.AvatarService
	AmountBasedAuth   service.AmountBasedAuthService

	Permission        service.PermissionService
	CPSUser           service.CPSUserService
	ServiceDetails    service.ServiceService
	AccountValidation service.AccountValidationService
	ProductCode       service.ProductCodeService

	DonationCategory service.DonationCategoryService

	DonationCompany service.DonationCompanyService

	Donation service.DonationService

	NotificationService service.NotificationService
	// BankVault           service.BankVaultService
	// VaultGroupCategory  service.VaultGroupCategoryService
	ArticleService         service.ArticleService
	ArticleCategoryService service.ArticleCategoryService
	ShortVideoService      service.ShortVideoService
}

var advertBucketName = "advert-bucket" // TODO: Add to config

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface, accountLookupService account_lookup.Account, redis storage.RedisRepository) ServiceLayer {
	// func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface, accountLookupService account_lookup.Account, oracle OraclePersistence) ServiceLayer {
	const minioPubUrl = "https://assetscbedev.eaglelionsystems.com"
	// Create CPS action service with the dispatcher
	// cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)

	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, logger)
	productService := productcode.NewProductCodeService(persistence.ProductCodePersistence, nil, logger) // cpsActionService updefine
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	// miniAppMerchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, cpsActionService, logger)
	avatarService := avatar.NewAvatarService(persistence.AvatarPersistence, nil, logger, minioClient, "avatar", minioPubUrl)

	accountValidation := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, nil, logger)
	eventService := event.NewEventService(persistence.EventPersistence, nil, nil, persistence.UserPersistence, minioClient, minioPubUrl, "events", cfg, logger) // Will be updated after CPS action service is created
	bulkService := bulk_service.NewBulkService(persistence.BulkService, nil, logger)                                                                            // Will be updated after CPS action service is created
	customerService := customer.NewCustomerService(persistence.CustomerService, nil, nil, nil, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, nil, minioClient, minioPubUrl, cfg, "banks")                                               // Will be updated after CPS action service is created
	walletService := wallet.NewWalletService(persistence.WalletPersistence, nil, minioClient, minioPubUrl, "wallets", cfg, logger)                                             // Will be updated after CPS action service is created
	topupService := topup.NewTopupService(persistence.TopupPersistence, nil, minioClient, minioPubUrl, "topups", cfg, logger)                                                  // Will be updated after CPS action service is created
	accountBlockService := accountblock.NewAccountService(persistence.AccountBlockPersistence, nil)                                                                            // Will be updated after CPS action service is created
	departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, nil, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger) // Will be updated after CPS action service is created
	passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, nil, logger)                                                                           // Will be updated after CPS action service is created
	hqService := hq.NewHQService(persistence.HQPersistence, nil, logger)                                                                                                       // Will be updated after CPS action service is created
	keygenService := keygen.NewKeyGenerator(logger, cfg)

	miniAppMerchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, nil, persistence.MiniAppPersistence, logger)
	miniAppService := miniapp.NewMiniAppService(persistence.MiniAppPersistence, nil, miniAppMerchantService, persistence.UserPersistence, keygenService, minioClient, minioPubUrl, "miniapps", cfg, logger)
	faydaService := fayda.NewFaydaService(persistence.FaydaPersistence, nil, logger)

	adService := advert.NewAdvertService(persistence.AdvertRepositoryPersistence, nil, minioClient, "", advertBucketName, cfg, logger)

	serviceDetails := service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, nil, logger) // Will be updated after CPS action service is created

	donationCategoryService := donation_category.NewDonationCategoryService(mongoClient, persistence.DonationCategoryPersistence, nil, logger, minioClient, "donation_icon", cfg, minioPubUrl)
	donationCompanyService := donation_company.NewDonationCompanyService(mongoClient, persistence.DonationCompanyPersistence, nil, logger, minioClient, "donation_company_logo", cfg, minioPubUrl, accountLookupService)
	donationService := donation.NewDonationService(mongoClient, persistence.DonationPersistence, persistence.DonationCategoryPersistence, persistence.DonationCompanyPersistence, nil, logger, minioClient, "donation", cfg, minioPubUrl)

	// productService := productcode.NewProductCodeService(persistence.ProductCodePersistence, nil, logger)

	permissionService := permission.InitPermissionService(
		persistence.PermissionPersistence,
		nil, // Will be updated after CPS action service is created
		logger,
	)

	budgetService := budget.NewBudgetService(persistence.IconPersistence, persistence.ColorPersistence, nil, "budget", minioClient, minioPubUrl, cfg, logger) // Will be updated after CPS action service is created
	cpsUserService := cpsusersvc.NewCPSUserService(
		persistence.CpsUserPersistence,
		persistence.DepartmentPersistence, // temporary it will replaced by department repo
		permissionService,
		nil, // Will be updated after CPS action service is created
		logger,
	)
	notificationsvc := notification.InitNotificationService(persistence.NotificationPersistence, logger, nil)
	amountBased := amount_based_auth.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, nil, minioClient, "amount_based_auth", cfg, logger)
	cpsactionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger, cpsaction.Dispatcher{}) // Will be updated after dispatcher is created

	unlinkService := unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, nil, logger)
	bpsUserService := bpsService.NewBPSUserService(persistence.BPSUserPersistence, nil, logger)
	articleService := media.NewMediaService(persistence.ArticlePersistence, logger)
	articleCategoryService := media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	ShortVideoService := media.NewShortVideoService(persistence.ShortVideoPersistence, logger)
	// Create the service container with all services
	serviceContainer := service.ServiceContainer{

		EventContainer:           eventService,
		FeedbackContainer:        feedbackService,
		UnlinkContainer:          unlinkService,  // Will be updated after CPS action service is created
		BPSUserContainer:         bpsUserService, // Will be updated after CPS action service is created
		CPSActionContainer:       cpsactionService,
		AdContainer:              adService,
		PortalCardContainer:      portalCardService,
		ServiceCheckContainer:    serviceDetails,
		BankContainer:            bank_service,
		WalletContainer:          walletService,
		TopupContainer:           topupService,
		PasswordRuleContainer:    passwordRule,
		AccountBlockContainer:    accountBlockService,
		DepartmentContainer:      departmentService,
		HQContainer:              hqService,
		MiniAppContainer:         miniAppService,
		MiniAppMerchantContainer: miniAppMerchantService,
		FaydaContainer:           faydaService,
		BulkServiceContainer:     bulkService,
		CustomerContainer:        customerService,
		PermissionContainer:      permissionService,
		CPSUserContainer:         cpsUserService,
		BudgetContainer:          budgetService, // Will be updated after CPS action service is created
		AccountContainer:         accountValidation,

		AmountBasedAuthContainer: amountBased,   // Not implemented yet
		AvatarDomian:             avatarService, // Not implemented yet
		BudgetCategoryContainer:  nil,           // Not implemented yet
		NotificationService:      notificationsvc,

		ProductCodeService:        productService, // Not implemented yet
		DonationContainer:         donationService,
		Unlink:                    unlinkService,
		DonationCategoryContainer: donationCategoryService,

		DonationCompanyContainer:   donationCompanyService,
		ArticleContainer:           articleService,
		ArticleCategoryContainer:   articleCategoryService,
		ShortVideoServiceContainer: ShortVideoService,
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
	accountValidation = accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, cpsActionService, logger)
	serviceContainer.AccountContainer = accountValidation
	bank_service = bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, minioPubUrl, cfg, "banks")
	serviceContainer.BankContainer = bank_service

	walletService = wallet.NewWalletService(persistence.WalletPersistence, cpsActionService, minioClient, minioPubUrl, "wallets", cfg, logger)
	serviceContainer.WalletContainer = walletService

	TopupService := topup.NewTopupService(persistence.TopupPersistence, cpsActionService, minioClient, minioPubUrl, "topups", cfg, logger)
	serviceContainer.TopupContainer = TopupService
	accountBlockService = accountblock.NewAccountService(persistence.AccountBlockPersistence, cpsActionService)
	serviceContainer.AccountBlockContainer = accountBlockService

	departmentService = department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	serviceContainer.DepartmentContainer = departmentService

	passwordRule = password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	serviceContainer.PasswordRuleContainer = passwordRule

	hqService = hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	serviceContainer.HQContainer = hqService

	notificationsvc = notification.InitNotificationService(persistence.NotificationPersistence, logger, cpsActionService)
	serviceContainer.NotificationService = notificationsvc

	miniAppService = miniapp.NewMiniAppService(persistence.MiniAppPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, keygenService, minioClient, minioPubUrl, "miniapps", cfg, logger)
	serviceContainer.MiniAppContainer = miniAppService

	faydaService = fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)
	serviceContainer.FaydaContainer = faydaService
	donationCategoryService = donation_category.NewDonationCategoryService(mongoClient, persistence.DonationCategoryPersistence, cpsActionService, logger, minioClient, "donation", cfg, minioPubUrl)
	serviceContainer.DonationCategoryContainer = donationCategoryService
	donationCompanyService = donation_company.NewDonationCompanyService(mongoClient, persistence.DonationCompanyPersistence, cpsActionService, logger, minioClient, "donation", cfg, minioPubUrl, accountLookupService)
	serviceContainer.DonationCompanyContainer = donationCompanyService
	donationService = donation.NewDonationService(mongoClient, persistence.DonationPersistence, persistence.DonationCategoryPersistence, persistence.DonationCompanyPersistence, cpsActionService, logger, minioClient, "donation", cfg, minioPubUrl)
	serviceContainer.DonationContainer = donationService
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

	serviceContainer.ProductCodeService = productcode.NewProductCodeService(persistence.ProductCodePersistence, nil, logger)

	serviceContainer.BudgetContainer = budget.NewBudgetService(persistence.IconPersistence, persistence.ColorPersistence, cpsActionService, "budget", minioClient, minioPubUrl, cfg, logger)

	serviceContainer.AdContainer = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, minioPubUrl, advertBucketName, cfg, logger)
	serviceContainer.AvatarDomian = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, "avatar", minioPubUrl)

	amountBased = amount_based_auth.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, cpsActionService, minioClient, "amount_based_auth", cfg, logger)
	serviceContainer.AmountBasedAuthContainer = amountBased

	// bankVaultSvc := bankvault.NewBankVaultService(oracle.BankVault, cpsActionService, logger)
	// serviceContainer.BankVaultContainer = bankVaultSvc

	// vaultGroupSvc := vaultGroupCategory.NewVaultGroupCategoryService(oracle.vaultGroupCategory, cpsActionService, logger)
	// serviceContainer.VaultGroupCategoryContainer = vaultGroupSvc

	// Rebuild dispatcher with the fully wired container so approvals route to BankVault
	dispatcher = cpsaction.NewDispatcher(serviceContainer)
	cpsActionService = cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger, *dispatcher)
	serviceContainer.CPSActionContainer = cpsActionService
	// Update miniAppMerchantService with the CPS action service
	miniAppMerchantService = mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, cpsActionService, persistence.MiniAppPersistence, logger)
	serviceContainer.MiniAppMerchantContainer = miniAppMerchantService
	serviceContainer.ProductCodeService = productcode.NewProductCodeService(persistence.ProductCodePersistence, cpsActionService, logger)
	unlinkService = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger)
	serviceContainer.Unlink = unlinkService
	articleService = media.NewMediaService(persistence.ArticlePersistence, logger)
	serviceContainer.ArticleContainer = articleService
	articleCategoryService = media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	serviceContainer.ArticleCategoryContainer = articleCategoryService
	ShortVideoService = media.NewShortVideoService(persistence.ShortVideoPersistence, logger)
	serviceContainer.ShortVideoServiceContainer = ShortVideoService
	customerService = customer.NewCustomerService(persistence.CustomerService, cpsActionService, redis, cfg, logger)

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
		AccountValidation: accountValidation,
		Wallet:            walletService,
		Topup:             TopupService,
		PasswordRule:      passwordRule,
		HQService:         hqService,
		AccountBlock:      accountBlockService,
		MiniAppService:    miniAppService,
		MiniAppMerchant:   miniAppMerchantService,
		Department:        departmentService,
		BulkService:       bulkService,
		CustomerService:   customerService,
		Fayda:             faydaService,
		Permission:        permissionService,
		CPSUser:           cpsUserService,

		ServiceDetails:   serviceDetails,
		Donation:         donationService,
		DonationCategory: donationCategoryService,
		DonationCompany:  donationCompanyService,

		ProductCode:         productService,
		NotificationService: notificationsvc,
		// BankVault:           bankVaultSvc,
		// VaultGroupCategory:  vaultGroupSvc,
		ArticleService:         articleService,
		ArticleCategoryService: articleCategoryService,
		ShortVideoService:      ShortVideoService,
	}
}
