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
	customer "cbe-super-app-cps-action/internal/service/customer"
	"cbe-super-app-cps-action/internal/service/department"
	"cbe-super-app-cps-action/internal/service/event"
	"cbe-super-app-cps-action/internal/service/fayda"
	feedback "cbe-super-app-cps-action/internal/service/feedback"
	"cbe-super-app-cps-action/internal/service/hq"
	mini_app_merchant "cbe-super-app-cps-action/internal/service/mini_app_merchant"
	password "cbe-super-app-cps-action/internal/service/password_rule"
	portalcard "cbe-super-app-cps-action/internal/service/portal_card"
	service_details "cbe-super-app-cps-action/internal/service/service_details"
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
	AccountBlock      service.AccountBlockService
	Department        service.DepartmentService
	PasswordRule      service.PasswordRuleService
	HQService         service.HQService
	Fayda             service.FaydaAccountService
	ServiceDetails    service.ServiceService
}

var advertBucketName = "advert-bucket" // TODO: Add to config

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {

	// Create CPS action service with the dispatcher
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)
	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, logger)
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	merchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, logger)
	accountValidation := accountvalidation.NewAccountValidationService(persistence.ValidationRulePersistence, logger)
	bulkService := bulk_service.NewBulkService(persistence.BulkService, cpsActionService, logger)
	customerSerice := customer.NewCustomerService(persistence.CustomerService, logger)
	eventService := event.NewEventService(persistence.EventPersistence, cpsActionService, merchantService, persistence.UserPersistence, minioClient, "events", cfg, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, cfg, "banks")
	walletService := wallet.NewWalletService(persistence.WalletPersistence, cpsActionService, minioClient, "wallets", cfg, logger)
	accountBlockService := accountblock.NewAccountService(persistence.AccountBlockPersistence, cpsActionService)
	departmentService := department.NewDepartmentService(persistence.DepartmentPersistence, cpsActionService, persistence.PortalCardPersistence, persistence.PermissionGroupPersistence, logger)
	passwordRule := password.NewPasswordRuleService(persistence.PasswordRulePersistent, cpsActionService, logger)
	hqService := hq.NewHQService(persistence.HQPersistence, cpsActionService, logger)
	fayda := fayda.NewFaydaService(persistence.FaydaPersistence, cpsActionService, logger)
	serviceDetails := service_details.NewServiceDetailsService(mongoClient , persistence.ServiceDetailsPersistence,persistence.HQPersistence,cpsActionService,logger)

	return ServiceLayer{
		Bank:              bank_service,
		EventService:      eventService,
		Feedback:          feedbackService,
		CPSAction:         cpsActionService,
		BpsUser:           bpsService.NewBPSUserService(persistence.BPSUserPersistence, cpsActionService, logger),
		Advert:            advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient, advertBucketName, cfg, logger),
		Unlink:            unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger),
		Budget:            budget.NewBudgetService(persistence.IconPersistence, persistence.ColorPersistence, cpsActionService, "budget", minioClient, cfg, logger),
		PortalCard:        portalCardService,
		ValidationService: accountValidation,
		Wallet:            walletService,
		AccountBlock:      accountBlockService,
		Department:        departmentService,
		BulkService:       bulkService,
		CustomerService:   customerSerice,

		PasswordRule: passwordRule,
		HQService:    hqService,
		Fayda:        fayda,
		ServiceDetails:serviceDetails,
	}
}
