package initiator

import (
	transactionpb "cbe-super-app-cps-action/grpc/sitota"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/service"
	access_list_segmentation_service "cbe-super-app-cps-action/internal/service/access_list_segmentation"
	accountblock "cbe-super-app-cps-action/internal/service/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/service/account_validation"
	advert "cbe-super-app-cps-action/internal/service/ad"
	amount_based_auth "cbe-super-app-cps-action/internal/service/amount_based_auth"
	bankService "cbe-super-app-cps-action/internal/service/bank"
	bps_action_service "cbe-super-app-cps-action/internal/service/bps_action"
	bpsuser "cbe-super-app-cps-action/internal/service/bps_user"
	event_merchant_service "cbe-super-app-cps-action/internal/service/event_merchant"
	logistics_merchant_service "cbe-super-app-cps-action/internal/service/logistics_merchant"
	"cbe-super-app-cps-action/internal/service/media"
	newscategory_service "cbe-super-app-cps-action/internal/service/news_category"
	newstag_service "cbe-super-app-cps-action/internal/service/news_tag"
	"cbe-super-app-cps-action/internal/service/transaction"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"time"

	bankvault "cbe-super-app-cps-action/internal/service/bankvault"
	vault_category "cbe-super-app-cps-action/internal/service/vault"

	bps_action_role_service "cbe-super-app-cps-action/internal/service/bps_action_role"
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
	ecommerce_merchant "cbe-super-app-cps-action/internal/service/ecommerce-merchant"
	"cbe-super-app-cps-action/internal/service/fayda"
	feedback "cbe-super-app-cps-action/internal/service/feedback"
	"cbe-super-app-cps-action/internal/service/hq"
	password "cbe-super-app-cps-action/internal/service/password_rule"
	permission "cbe-super-app-cps-action/internal/service/permission"
	portalcard "cbe-super-app-cps-action/internal/service/portal_card"

	cps_action_role_service "cbe-super-app-cps-action/internal/service/cps_action_role"
	cps_role "cbe-super-app-cps-action/internal/service/cps_roles"
	kyc_service "cbe-super-app-cps-action/internal/service/customer_kyc"
	customer_segmentation "cbe-super-app-cps-action/internal/service/customer_segmentation"
	donation "cbe-super-app-cps-action/internal/service/donation"
	donation_category "cbe-super-app-cps-action/internal/service/donation_category"
	donation_company "cbe-super-app-cps-action/internal/service/donation_company"
	encryption_service "cbe-super-app-cps-action/internal/service/encryption"
	job_role "cbe-super-app-cps-action/internal/service/job_role"
	kycsvc "cbe-super-app-cps-action/internal/service/kyc_verifier"
	"cbe-super-app-cps-action/internal/service/notification"
	"cbe-super-app-cps-action/internal/service/roles"
	services_svc "cbe-super-app-cps-action/internal/service/services"
	sitota_service "cbe-super-app-cps-action/internal/service/sitota"
	"cbe-super-app-cps-action/internal/service/topup"
	"cbe-super-app-cps-action/internal/service/unlink"
	ussd_merchant "cbe-super-app-cps-action/internal/service/ussd_merchant"
	vault_amount_tier "cbe-super-app-cps-action/internal/service/vault_amount_tier"
	"cbe-super-app-cps-action/internal/service/wallet"
	"cbe-super-app-cps-action/internal/storage/persistance"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, oracle OraclePersistence, logger utils.Logger, sitotagRPCClient transactionpb.TransactionServiceClient, cfg *config.VaultConfig, minioClient *s3.Client, redis storage.RedisRepository, smsService *lib.NotificationStore, clientOrchestrationProducer *kafka.ClientOrchestrationProducer,presignClient *s3.PresignClient) service.ServiceLayer {

	// Initiate Service Layer
	// Assign variable for minio public url
	minioPubUrl := cfg.MinioPublicEndPoint

	mediaProducer := media.CreateKafkaProducer(logger, cfg)
	accountLookupAdapter := account_lookup.NewCoreAccountLookupAdapter(persistence.AccountLookup, cfg.CbeCoreUrl, time.Duration(cfg.ServerTimeout))
	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, persistence.CustomerService, logger)
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	avatarService := avatar.NewAvatarService(persistence.AvatarPersistence, nil, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	accountValidation := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, nil, logger)
	eventService := event.NewEventService(persistence.EventPersistence, nil, nil, persistence.UserPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	bulkService := bulk_service.NewBulkService(persistence.BulkService, nil, persistence.AccessListSegmentationPersistence, logger)
	ussdMerchant := ussd_merchant.NewUssdMerchantService(persistence.UssdMerchantPersistence, persistence.ServicesPersistence, minioClient, nil, cfg.S3BucketName, *cfg, logger)

	customerService := customer.NewCustomerService(persistence.CustomerService, persistence.BpsActionPersistence, nil, nil, nil, nil, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, nil, minioClient, minioPubUrl, cfg, cfg.S3BucketName)
	walletService := wallet.NewWalletService(persistence.WalletPersistence, nil, persistence.ServicesPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	bpsActionService := bps_action_service.NewBPSActionService(persistence.BPSActionRolePersistence, persistence.BpsActionPersistence, logger, bps_action_service.Dispatcher{})
	topupService := topup.NewTopupService(persistence.TopupPersistence, nil, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	accountBlockService := accountblock.NewAccountService(persistence.AccountBlockPersistence, nil, logger)
	departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, nil, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, nil, logger)
	hqService := hq.NewHQService(persistence.HQPersistence, nil, logger)
	// miniAppMerchantService := ecommerce_merchant.NewEcommerceMerchantService(persistence.EcommerceMerchantPersistence, nil, persistence.MerchantLookup, logger, persistence.Account_lookup_external)
	miniMerchant := miniapp.NewMiniAppMerchantService(persistence.MiniAppMerchant, persistence.MiniAppPersistence, persistence.MerchantLookup, logger)
	miniAppService := miniapp.NewMiniAppService(persistence.MiniAppPersistence, miniMerchant, logger)
	ecommerceMerchantService := ecommerce_merchant.NewEcommerceMerchantService(persistence.EcommerceMerchantPersistence, nil, persistence.MerchantLookup, logger, persistence.Account_lookup_external, *cfg)
	faydaService := fayda.NewFaydaService(persistence.FaydaPersistence, nil, logger)

	adService := advert.NewAdvertService(persistence.AdvertRepositoryPersistence, nil, minioClient, cfg.S3BucketName, cfg, logger)
	// serviceDetails := service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, nil, logger)
	donationCategoryService := donation_category.NewDonationCategoryService(mongoClient, persistence.DonationCategoryPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	donationCompanyService := donation_company.NewDonationCompanyService(mongoClient, persistence.DonationCompanyPersistence, persistence.DonationPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl, accountLookupAdapter)
	donationService := donation.NewDonationService(mongoClient, persistence.DonationPersistence, persistence.DonationCategoryPersistence, persistence.DonationCompanyPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	kycService := kycsvc.NewKYCVerifierService(mongoClient, persistence.KYCVerifierPersistence, persistence.UserPersistence, accountLookupAdapter, nil, persistence.LinkedAccountPersistence, *cfg, logger)
	permissionService := permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, nil, logger)
	budgetCategoryService := budgetCategorySvc.NewBudgetCategoryService(persistence.BudgetCategoryPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	cpsUserService := cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.RolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.DepartmentPersistence, permissionService, nil, nil, logger)
	notificationsvc := notification.InitNotificationService(persistence.NotificationPersistence, logger, nil, smsService)
	amountBased := amount_based_auth.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, nil, cfg, logger)
	unlinkService := unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, persistence.AccountBlockPersistence, nil, logger)
	bpsUserService := bpsuser.NewBPSUserService(persistence.BPSUserPersistence, persistence.RolePersistence, nil, persistence.CpsUserPersistence, logger)
	articleService := media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	articleCategoryService := media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	ShortVideoService := media.NewShortVideoService(persistence.ShortVideoPersistence, redis, mediaProducer, logger)
	newsTagService := newstag_service.NewNewsTagService(persistence.NewsTagPersistence, nil, logger)
	newsCategoryService := newscategory_service.NewNewsCategoryService(persistence.NewsCategoryPersistence, nil, logger)
	newsTagsService := media.NewMediaTagsService(persistence.NewsTagsServiceContainer, logger)
	deviceVersionService := deviceversion.NewDeviceVersionService(persistence.DeviceVersionControlPersistence, nil, clientOrchestrationProducer, logger)
	sitotaService := sitota_service.NewSitotaTransactionService(oracle.Sitota, logger)
	transactionService := transaction.NewTransactionService(oracle.Transaction, logger)
	encryptionService := encryption_service.NewEncryptionService(cfg, logger)
	bankVaultProductService := bankvault.NewBankVaultService(oracle.BankVault, nil, logger)
	vaultCategoryService := vault_category.NewVaultCategoryService(oracle.vaultCategory, nil, logger, minioClient, minioPubUrl, VaultCategoryBucketName, cfg)
	bpsActionRoleService := bps_action_role_service.NewBPSActionRoleService(persistence.BPSActionRolePersistence, persistence.BPSActionApproveIndexPersistence, persistence.JobRolePersistence, nil, logger)
	miniAppCategory := miniapp.NewMiniAppCategoryService(persistence.MiniAppCategoryPersistence, logger)
	cpsActionRoleService := cps_action_role_service.NewCPSActionRoleService(persistence.CPSActionRolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.JobRolePersistence, nil, logger)
	eventMerchantService := event_merchant_service.NewEventMerchantService(persistence.EventMerchantPersistence, nil, accountLookupAdapter, cfg, logger)
	logisticsMerchantService := logistics_merchant_service.NewLogisticsMerchantService(persistence.LogisticsMerchantPersistence, nil, accountLookupAdapter, cfg, logger)
	servicesService := services_svc.NewServicesService(persistence.ServicesPersistence, nil, logger)
	vaultAmountTierSrv := vault_amount_tier.NewVaultAmountTierService(oracle.VaultAmountTier, nil, logger)
	miniAppProductCodeContainer := miniapp.NewMiniAppProductCodeService(persistence.MiniAppProductCodePersistence, logger)
	accessListSegmentationService := access_list_segmentation_service.NewAccessListSegmentationService(persistence.AccessListSegmentationPersistence, nil, persistence.AccessListPersistence, persistence.AccountBlockPersistence, persistence.CustomerService, persistence.CPSRoles, logger)
	customerSegmentationService := customer_segmentation.NewCustomerSegmentation(persistence.CustomerSegmentation, persistence.CPSRoles, nil, logger)
	jobRoleService := job_role.NewJobRoleService(persistence.JobRolePersistence, persistence.RolePersistence, nil, *cfg, logger)
	RoleService := roles.NewRoleService(persistence.JobRolePersistence, persistence.PortalCardPersistence, persistence.CPSActionApproveIndexPersistence, nil, *cfg, logger)
	CPSRolesService := cps_role.NewCPSRoleService(persistence.CPSRoles, nil, logger)
	customerKYCService := kyc_service.NewCustomerKYCService(persistence.CustomerKYCPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)

	// Attach Service to Container
	serviceContainer := service.ServiceContainer{
		RoleContainer:       RoleService,
		JobRoleContainer:    jobRoleService,
		EventContainer:      eventService,
		FeedbackContainer:   feedbackService,
		UnlinkContainer:     unlinkService,
		BPSUserContainer:    bpsUserService,
		CPSActionContainer:  nil,
		AdContainer:         adService,
		PortalCardContainer: portalCardService,
		// ServiceCheckContainer:         serviceDetails,
		BankContainer:                     bank_service,
		WalletContainer:                   walletService,
		TopupContainer:                    topupService,
		PasswordRuleContainer:             passwordRule,
		AccountBlockContainer:             accountBlockService,
		DepartmentContainer:               departmentService,
		HQContainer:                       hqService,
		MiniAppContainer:                  miniAppService,
		FaydaContainer:                    faydaService,
		BulkServiceContainer:              bulkService,
		CustomerContainer:                 customerService,
		PermissionContainer:               permissionService,
		CPSUserContainer:                  cpsUserService,
		BudgetCategoryContainer:           budgetCategoryService,
		AccountContainer:                  accountValidation,
		AmountBasedAuthContainer:          amountBased,
		AvatarDomian:                      avatarService,
		NotificationService:               notificationsvc,
		DonationContainer:                 donationService,
		DonationCategoryContainer:         donationCategoryService,
		DonationCompanyContainer:          donationCompanyService,
		ArticleContainer:                  articleService,
		ArticleCategoryContainer:          articleCategoryService,
		ShortVideoServiceContainer:        ShortVideoService,
		NewsTagContainer:                  newsTagService,
		NewsCategoryContainer:             newsCategoryService,
		SitotaContainer:                   sitotaService,
		TransactionContainer:              transactionService,
		KYCVerifierContainer:              kycService,
		DeviceVersionContainer:            deviceVersionService,
		NewsTagsServiceContainer:          newsTagsService,
		EncryptionContainer:               encryptionService,
		BankProductContainer:              bankVaultProductService,
		VaultCategoryContainer:            vaultCategoryService,
		BPSActionRoleContainer:            bpsActionRoleService,
		MiniAppCategoryContainer:          miniAppCategory,
		CPSActionRoleContainer:            cpsActionRoleService,
		EventMerchantServiceContainer:     eventMerchantService,
		LogisticsMerchantServiceContainer: logisticsMerchantService,
		ServiceContainer:                  servicesService,
		VaultAmountTierContainer:          vaultAmountTierSrv,
		MiniAppProductCodeContainer:       miniAppProductCodeContainer,
		AccessListSegmentationContainer:   accessListSegmentationService,
		CustomerSegmentationContainer:     customerSegmentationService,
		MiniAppMerchantContainer:          miniMerchant,
		EcommerceMerchantContainer:        ecommerceMerchantService,
		CPSRolesContainer:                 CPSRolesService,
		CustomerKYCContainer:              customerKYCService,
		UssdMerchantContainer:             ussdMerchant,
	}

	// CPSActionService Appended
	dispatcher := cpsaction.NewDispatcher(serviceContainer)
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSActionRolePersistence, persistence.CPSAction, logger, *dispatcher,minioClient,cfg.S3BucketName,cfg.MinioEndPoint,presignClient)
	cpsActionService = cpsaction.WithActionRolePolicy(cpsActionService, persistence.CPSActionRolePersistence, logger)
	ussdMerchant = ussd_merchant.NewUssdMerchantService(persistence.UssdMerchantPersistence, persistence.ServicesPersistence, minioClient, cpsActionService, cfg.S3BucketName, *cfg, logger)

	// eventService = event.NewEventService(persistence.EventPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	// cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger, *dispatcher)
	eventService = event.NewEventService(persistence.EventPersistence, cpsActionService, ecommerceMerchantService, persistence.UserPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	bulkService = bulk_service.NewBulkService(persistence.BulkService, cpsActionService, persistence.AccessListSegmentationPersistence, logger)
	// customerService = customer.NewCustomerService(persistence.CustomerService, persistence.BpsActionPersistence, cpsActionService, redis, smsService, cfg, logger)
	bank_service = bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, minioPubUrl, cfg, cfg.S3BucketName)
	walletService = wallet.NewWalletService(persistence.WalletPersistence, cpsActionService, persistence.ServicesPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	topupService = topup.NewTopupService(persistence.TopupPersistence, cpsActionService, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	accountBlockService = accountblock.NewAccountService(persistence.AccountBlockPersistence, cpsActionService, logger)
	departmentService = department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	passwordRule = password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	hqService = hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	miniAppService = miniapp.NewMiniAppService(persistence.MiniAppPersistence, miniMerchant, logger)
	// miniAppMerchantService = ecommerce_merchant.NewEcommerceMerchantService(persistence.EcommerceMerchantPersistence, cpsActionService, persistence.MerchantLookup, logger, accountLookupAdapter)
	// ecommerceMerchantService = ecommerce_merchant.NewEcommerceMerchantService(persistence.EcommerceMerchantPersistence, cpsActionService, persistence.MerchantLookup, logger, accountLookupAdapter, *cfg)
	faydaService = fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)
	adService = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, cfg.S3BucketName, cfg, logger)
	// serviceDetails = service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, cpsActionService, logger)
	donationCategoryService = donation_category.NewDonationCategoryService(mongoClient, persistence.DonationCategoryPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	donationCompanyService = donation_company.NewDonationCompanyService(mongoClient, persistence.DonationCompanyPersistence, persistence.DonationPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl, accountLookupAdapter)
	donationService = donation.NewDonationService(mongoClient, persistence.DonationPersistence, persistence.DonationCategoryPersistence, persistence.DonationCompanyPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	permissionService = permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, cpsActionService, logger)
	serviceContainer.PermissionContainer = permissionService
	accountValidationService := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, cpsActionService, logger)
	serviceContainer.AccountContainer = accountValidationService
	cpsUserService = cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.RolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.DepartmentPersistence, permissionService, cpsActionService, persistence.BPSUserPersistence, logger)
	serviceContainer.CPSUserContainer = cpsUserService
	notificationsvc = notification.InitNotificationService(persistence.NotificationPersistence, logger, cpsActionService, smsService)
	deviceVersionService = deviceversion.NewDeviceVersionService(persistence.DeviceVersionControlPersistence, cpsActionService, clientOrchestrationProducer, logger)
	serviceContainer.DeviceVersionContainer = deviceVersionService
	budgetCategoryService = budgetCategorySvc.NewBudgetCategoryService(persistence.BudgetCategoryPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	serviceContainer.UnlinkContainer = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, persistence.AccountBlockPersistence, cpsActionService, logger)
	bpsUserService = bpsuser.NewBPSUserService(persistence.BPSUserPersistence, persistence.RolePersistence, cpsActionService, persistence.CpsUserPersistence, logger)
	serviceContainer.BPSUserContainer = bpsUserService
	serviceContainer.AdContainer = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, cfg.S3BucketName, cfg, logger)
	serviceContainer.AvatarDomian = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	avatarService = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	permissionService = permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, cpsActionService, logger)
	cpsUserService = cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.RolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.DepartmentPersistence, permissionService, cpsActionService, persistence.BPSUserPersistence, logger)
	amountBased = amount_based_auth.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, cpsActionService, cfg, logger)
	serviceContainer.AmountBasedAuthContainer = amountBased
	bpsActionRoleService = bps_action_role_service.NewBPSActionRoleService(persistence.BPSActionRolePersistence, persistence.BPSActionApproveIndexPersistence, persistence.JobRolePersistence, cpsActionService, logger)
	serviceContainer.BPSActionRoleContainer = bpsActionRoleService
	cpsActionRoleService = cps_action_role_service.NewCPSActionRoleService(persistence.CPSActionRolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.JobRolePersistence, cpsActionService, logger)
	serviceContainer.CPSActionRoleContainer = cpsActionRoleService
	dispatcher = cpsaction.NewDispatcher(serviceContainer)
	cpsActionService = cpsaction.NewCPSActionService(persistence.CPSActionRolePersistence, persistence.CPSAction, logger, *dispatcher,minioClient,cfg.S3BucketName,cfg.MinioEndPoint,presignClient)
	cpsActionService = cpsaction.WithActionRolePolicy(cpsActionService, persistence.CPSActionRolePersistence, logger)
	serviceContainer.CPSActionContainer = cpsActionService

	// Services catalog service (uses CPSAction for maker-checker)
	servicesService = services_svc.NewServicesService(persistence.ServicesPersistence, cpsActionService, logger)
	// serviceContainer.ServicesContainer = servicesService
	// miniAppMerchantService = ecommerce_merchant.NewEcommerceMerchantService(persistence.EcommerceMerchantPersistence, cpsActionService, persistence.MerchantLookup, logger, accountLookupAdapter)
	serviceContainer.MiniAppMerchantContainer = miniMerchant
	ecommerceMerchantService = ecommerce_merchant.NewEcommerceMerchantService(persistence.EcommerceMerchantPersistence, cpsActionService, persistence.MerchantLookup, logger, accountLookupAdapter, *cfg)
	unlinkService = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, persistence.AccountBlockPersistence, cpsActionService, logger)
	serviceContainer.Unlink = unlinkService
	articleService = media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	serviceContainer.ArticleContainer = articleService
	articleCategoryService = media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	serviceContainer.ArticleCategoryContainer = articleCategoryService
	ShortVideoService = media.NewShortVideoService(persistence.ShortVideoPersistence, redis, mediaProducer, logger)
	serviceContainer.ShortVideoServiceContainer = ShortVideoService
	customerService = customer.NewCustomerService(persistence.CustomerService, persistence.BpsActionPersistence, cpsActionService, redis, smsService, cfg, logger)

	kycService = kycsvc.NewKYCVerifierService(mongoClient, persistence.KYCVerifierPersistence, persistence.UserPersistence, accountLookupAdapter, cpsActionService, persistence.LinkedAccountPersistence, *cfg, logger)
	sitotaService = sitota_service.NewSitotaTransactionService(oracle.Sitota, logger)
	transactionService = transaction.NewTransactionService(oracle.Transaction, logger)
	newsTagService = newstag_service.NewNewsTagService(persistence.NewsTagPersistence, cpsActionService, logger)
	encryptionService = encryption_service.NewEncryptionService(cfg, logger)
	newsCategoryService = newscategory_service.NewNewsCategoryService(persistence.NewsCategoryPersistence, cpsActionService, logger)
	ShortVideoService = media.NewShortVideoService(persistence.ShortVideoPersistence, redis, mediaProducer, logger)
	articleService = media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	unlinkService = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, persistence.AccountBlockPersistence, cpsActionService, logger)
	bankVaultProductService = bankvault.NewBankVaultService(oracle.BankVault, cpsActionService, logger)
	vaultCategoryService = vault_category.NewVaultCategoryService(oracle.vaultCategory, cpsActionService, logger, minioClient, minioPubUrl, cfg.S3BucketName, cfg)
	eventMerchantService = event_merchant_service.NewEventMerchantService(persistence.EventMerchantPersistence, cpsActionService, accountLookupAdapter, cfg, logger)

	logisticsMerchantService = logistics_merchant_service.NewLogisticsMerchantService(persistence.LogisticsMerchantPersistence, cpsActionService, accountLookupAdapter, cfg, logger)

	serviceContainer.EventMerchantServiceContainer = eventMerchantService
	serviceContainer.LogisticsMerchantServiceContainer = logisticsMerchantService
	vaultAmountTierSrv = vault_amount_tier.NewVaultAmountTierService(oracle.VaultAmountTier, cpsActionService, logger)
	accessListSegmentationService = access_list_segmentation_service.NewAccessListSegmentationService(persistence.AccessListSegmentationPersistence, cpsActionService, persistence.AccessListPersistence, persistence.AccountBlockPersistence, persistence.CustomerService, persistence.CPSRoles, logger)
	customerSegmentationService = customer_segmentation.NewCustomerSegmentation(persistence.CustomerSegmentation, persistence.CPSRoles, cpsActionService, logger)
	jobRoleService = job_role.NewJobRoleService(persistence.JobRolePersistence, persistence.RolePersistence, cpsActionService, *cfg, logger)
	RoleService = roles.NewRoleService(persistence.JobRolePersistence, persistence.PortalCardPersistence, persistence.CPSActionApproveIndexPersistence, cpsActionService, *cfg, logger)
	CPSRolesService = cps_role.NewCPSRoleService(persistence.CPSRoles, cpsActionService, logger)
	customerKYCService = kyc_service.NewCustomerKYCService(persistence.CustomerKYCPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)

	return service.ServiceLayer{
		RoleService:       RoleService,
		JobRoleService:    jobRoleService,
		CPSAction:         cpsActionService,
		Feedback:          feedbackService,
		EventService:      eventService,
		Avatar:            avatarService,
		Advert:            adService,
		BpsUser:           bpsUserService,
		Bank:              bank_service,
		Unlink:            unlinkService,
		BudgetCategory:    budgetCategoryService,
		PortalCard:        portalCardService,
		AmountBasedAuth:   amountBased,
		AccountValidation: accountValidationService,
		Wallet:            walletService,
		Topup:             topupService,
		PasswordRule:      passwordRule,
		HQService:         hqService,
		AccountBlock:      accountBlockService,
		MiniAppService:    miniAppService,
		EcommerceMerchant: ecommerceMerchantService,
		Department:        departmentService,
		BulkService:       bulkService,
		CustomerService:   customerService,
		Fayda:             faydaService,
		Permission:        permissionService,
		CPSUser:           cpsUserService,
		// ServiceDetails:         serviceDetails,
		Services:                      servicesService,
		Donation:                      donationService,
		DonationCategory:              donationCategoryService,
		DonationCompany:               donationCompanyService,
		NotificationService:           notificationsvc,
		ArticleService:                articleService,
		ArticleCategoryService:        articleCategoryService,
		ShortVideoService:             ShortVideoService,
		NewsTagService:                newsTagService,
		NewsCategoryService:           newsCategoryService,
		Sitota:                        sitotaService,
		TransactionService:            transactionService,
		KYCVerifier:                   kycService,
		NewsTagsService:               newsTagsService,
		DeviceVersion:                 deviceVersionService,
		Encryption:                    encryptionService,
		BankVault:                     bankVaultProductService,
		VaultCategoryService:          vaultCategoryService,
		BPSActionRole:                 bpsActionRoleService,
		MiniAppCategory:               miniAppCategory,
		CPSActionRole:                 cpsActionRoleService,
		EventMerchantService:          eventMerchantService,
		LogisticsMerchantService:      logisticsMerchantService,
		VaultAmountTierService:        vaultAmountTierSrv,
		MiniappProductCode:            miniAppProductCodeContainer,
		AccessListSegmentationService: accessListSegmentationService,
		CustomerSegmentation:          customerSegmentationService,
		CPSRoles:                      CPSRolesService,
		CustomerKYC:                   customerKYCService,
		UssdMerchantService:           ussdMerchant,
		BPSActionService:              bpsActionService,
	}
}
