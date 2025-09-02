package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/service"
	accountblock "cbe-super-app-cps-action/internal/service/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/service/account_validation"
	advert "cbe-super-app-cps-action/internal/service/ad"
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

	Permission service.PermissionService
	CPSUser    service.CPSUserService

	ServiceDetails service.ServiceService

	ProductCode service.ProductCodeService
}

var advertBucketName = "advert-bucket" // TODO: Add to config

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {
	const minioPubUrl = "https://assetscbedev.eaglelionsystems.com"
	// Create CPS action service with the dispatcher
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)
	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, logger)
	productService := productcode.NewProductCodeService(persistence.ProductCodePersistence, cpsActionService, logger)
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	miniAppMerchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, cpsActionService, logger)
	accountValidation := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, cpsActionService, logger)
	avatar := avatar.NewAvatarService(persistence.AvatarPersistence, cpsActionService, logger, minioClient, "avatar", minioPubUrl)
	eventService := event.NewEventService(persistence.EventPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, minioClient, minioPubUrl, "events", cfg, logger)

	bulkService := bulk_service.NewBulkService(persistence.BulkService, cpsActionService, logger)
	customerSerice := customer.NewCustomerService(persistence.CustomerService, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, minioPubUrl, cfg, "banks")
	walletService := wallet.NewWalletService(persistence.WalletPersistence, cpsActionService, minioClient, "wallets", cfg, logger)
	accountBlockService := accountblock.NewAccountService(persistence.AccountBlockPersistence, cpsActionService)
	departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionPersistence, logger)
	passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	hqService := hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	keygenService := keygen.NewKeyGenerator(logger, cfg)

	miniAppService := miniapp.NewMiniAppService(persistence.MiniAppPersistence, cpsActionService, miniAppMerchantService, persistence.UserPersistence, keygenService, minioClient, minioPubUrl, "miniapps", cfg, logger)
	fayda := fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)

	serviceDetails := service_details.NewServiceDetailsService(mongoClient, persistence.ServiceDetailsPersistence, persistence.HQPersistence, cpsActionService, logger)

	permissionService := permission.InitPermissionService(
		persistence.PermissionPersistence,
		cpsActionService,
		logger,
	)

	cpsUserService := cpsusersvc.NewCPSUserService(
		persistence.CpsUserPersistence,
		persistence.AccessListPersistence, // temporary it will replaced by department repo
		permissionService,
		cpsActionService,
		logger,
	)

	return ServiceLayer{
		CPSAction:         cpsActionService,
		Feedback:          feedbackService,
		EventService:      eventService,
		Avatar:            avatar,
		Advert:            advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, minioPubUrl, advertBucketName, cfg, logger),
		BpsUser:           bpsService.NewBPSUserService(persistence.BPSUserPersistence, cpsActionService, logger),
		Bank:              bank_service,
		Unlink:            unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger),
		Budget:            budget.NewBudgetService(persistence.IconPersistence, persistence.ColorPersistence, cpsActionService, "budget", minioClient, minioPubUrl, cfg, logger),
		PortalCard:        portalCardService,
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
		Fayda:             fayda,
		Permission:        permissionService,
		CPSUser:           cpsUserService,

		ServiceDetails: serviceDetails,

		ProductCode: productService,
	}
}
