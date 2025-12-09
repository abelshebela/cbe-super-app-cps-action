package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	transactionpb "cbe-super-app-cps-action/grpc/sitota"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/service"
	accountblock "cbe-super-app-cps-action/internal/service/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/service/account_validation"
	advert "cbe-super-app-cps-action/internal/service/ad"
	amount_based_auth "cbe-super-app-cps-action/internal/service/amount_based_auth"
	bankService "cbe-super-app-cps-action/internal/service/bank"
	"cbe-super-app-cps-action/internal/service/media"
	newscategory_service "cbe-super-app-cps-action/internal/service/news_category"
	newstag_service "cbe-super-app-cps-action/internal/service/news_tag"
	"cbe-super-app-cps-action/internal/service/transaction"
	"cbe-super-app-cps-action/internal/storage"
	"time"

	bankvault "cbe-super-app-cps-action/internal/service/bankvault"
	vaultGroupCategory "cbe-super-app-cps-action/internal/service/vaultgroup_category"

	bps_action_role_service "cbe-super-app-cps-action/internal/service/bps_action_role"
	bpsService "cbe-super-app-cps-action/internal/service/bps_user"
	budgetCategorySvc "cbe-super-app-cps-action/internal/service/budget_category"
	bulk_service "cbe-super-app-cps-action/internal/service/bulk"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	cpsusersvc "cbe-super-app-cps-action/internal/service/cps_user"
	customer "cbe-super-app-cps-action/internal/service/customer"
	"cbe-super-app-cps-action/internal/service/department"
	deviceversion "cbe-super-app-cps-action/internal/service/device_version"
	"cbe-super-app-cps-action/internal/service/event"
	miniapp "cbe-super-app-cps-action/internal/service/mini_app"

	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"

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

	"cbe-super-app-cps-action/internal/service/notification"
	"cbe-super-app-cps-action/internal/service/productcode"
	"cbe-super-app-cps-action/internal/service/topup"
	"cbe-super-app-cps-action/internal/service/unlink"
	"cbe-super-app-cps-action/internal/service/wallet"
	"cbe-super-app-cps-action/internal/storage/persistance"

	encryption_service "cbe-super-app-cps-action/internal/service/encryption"
	kycsvc "cbe-super-app-cps-action/internal/service/kyc_verifier"
	sitota_service "cbe-super-app-cps-action/internal/service/sitota"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, oracle OraclePersistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, sitotagRPCClient transactionpb.TransactionServiceClient, cfg *config.VaultConfig, minioClient *s3.Client, redis storage.RedisRepository, smsService *lib.NotificationStore) service.ServiceLayer {

	// Initiate Service Layer
	// Assign variable for minio public url
	minioPubUrl := cfg.MinioPublicEndPoint

	accountLookupAdapter := account_lookup.NewCoreAccountLookupAdapter(persistence.AccountLookup, cfg.CbeCoreUrl, time.Duration(cfg.ServerTimeout))
	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, logger)
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	avatarService := avatar.NewAvatarService(persistence.AvatarPersistence, nil, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	accountValidation := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, nil, logger)
	eventService := event.NewEventService(persistence.EventPersistence, nil, nil, persistence.UserPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	bulkService := bulk_service.NewBulkService(persistence.BulkService, nil, logger)

	customerService := customer.NewCustomerService(persistence.CustomerService, nil, nil, nil, nil, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, nil, minioClient, minioPubUrl, cfg, cfg.S3BucketName)
	walletService := wallet.NewWalletService(persistence.WalletPersistence, nil, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	topupService := topup.NewTopupService(persistence.TopupPersistence, nil, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	accountBlockService := accountblock.NewAccountService(persistence.AccountBlockPersistence, nil)
	departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, nil, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, nil, logger)
	hqService := hq.NewHQService(persistence.HQPersistence, nil, logger)
	miniAppMerchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, nil, persistence.MiniAppPersistence, persistence.MerchantLookup, logger, accountLookupAdapter)
	miniAppService := miniapp.NewMiniAppService(persistence.MiniAppPersistence, nil, miniAppMerchantService, logger)
	faydaService := fayda.NewFaydaService(persistence.FaydaPersistence, nil, logger)

	adService := advert.NewAdvertService(persistence.AdvertRepositoryPersistence, nil, minioClient, cfg.S3BucketName, cfg, logger)
	serviceDetails := service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, nil, logger)
	donationCategoryService := donation_category.NewDonationCategoryService(mongoClient, persistence.DonationCategoryPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	donationCompanyService := donation_company.NewDonationCompanyService(mongoClient, persistence.DonationCompanyPersistence, persistence.DonationPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl, accountLookupAdapter)
	donationService := donation.NewDonationService(mongoClient, persistence.DonationPersistence, persistence.DonationCategoryPersistence, persistence.DonationCompanyPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	kycService := kycsvc.NewKYCVerifierService(mongoClient, persistence.KYCVerifierPersistence, persistence.UserPersistence, accountLookupAdapter, nil, persistence.LinkedAccountPersistence, *cfg, logger)
	permissionService := permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, nil, logger)
	budgetCategoryService := budgetCategorySvc.NewBudgetCategoryService(persistence.BudgetCategoryPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	cpsUserService := cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.DepartmentPersistence, permissionService, nil, logger)
	notificationsvc := notification.InitNotificationService(persistence.NotificationPersistence, logger, nil)
	amountBased := amount_based_auth.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, nil, cfg, logger)
	unlinkService := unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, nil, logger)
	bpsUserService := bpsService.NewBPSUserService(persistence.BPSUserPersistence, nil, logger)
	articleService := media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	articleCategoryService := media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	ShortVideoService := media.NewShortVideoService(persistence.ShortVideoPersistence, redis, logger)
	newsTagService := newstag_service.NewNewsTagService(persistence.NewsTagPersistence, nil, logger)
	newsCategoryService := newscategory_service.NewNewsCategoryService(persistence.NewsCategoryPersistence, nil, logger)
	newsTagsService := media.NewMediaTagsService(persistence.NewsTagsServiceContainer, logger)
	deviceVersionService := deviceversion.NewDeviceVersionService(persistence.DeviceVersionControlPersistence, nil, logger)
	sitotaService := sitota_service.NewSitotaTransactionService(oracle.Sitota, logger)
	transactionService := transaction.NewTransactionService(oracle.Transaction, logger)
	encryptionService := encryption_service.NewEncryptionService(cfg, logger)
	bankVaultProductService := bankvault.NewBankVaultService(oracle.BankVault, nil, logger)
	vaultGroupCategoryService := vaultGroupCategory.NewVaultGroupCategoryService(oracle.vaultGroupCategory, nil, logger, minioClient, minioPubUrl, VaultCategoryBucketName, cfg)
	productCodeService := productcode.NewProductCodeService(persistence.ProductCodePersistence, nil, logger)
	bpsActionRoleService := bps_action_role_service.NewBPSActionRoleService(persistence.BPSActionRolePersistence, persistence.BPSActionApproveIndexPersistence, persistence.RolePersistence, nil, logger)
  	miniAppCategory := miniapp.NewMiniAppCategoryService(persistence.MiniAppCategoryPersistence, logger)

	// Attach Service to Container
	serviceContainer := service.ServiceContainer{
		EventContainer:             eventService,
		FeedbackContainer:          feedbackService,
		UnlinkContainer:            unlinkService,
		BPSUserContainer:           bpsUserService,
		CPSActionContainer:         nil,
		AdContainer:                adService,
		PortalCardContainer:        portalCardService,
		ServiceCheckContainer:      serviceDetails,
		BankContainer:              bank_service,
		WalletContainer:            walletService,
		TopupContainer:             topupService,
		PasswordRuleContainer:      passwordRule,
		AccountBlockContainer:      accountBlockService,
		DepartmentContainer:        departmentService,
		HQContainer:                hqService,
		MiniAppContainer:           miniAppService,
		MiniAppMerchantContainer:   miniAppMerchantService,
		FaydaContainer:             faydaService,
		BulkServiceContainer:       bulkService,
		CustomerContainer:          customerService,
		PermissionContainer:        permissionService,
		CPSUserContainer:           cpsUserService,
		BudgetCategoryContainer:    budgetCategoryService,
		AccountContainer:           accountValidation,
		AmountBasedAuthContainer:   amountBased,
		AvatarDomian:               avatarService,
		NotificationService:        notificationsvc,
		ProductCodeService:         productCodeService,
		DonationContainer:          donationService,
		DonationCategoryContainer:  donationCategoryService,
		DonationCompanyContainer:   donationCompanyService,
		ArticleContainer:           articleService,
		ArticleCategoryContainer:   articleCategoryService,
		ShortVideoServiceContainer: ShortVideoService,
		NewsTagContainer:           newsTagService,
		NewsCategoryContainer:      newsCategoryService,
		SitotaContainer:            sitotaService,
		TransactionContainer:       transactionService,
		KYCVerifierContainer:       kycService,
		DeviceVersionContainer:     deviceVersionService,
		NewsTagsServiceContainer:   newsTagsService,
		EncryptionContainer:        encryptionService,
		BankProductContainer:       bankVaultProductService,
		VaultCategoryContainer:     vaultGroupCategoryService,
		BPSActionRoleContainer:     bpsActionRoleService,
		MiniAppCategoryContainer:   miniAppCategory,
	}

	// CPSActionService Appended
	dispatcher := cpsaction.NewDispatcher(serviceContainer)
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger, *dispatcher)
	eventService = event.NewEventService(persistence.EventPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	bulkService = bulk_service.NewBulkService(persistence.BulkService, cpsActionService, logger)
	customerService = customer.NewCustomerService(persistence.CustomerService, cpsActionService, redis, smsService, cfg, logger)
	bank_service = bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, minioPubUrl, cfg, cfg.S3BucketName)
	walletService = wallet.NewWalletService(persistence.WalletPersistence, cpsActionService, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	topupService = topup.NewTopupService(persistence.TopupPersistence, cpsActionService, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	accountBlockService = accountblock.NewAccountService(persistence.AccountBlockPersistence, cpsActionService)
	departmentService = department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	passwordRule = password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	hqService = hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	miniAppService = miniapp.NewMiniAppService(persistence.MiniAppPersistence, cpsActionService, miniAppMerchantService, logger)
	miniAppMerchantService = mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, cpsActionService, persistence.MiniAppPersistence, persistence.MerchantLookup, logger, accountLookupAdapter)
	faydaService = fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)
	adService = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, cfg.S3BucketName, cfg, logger)
	serviceDetails = service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, cpsActionService, logger)
	donationCategoryService = donation_category.NewDonationCategoryService(mongoClient, persistence.DonationCategoryPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	donationCompanyService = donation_company.NewDonationCompanyService(mongoClient, persistence.DonationCompanyPersistence, persistence.DonationPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl, accountLookupAdapter)
	donationService = donation.NewDonationService(mongoClient, persistence.DonationPersistence, persistence.DonationCategoryPersistence, persistence.DonationCompanyPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	permissionService = permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, cpsActionService, logger)
	serviceContainer.PermissionContainer = permissionService
	accountValidationService := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, cpsActionService, logger)
	serviceContainer.AccountContainer = accountValidationService
	cpsUserService = cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.DepartmentPersistence, permissionService, cpsActionService, logger)
	serviceContainer.CPSUserContainer = cpsUserService
	notificationsvc = notification.InitNotificationService(persistence.NotificationPersistence, logger, cpsActionService)
	deviceVersionService = deviceversion.NewDeviceVersionService(persistence.DeviceVersionControlPersistence, cpsActionService, logger)
	serviceContainer.DeviceVersionContainer = deviceVersionService
	budgetCategoryService = budgetCategorySvc.NewBudgetCategoryService(persistence.BudgetCategoryPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	serviceContainer.UnlinkContainer = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger)
	bpsUserService = bpsService.NewBPSUserService(persistence.BPSUserPersistence, cpsActionService, logger)
	serviceContainer.BPSUserContainer = bpsUserService
	serviceContainer.ProductCodeService = productcode.NewProductCodeService(persistence.ProductCodePersistence, nil, logger)
	serviceContainer.AdContainer = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, cfg.S3BucketName, cfg, logger)
	serviceContainer.AvatarDomian = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	avatarService = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	permissionService = permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, cpsActionService, logger)
	cpsUserService = cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.DepartmentPersistence, permissionService, cpsActionService, logger)
	amountBased = amount_based_auth.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, cpsActionService, cfg, logger)
	serviceContainer.AmountBasedAuthContainer = amountBased
	bpsActionRoleService = bps_action_role_service.NewBPSActionRoleService(persistence.BPSActionRolePersistence, persistence.BPSActionApproveIndexPersistence, persistence.RolePersistence, cpsActionService, logger)
	serviceContainer.BPSActionRoleContainer = bpsActionRoleService

	dispatcher = cpsaction.NewDispatcher(serviceContainer)
	cpsActionService = cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger, *dispatcher)
	serviceContainer.CPSActionContainer = cpsActionService
	miniAppMerchantService = mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, cpsActionService, persistence.MiniAppPersistence, persistence.MerchantLookup, logger, accountLookupAdapter)
	serviceContainer.MiniAppMerchantContainer = miniAppMerchantService
	serviceContainer.ProductCodeService = productcode.NewProductCodeService(persistence.ProductCodePersistence, cpsActionService, logger)
	unlinkService = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger)
	serviceContainer.Unlink = unlinkService
	articleService = media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	serviceContainer.ArticleContainer = articleService
	articleCategoryService = media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	serviceContainer.ArticleCategoryContainer = articleCategoryService
	ShortVideoService = media.NewShortVideoService(persistence.ShortVideoPersistence, redis, logger)
	serviceContainer.ShortVideoServiceContainer = ShortVideoService
	customerService = customer.NewCustomerService(persistence.CustomerService, cpsActionService, redis, smsService, cfg, logger)

	kycService = kycsvc.NewKYCVerifierService(mongoClient, persistence.KYCVerifierPersistence, persistence.UserPersistence, accountLookupAdapter, cpsActionService, persistence.LinkedAccountPersistence, *cfg, logger)
	sitotaService = sitota_service.NewSitotaTransactionService(oracle.Sitota, logger)
	transactionService = transaction.NewTransactionService(oracle.Transaction, logger)
	newsTagService = newstag_service.NewNewsTagService(persistence.NewsTagPersistence, cpsActionService, logger)
	encryptionService = encryption_service.NewEncryptionService(cfg, logger)
	newsCategoryService = newscategory_service.NewNewsCategoryService(persistence.NewsCategoryPersistence, cpsActionService, logger)
	ShortVideoService = media.NewShortVideoService(persistence.ShortVideoPersistence, redis, logger)
	articleService = media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	unlinkService = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger)

	bankVaultProductService = bankvault.NewBankVaultService(oracle.BankVault, cpsActionService, logger)
	vaultGroupCategoryService = vaultGroupCategory.NewVaultGroupCategoryService(oracle.vaultGroupCategory, cpsActionService, logger, minioClient, minioPubUrl, cfg.S3BucketName, cfg)
	productCodeService = productcode.NewProductCodeService(persistence.ProductCodePersistence, cpsActionService, logger)

	return service.ServiceLayer{
		CPSAction:              cpsActionService,
		Feedback:               feedbackService,
		EventService:           eventService,
		Avatar:                 avatarService,
		Advert:                 adService,
		BpsUser:                bpsUserService,
		Bank:                   bank_service,
		Unlink:                 unlinkService,
		BudgetCategory:         budgetCategoryService,
		PortalCard:             portalCardService,
		AmountBasedAuth:        amountBased,
		AccountValidation:      accountValidationService,
		Wallet:                 walletService,
		Topup:                  topupService,
		PasswordRule:           passwordRule,
		HQService:              hqService,
		AccountBlock:           accountBlockService,
		MiniAppService:         miniAppService,
		MiniAppMerchant:        miniAppMerchantService,
		Department:             departmentService,
		BulkService:            bulkService,
		CustomerService:        customerService,
		Fayda:                  faydaService,
		Permission:             permissionService,
		CPSUser:                cpsUserService,
		ServiceDetails:         serviceDetails,
		Donation:               donationService,
		DonationCategory:       donationCategoryService,
		DonationCompany:        donationCompanyService,
		ProductCode:            productCodeService,
		NotificationService:    notificationsvc,
		ArticleService:         articleService,
		ArticleCategoryService: articleCategoryService,
		ShortVideoService:      ShortVideoService,
		NewsTagService:         newsTagService,
		NewsCategoryService:    newsCategoryService,
		Sitota:                 sitotaService,
		TransactionService:     transactionService,
		KYCVerifier:            kycService,
		NewsTagsService:        newsTagsService,
		DeviceVersion:          deviceVersionService,
		Encryption:             encryptionService,
		BankVault:              bankVaultProductService,
		VaultGroupCategory:     vaultGroupCategoryService,
		BPSActionRole:          bpsActionRoleService,
		MiniAppCategory:        miniAppCategory,
	}
}
