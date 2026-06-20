package initiator

import (
	transactionpb "cbe-super-app-cps-action/grpc/sitota"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/service"
	access_list_segmentation_service "cbe-super-app-cps-action/internal/service/access_list_segmentation"
	accountblock "cbe-super-app-cps-action/internal/service/account_block"
	ap_svc "cbe-super-app-cps-action/internal/service/account_product"
	apc_svc "cbe-super-app-cps-action/internal/service/account_product_category"
	account_sub_type_svc "cbe-super-app-cps-action/internal/service/account_sub_type"
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
	role_delegation_service "cbe-super-app-cps-action/internal/service/role_delegation"
	tac_svc "cbe-super-app-cps-action/internal/service/term_and_condition"
	token_provider_svc "cbe-super-app-cps-action/internal/service/token_provider"
	"cbe-super-app-cps-action/internal/storage"
	tp_client "cbe-super-app-cps-action/internal/storage/external_call/token_provider"
	"cbe-super-app-cps-action/internal/storage/kafka"
	queue "cbe-super-app-cps-action/internal/storage/queue_system"
	"time"

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
	customer_group "cbe-super-app-cps-action/internal/service/customer_group"
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
	superapp_role "cbe-super-app-cps-action/internal/service/superapp_role"
	"cbe-super-app-cps-action/internal/service/topup"
	"cbe-super-app-cps-action/internal/service/unlink"
	ussd_merchant "cbe-super-app-cps-action/internal/service/ussd_merchant"
	"cbe-super-app-cps-action/internal/service/wallet"
	"cbe-super-app-cps-action/internal/storage/persistance"


	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hugokessem/coreio/core"
	access_list_cache "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/catch/access_list"
	service_cache "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/catch/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	sharedRedis "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config/redis"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	logge "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	walletCatch "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/catch/wallet"

)

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, oracle OraclePersistence, coreInterface core.CBECoreAPIInterface,
	logger utils.Logger, sitotagRPCClient transactionpb.TransactionServiceClient,
	cfg *config.VaultConfig, minioClient *s3.Client, redis storage.RedisRepository, smsService *lib.NotificationStore,
	clientOrchestrationProducer *kafka.ClientOrchestrationProducer, presignClient *s3.PresignClient, queueManager *queue.QueueManager,
	sharedRedisClient *sharedRedis.RedisClient, cacheLogger logge.Logger) service.ServiceLayer {

	persistence.AccessListSegmentationPersistence.SetAccountBlockRepository(oracle.AccountBlock)

	minioPubUrl := cfg.MinioPublicEndPoint

	mediaProducer := media.CreateKafkaProducer(logger, cfg)
	accountLookupAdapter := account_lookup.NewCoreAccountLookupAdapter(persistence.AccountLookup, cfg.CbeCoreUrl, time.Duration(cfg.ServerTimeout)*time.Second, logger)

	const (
		tokenURL          = "https://devapisuperapp.cbe.com.et/superapp/parser/proxy/cbe-dev/sandbox/oauth-mb-cbebirr/oauth2/token?target=https%3A%2F%2Fapi-gw-uat-gateway-apic-nonprod.apps.cp4itest.cbe.local"
		tokenClientID     = "f1ceebd8d6d5b802dc7fd8332ab33603"
		tokenClientSecret = "05ac01f2f134bfff2549669bc11bd6cc"
		tokenScope        = "mb-cbebirr-scope"
	)
	tokenProviderClient := tp_client.NewTokenProviderClient(tokenURL, tokenClientID, tokenClientSecret, tokenScope, logger)
	tokenProviderService := token_provider_svc.NewTokenProviderService(tokenProviderClient, redis, logger)
	serviceCache := service_cache.NewServiceCatch(*sharedRedisClient, cacheLogger)
	accessListCache := access_list_cache.NewAccessListCatch(*sharedRedisClient, cacheLogger)

	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, oracle.Customer, logger)
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	avatarService := avatar.NewAvatarService(persistence.AvatarPersistence, nil, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	accountValidation := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, nil, logger)
	eventService := event.NewEventService(persistence.EventPersistence, nil, nil, persistence.UserPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	bulkService := bulk_service.NewBulkService(oracle.AccessListOracle, nil, oracle.AccessListSegmentaion, logger)
	ussdMerchant := ussd_merchant.NewUssdMerchantService(persistence.UssdMerchantPersistence, oracle.ServicesPersistence, minioClient, nil, cfg.S3BucketName, *cfg, logger)

	customerService := customer.NewCustomerService(oracle.Customer, persistence.BpsActionPersistence, nil, nil, nil, accountLookupAdapter, nil, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, oracle.BankOracle, nil, minioClient, minioPubUrl, cfg, cfg.S3BucketName)
	accountSubTypeService := account_sub_type_svc.NewAccountSubTypeService(oracle.AccountSubType, nil, logger)
	accountProductCategoryService := apc_svc.NewAccountProductCategoryService(oracle.AccountProductCategory, nil, logger)
	accountProductService := ap_svc.NewAccountProductService(oracle.AccountProduct, oracle.AccountProductCategory, nil, logger, minioClient, cfg.S3BucketName, cfg)
	accountOpeningTermsService := tac_svc.NewAccountOpeningTermsService(oracle.AccountOpeningTerms, oracle.AccountProduct, nil, logger, minioClient, cfg.S3BucketName, cfg)
	walletCatch :=  walletCatch.NewWalletCatch(*sharedRedisClient,cacheLogger)

	walletService := wallet.NewWalletService(oracle.WalletOracle, nil, oracle.ServicesPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, walletCatch,logger)
	bpsActionService := bps_action_service.NewBPSActionService(persistence.BPSActionRolePersistence, persistence.BpsActionPersistence, oracle.Customer, persistence.ArchivedLinkedAccountPersistence, persistence.BPSUserPersistence, persistence.CpsUserPersistence, persistence.UserActionLogPersistence, logger, bps_action_service.Dispatcher{}, *cfg)
	topupService := topup.NewTopupService(persistence.TopupPersistence, nil, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	accountBlockService := accountblock.NewAccountService(oracle.AccountBlock, nil, logger)
	departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, nil, persistence.PortalCardPersistence, persistence.PermissionPersistence, persistence.CpsUserPersistence, logger)
	passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, nil, logger)
	hqService := hq.NewHQService(persistence.HQPersistence, nil, logger)
	miniMerchant := miniapp.NewMiniAppMerchantService(oracle.MiniAppMerchant, persistence.MiniAppPersistence, persistence.MerchantLookup, logger)
	miniAppService := miniapp.NewMiniAppService(oracle.MiniApp, miniMerchant, logger)
	ecommerceMerchantService := ecommerce_merchant.NewEcommerceMerchantService(oracle.EcommerceMerchant, oracle.ServicesPersistence, nil, persistence.MerchantLookup, coreInterface, logger, persistence.Account_lookup_external, *cfg)
	faydaService := fayda.NewFaydaService(persistence.FaydaPersistence, nil, logger)

	adService := advert.NewAdvertService(persistence.AdvertRepositoryPersistence, nil, minioClient, cfg.S3BucketName, cfg, logger)
	donationCategoryService := donation_category.NewDonationCategoryService(mongoClient, oracle.DonationCategory, oracle.Donation, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	donationCompanyService := donation_company.NewDonationCompanyService(mongoClient, oracle.DonationCompany, oracle.Donation, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl, accountLookupAdapter)
	donationService := donation.NewDonationService(mongoClient, oracle.Donation, oracle.DonationCategory, oracle.DonationCompany, oracle.ServicesPersistence, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	kycService := kycsvc.NewKYCVerifierService(mongoClient, persistence.KYCVerifierPersistence, persistence.UserPersistence, accountLookupAdapter, nil, persistence.LinkedAccountPersistence, *cfg, logger)
	permissionService := permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, nil, logger)
	budgetCategoryService := budgetCategorySvc.NewBudgetCategoryService(oracle.BudgetCategoryOracle, nil, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	cpsUserService := cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.JobRolePersistence, persistence.RolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.BPSActionApproveIndexPersistence, persistence.DepartmentPersistence, permissionService, nil, nil, persistence.RoleDelegationPersistence, minioClient, cfg.S3BucketName, *cfg, logger)
	notificationsvc := notification.InitNotificationService(persistence.NotificationPersistence, logger, nil, smsService)
	amountBased := amount_based_auth.NewAmountBasedAuthService(oracle.AmountBasedAuthOracle, nil, cfg, logger)
	unlinkService := unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, persistence.AccountBlockPersistence, nil, logger)
	bpsUserService := bpsuser.NewBPSUserService(persistence.BPSUserPersistence, persistence.JobRolePersistence, persistence.RolePersistence, nil, persistence.CpsUserPersistence, persistence.AccountBlockPersistence, minioClient, cfg.S3BucketName, *cfg, logger)
	articleService := media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	articleCategoryService := media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	ShortVideoService := media.NewShortVideoService(persistence.ShortVideoPersistence, redis, mediaProducer, logger)
	newsTagService := newstag_service.NewNewsTagService(persistence.NewsTagPersistence, nil, logger)
	newsCategoryService := newscategory_service.NewNewsCategoryService(persistence.NewsCategoryPersistence, nil, logger)
	newsTagsService := media.NewMediaTagsService(persistence.NewsTagsServiceContainer, logger)
	deviceVersionService := deviceversion.NewDeviceVersionService(persistence.DeviceVersionControlPersistence, nil, clientOrchestrationProducer, logger)
	sitotaService := sitota_service.NewSitotaTransactionService(oracle.Sitota, logger)
	encryptionService := encryption_service.NewEncryptionService(cfg, logger)
	vault := vault_category.NewVaultCategoryService(oracle.Vault, nil, logger, minioClient, minioPubUrl, VaultCategoryBucketName, cfg)
	bpsActionRoleService := bps_action_role_service.NewBPSActionRoleService(persistence.BPSActionRolePersistence, persistence.BPSActionApproveIndexPersistence, persistence.JobRolePersistence, nil, logger)
	miniAppCategory := miniapp.NewMiniAppCategoryService(oracle.MiniAppCategory, logger)
	cpsActionRoleService := cps_action_role_service.NewCPSActionRoleService(persistence.CPSActionRolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.JobRolePersistence, nil, logger)
	eventMerchantService := event_merchant_service.NewEventMerchantService(oracle.EventMerchant, nil, accountLookupAdapter, persistence.MerchantLookup, cfg, logger)
	logisticsMerchantService := logistics_merchant_service.NewLogisticsMerchantService(oracle.LogisticsMerchantOracle, nil, accountLookupAdapter, cfg, logger)
	servicesService := services_svc.NewServicesService(oracle.ServicesPersistence, persistence.UssdMerchantPersistence, nil, coreInterface, serviceCache, accessListCache, logger)
	miniAppProductCodeContainer := miniapp.NewMiniAppProductCodeService(persistence.MiniAppProductCodePersistence, logger)
	accessListSegmentationService := access_list_segmentation_service.NewAccessListSegmentationService(oracle.AccessListSegmentaion, nil, oracle.AccessListOracle, persistence.AccountBlockPersistence, persistence.CustomerService, persistence.CPSRoles, logger)
	customerSegmentationService := customer_segmentation.NewCustomerSegmentation(oracle.CustomerSegmentation, oracle.NewCPSRolesStorage, nil, nil, logger)
	jobRoleService := job_role.NewJobRoleService(persistence.JobRolePersistence, persistence.RolePersistence, nil, persistence.CpsUserPersistence, *cfg, logger)
	RoleService := roles.NewRoleService(persistence.RolePersistence, persistence.PortalCardPersistence, persistence.CPSActionApproveIndexPersistence, persistence.JobRolePersistence, persistence.BPSActionApproveIndexPersistence, nil, *cfg, logger)
	CPSRolesService := cps_role.NewCPSRoleService(oracle.NewCPSRolesStorage, nil, coreInterface, logger)
	customerKYCService := kyc_service.NewCustomerKYCService(oracle.CustomerKYC, nil, persistence.CpsUserPersistence, accountLookupAdapter, coreInterface, tokenProviderService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	customerGroupService := customer_group.NewCustomerGroupService(oracle.CustomerGroup, nil, logger)
	superAppRoleService := superapp_role.NewSuperAppRoleService(oracle.SuperAppRole, nil, coreInterface, logger)
	roleDelegationService := role_delegation_service.NewRoleDelegationService(persistence.RoleDelegationPersistence, persistence.JobRolePersistence, persistence.CpsUserPersistence, persistence.BPSUserPersistence, persistence.DepartmentPersistence, persistence.RolePersistence, oracle.AccountBlock, nil, minioClient, cfg.S3BucketName, *cfg, logger)

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
		BankContainer:              bank_service,
		AccountSubTypeContainer:    accountSubTypeService,
		WalletContainer:            walletService,
		TopupContainer:             topupService,
		PasswordRuleContainer:      passwordRule,
		AccountBlockContainer:      accountBlockService,
		DepartmentContainer:        departmentService,
		HQContainer:                hqService,
		MiniAppContainer:           miniAppService,
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
		DonationContainer:          donationService,
		DonationCategoryContainer:  donationCategoryService,
		DonationCompanyContainer:   donationCompanyService,
		ArticleContainer:           articleService,
		ArticleCategoryContainer:   articleCategoryService,
		ShortVideoServiceContainer: ShortVideoService,
		NewsTagContainer:           newsTagService,
		NewsCategoryContainer:      newsCategoryService,
		SitotaContainer:            sitotaService,
		KYCVerifierContainer:     kycService,
		DeviceVersionContainer:   deviceVersionService,
		NewsTagsServiceContainer: newsTagsService,
		EncryptionContainer:      encryptionService,
		VaultCategoryContainer:            vault,
		BPSActionRoleContainer:            bpsActionRoleService,
		MiniAppCategoryContainer:          miniAppCategory,
		CPSActionRoleContainer:            cpsActionRoleService,
		EventMerchantServiceContainer:     eventMerchantService,
		LogisticsMerchantServiceContainer: logisticsMerchantService,
		ServiceContainer:                  servicesService,
		MiniAppProductCodeContainer:       miniAppProductCodeContainer,
		AccessListSegmentationContainer:   accessListSegmentationService,
		CustomerSegmentationContainer:     customerSegmentationService,
		MiniAppMerchantContainer:          miniMerchant,
		EcommerceMerchantContainer:        ecommerceMerchantService,
		CPSRolesContainer:                 CPSRolesService,
		CustomerKYCContainer:              customerKYCService,
		UssdMerchantContainer:             ussdMerchant,
		CustomerGroupContainer:            customerGroupService,
		SuperAppRoleContainer:             superAppRoleService,
		RoleDelegationContainer:           roleDelegationService,
	}

	dispatcher := cpsaction.NewDispatcher(serviceContainer)
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSActionRolePersistence, persistence.CPSAction, persistence.UserActionLogPersistence, logger, *dispatcher, minioClient, cfg.S3BucketName, cfg.MinioPublicEndPoint, *cfg)
	cpsActionService = cpsaction.WithActionRolePolicy(cpsActionService, persistence.CPSActionRolePersistence, logger)
	ussdMerchant = ussd_merchant.NewUssdMerchantService(persistence.UssdMerchantPersistence, oracle.ServicesPersistence, minioClient, cpsActionService, cfg.S3BucketName, *cfg, logger)

	eventService = event.NewEventService(persistence.EventPersistence, cpsActionService, ecommerceMerchantService, persistence.UserPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	bulkService = bulk_service.NewBulkService(oracle.AccessListOracle, cpsActionService, oracle.AccessListSegmentaion, logger)
	bank_service = bankService.NewBankService(logger, persistence.BankPersistence, oracle.BankOracle, cpsActionService, minioClient, minioPubUrl, cfg, cfg.S3BucketName)
	accountSubTypeService = account_sub_type_svc.NewAccountSubTypeService(oracle.AccountSubType, cpsActionService, logger)
	serviceContainer.AccountSubTypeContainer = accountSubTypeService
	accountProductCategoryService = apc_svc.NewAccountProductCategoryService(oracle.AccountProductCategory, cpsActionService, logger)
	serviceContainer.AccountProductCategoryContainer = accountProductCategoryService
	accountProductService = ap_svc.NewAccountProductService(oracle.AccountProduct, oracle.AccountProductCategory, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg)
	serviceContainer.AccountProductContainer = accountProductService
	accountOpeningTermsService = tac_svc.NewAccountOpeningTermsService(oracle.AccountOpeningTerms, oracle.AccountProduct, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg)
	serviceContainer.AccountOpeningTermsContainer = accountOpeningTermsService
	walletService = wallet.NewWalletService(oracle.WalletOracle, cpsActionService, oracle.ServicesPersistence, minioClient, minioPubUrl, cfg.S3BucketName, cfg,walletCatch,logger)
	topupService = topup.NewTopupService(persistence.TopupPersistence, cpsActionService, minioClient, minioPubUrl, cfg.S3BucketName, cfg, logger)
	accountBlockService = accountblock.NewAccountService(oracle.AccountBlock, cpsActionService, logger)
	departmentService = department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionPersistence, persistence.CpsUserPersistence, logger)
	passwordRule = password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	hqService = hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	miniAppService = miniapp.NewMiniAppService(persistence.MiniAppPersistence, miniMerchant, logger)
	faydaService = fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)
	adService = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, cfg.S3BucketName, cfg, logger)
	donationCategoryService = donation_category.NewDonationCategoryService(mongoClient, oracle.DonationCategory, oracle.Donation, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	donationCompanyService = donation_company.NewDonationCompanyService(mongoClient, oracle.DonationCompany, oracle.Donation, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl, accountLookupAdapter)
	donationService = donation.NewDonationService(mongoClient, oracle.Donation, oracle.DonationCategory, oracle.DonationCompany, oracle.ServicesPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	permissionService = permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, cpsActionService, logger)
	serviceContainer.PermissionContainer = permissionService
	accountValidationService := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, cpsActionService, logger)
	serviceContainer.AccountContainer = accountValidationService
	cpsUserService = cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.JobRolePersistence, persistence.RolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.BPSActionApproveIndexPersistence, persistence.DepartmentPersistence, permissionService, cpsActionService, persistence.BPSUserPersistence, persistence.RoleDelegationPersistence, minioClient, cfg.S3BucketName, *cfg, logger)
	serviceContainer.CPSUserContainer = cpsUserService
	notificationsvc = notification.InitNotificationService(persistence.NotificationPersistence, logger, cpsActionService, smsService)
	deviceVersionService = deviceversion.NewDeviceVersionService(persistence.DeviceVersionControlPersistence, cpsActionService, clientOrchestrationProducer, logger)
	serviceContainer.DeviceVersionContainer = deviceVersionService
	budgetCategoryService = budgetCategorySvc.NewBudgetCategoryService(oracle.BudgetCategoryOracle, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	serviceContainer.UnlinkContainer = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, persistence.AccountBlockPersistence, cpsActionService, logger)
	bpsUserService = bpsuser.NewBPSUserService(persistence.BPSUserPersistence, persistence.JobRolePersistence, persistence.RolePersistence, cpsActionService, persistence.CpsUserPersistence, persistence.AccountBlockPersistence, minioClient, cfg.S3BucketName, *cfg, logger)
	serviceContainer.BPSUserContainer = bpsUserService
	serviceContainer.AdContainer = advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, cfg.S3BucketName, cfg, logger)
	serviceContainer.AvatarDomian = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	avatarService = avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, cfg.S3BucketName, minioPubUrl, *cfg)
	permissionService = permission.InitPermissionService(persistence.PermissionPersistence, persistence.DepartmentPersistence, cpsActionService, logger)
	cpsUserService = cpsusersvc.NewCPSUserService(persistence.CpsUserPersistence, persistence.JobRolePersistence, persistence.RolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.BPSActionApproveIndexPersistence, persistence.DepartmentPersistence, permissionService, cpsActionService, persistence.BPSUserPersistence, persistence.RoleDelegationPersistence, minioClient, cfg.S3BucketName, *cfg, logger)
	amountBased = amount_based_auth.NewAmountBasedAuthService(oracle.AmountBasedAuthOracle, cpsActionService, cfg, logger)
	serviceContainer.AmountBasedAuthContainer = amountBased
	bpsActionRoleService = bps_action_role_service.NewBPSActionRoleService(persistence.BPSActionRolePersistence, persistence.BPSActionApproveIndexPersistence, persistence.JobRolePersistence, cpsActionService, logger)
	serviceContainer.BPSActionRoleContainer = bpsActionRoleService
	cpsActionRoleService = cps_action_role_service.NewCPSActionRoleService(persistence.CPSActionRolePersistence, persistence.CPSActionApproveIndexPersistence, persistence.JobRolePersistence, cpsActionService, logger)
	serviceContainer.CPSActionRoleContainer = cpsActionRoleService
	dispatcher = cpsaction.NewDispatcher(serviceContainer)
	cpsActionService = cpsaction.NewCPSActionService(persistence.CPSActionRolePersistence, persistence.CPSAction, persistence.UserActionLogPersistence, logger, *dispatcher, minioClient, cfg.S3BucketName, cfg.MinioPublicEndPoint, *cfg)
	cpsActionService = cpsaction.WithActionRolePolicy(cpsActionService, persistence.CPSActionRolePersistence, logger)
	serviceContainer.CPSActionContainer = cpsActionService

	accountProductCategoryService = apc_svc.NewAccountProductCategoryService(oracle.AccountProductCategory, cpsActionService, logger)
	serviceContainer.AccountProductCategoryContainer = accountProductCategoryService
	accountProductService = ap_svc.NewAccountProductService(oracle.AccountProduct, oracle.AccountProductCategory, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg)
	serviceContainer.AccountProductContainer = accountProductService
	accountOpeningTermsService = tac_svc.NewAccountOpeningTermsService(oracle.AccountOpeningTerms, oracle.AccountProduct, cpsActionService, logger, minioClient, cfg.S3BucketName, cfg)
	serviceContainer.AccountOpeningTermsContainer = accountOpeningTermsService

	// Services catalog service (uses CPSAction for maker-checker)
	// servicesService = services_svc.NewServicesService(persistence.ServicesPersistence, cpsActionService, logger)
	servicesService = services_svc.NewServicesService(oracle.ServicesPersistence, persistence.UssdMerchantPersistence, cpsActionService, coreInterface, serviceCache, accessListCache, logger)
	// serviceContainer.ServicesContainer = servicesService
	serviceContainer.MiniAppMerchantContainer = miniMerchant
	ecommerceMerchantService = ecommerce_merchant.NewEcommerceMerchantService(oracle.EcommerceMerchant, oracle.ServicesPersistence, cpsActionService, persistence.MerchantLookup, coreInterface, logger, accountLookupAdapter, *cfg)
	unlinkService = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, oracle.AccountBlock, cpsActionService, logger)
	serviceContainer.Unlink = unlinkService
	articleService = media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	serviceContainer.ArticleContainer = articleService
	articleCategoryService = media.NewMediaCategoryService(persistence.ArticleCategoryPersistence, logger)
	serviceContainer.ArticleCategoryContainer = articleCategoryService
	ShortVideoService = media.NewShortVideoService(persistence.ShortVideoPersistence, redis, mediaProducer, logger)
	serviceContainer.ShortVideoServiceContainer = ShortVideoService
	customerService = customer.NewCustomerService(oracle.Customer, persistence.BpsActionPersistence, cpsActionService, redis, smsService, accountLookupAdapter, cfg, logger)
	kycService = kycsvc.NewKYCVerifierService(mongoClient, persistence.KYCVerifierPersistence, persistence.UserPersistence, accountLookupAdapter, cpsActionService, persistence.LinkedAccountPersistence, *cfg, logger)
	sitotaService = sitota_service.NewSitotaTransactionService(oracle.Sitota, logger)
	newsTagService = newstag_service.NewNewsTagService(persistence.NewsTagPersistence, cpsActionService, logger)
	encryptionService = encryption_service.NewEncryptionService(cfg, logger)
	newsCategoryService = newscategory_service.NewNewsCategoryService(persistence.NewsCategoryPersistence, cpsActionService, logger)
	ShortVideoService = media.NewShortVideoService(persistence.ShortVideoPersistence, redis, mediaProducer, logger)
	articleService = media.NewMediaService(persistence.ArticlePersistence, redis, logger)
	unlinkService = unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, oracle.AccountBlock, cpsActionService, logger)
	vault = vault_category.NewVaultCategoryService(oracle.Vault, cpsActionService, logger, minioClient, minioPubUrl, cfg.S3BucketName, cfg)
	eventMerchantService = event_merchant_service.NewEventMerchantService(oracle.EventMerchant, cpsActionService, accountLookupAdapter, persistence.MerchantLookup, cfg, logger)

	logisticsMerchantService = logistics_merchant_service.NewLogisticsMerchantService(oracle.LogisticsMerchantOracle, cpsActionService, accountLookupAdapter, cfg, logger)

	serviceContainer.EventMerchantServiceContainer = eventMerchantService
	serviceContainer.LogisticsMerchantServiceContainer = logisticsMerchantService
	accessListSegmentationService = access_list_segmentation_service.NewAccessListSegmentationService(oracle.AccessListSegmentaion, cpsActionService, oracle.AccessListOracle, oracle.AccountBlock, persistence.CustomerService, oracle.NewCPSRolesStorage, logger)
	customerSegmentationService = customer_segmentation.NewCustomerSegmentation(oracle.CustomerSegmentation, oracle.NewCPSRolesStorage, cpsActionService, nil, logger)
	jobRoleService = job_role.NewJobRoleService(persistence.JobRolePersistence, persistence.RolePersistence, cpsActionService, persistence.CpsUserPersistence, *cfg, logger)
	RoleService = roles.NewRoleService(persistence.RolePersistence, persistence.PortalCardPersistence, persistence.CPSActionApproveIndexPersistence, persistence.JobRolePersistence, persistence.BPSActionApproveIndexPersistence, cpsActionService, *cfg, logger)
	CPSRolesService = cps_role.NewCPSRoleService(oracle.NewCPSRolesStorage, cpsActionService, coreInterface, logger)
	customerKYCService = kyc_service.NewCustomerKYCService(oracle.CustomerKYC, cpsActionService, persistence.CpsUserPersistence, accountLookupAdapter, coreInterface, tokenProviderService, logger, minioClient, cfg.S3BucketName, cfg, minioPubUrl)
	customerGroupService = customer_group.NewCustomerGroupService(oracle.CustomerGroup, cpsActionService, logger)
	serviceContainer.CustomerGroupContainer = customerGroupService
	superAppRoleService = superapp_role.NewSuperAppRoleService(oracle.SuperAppRole, cpsActionService, coreInterface, logger)
	serviceContainer.SuperAppRoleContainer = superAppRoleService
	roleDelegationService = role_delegation_service.NewRoleDelegationService(persistence.RoleDelegationPersistence, persistence.JobRolePersistence, persistence.CpsUserPersistence, persistence.BPSUserPersistence, persistence.DepartmentPersistence, persistence.RolePersistence, oracle.AccountBlock, cpsActionService, minioClient, cfg.S3BucketName, *cfg, logger)
	serviceContainer.RoleDelegationContainer = roleDelegationService

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
		AccountSubType:    accountSubTypeService,
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
		Services:               servicesService,
		Donation:               donationService,
		DonationCategory:       donationCategoryService,
		DonationCompany:        donationCompanyService,
		NotificationService:    notificationsvc,
		ArticleService:         articleService,
		ArticleCategoryService: articleCategoryService,
		ShortVideoService:      ShortVideoService,
		NewsTagService:         newsTagService,
		NewsCategoryService:    newsCategoryService,
		Sitota:                 sitotaService,
		KYCVerifier:     kycService,
		NewsTagsService: newsTagsService,
		DeviceVersion:   deviceVersionService,
		Encryption:      encryptionService,
		VaultCategoryService:          vault,
		BPSActionRole:                 bpsActionRoleService,
		MiniAppCategory:               miniAppCategory,
		CPSActionRole:                 cpsActionRoleService,
		EventMerchantService:          eventMerchantService,
		LogisticsMerchantService:      logisticsMerchantService,
		MiniappProductCode:            miniAppProductCodeContainer,
		AccessListSegmentationService: accessListSegmentationService,
		CustomerSegmentation:          customerSegmentationService,
		CPSRoles:                      CPSRolesService,
		CustomerKYC:                   customerKYCService,
		CustomerGroup:                 customerGroupService,
		SuperAppRole:                  superAppRoleService,
		UssdMerchantService:           ussdMerchant,
		BPSActionService:              bpsActionService,
		QueueManager:                  queueManager,
		RoleDelegationService:         roleDelegationService,
		AccountProductCategory:        accountProductCategoryService,
		AccountProduct:                accountProductService,
		AccountOpeningTerms:           accountOpeningTermsService,
		TokenProvider:                 tokenProviderService,
	}
}
