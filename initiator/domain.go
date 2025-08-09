package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	ad_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	avatar_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	bank_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bps_user"
	budget_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/services"
	customer_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"
	department "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	fayda_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	feedback_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	password_rule_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/repository"
	password_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	permission "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	service_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"

	account_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
	cps_action_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	event_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	mini_app_merchant_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	notification_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	productcode "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
	unlink_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	wallet_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"
	keyGen_service "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/keygen"

	BulkServiceDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulk_service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitDomain(minioClient config.MinioClientInterface, persistence Persitence, logger utils.Logger,
	cfg *config.VaultConfig) application.Domain {
	permissionDomain := permission.InitPermissionDomain(persistence.PermissionPersistence, persistence.PermissionPersistence, persistence.PermissionPersistence, logger)
	customerDomain := customer_service.IntiCustomerDomain(persistence.CustomerPersistence, logger)
	keygenService := keyGen_service.NewKeyGenerator(logger, cfg)
	miniAppMerchantDomain := mini_app_merchant_service.NewMiniAppMerchantService(persistence.MiniAppMerchantPersisitenct, logger)

	return application.Domain{
		AdDomain:              ad_service.NewAdvertService(persistence.advertPersistence, minioClient, "adverts", cfg, logger),
		AvatarDomian:          avatar_domain.InitAvatarDomain(persistence.avatarPersitence, minioClient, "avatars", cfg, logger),
		CustomerDomain:        customerDomain,
		BPSUserDomain:         bps_user.NewBPSUserService(persistence.BPSUserPersistence, persistence.CPSActionsPersistance, logger),
		FeedbackDomain:        feedback_service.InitFeedbackDomain(persistence.FeedBackPersistence, logger),
		UnlinkDomain:          unlink_service.NewUnlinkServiceDomain(persistence.UnlinkPersistence, logger),
		BudgetDomain:          budget_service.InitBudgetDomain(persistence.BudgetPersistence, logger),
		CPSUserDomain:         services.NewCPSUserService(persistence.CPSUserPersistence, permissionDomain, persistence.DepartmentPersistence, logger),
		PasswordRuleDomain:    password_service.NewPasswordRuleService(persistence.PasswordRulesPersistence.(password_rule_repo.PasswordRuleRepository)),
		BankDomain:            bank_service.InitBankDomain(persistence.BankPersistance, minioClient, "banks", logger, cfg),
		DepartmentDomain:      *department.InitDepartmentDomain(persistence.DepartmentPersistence, persistence.DepartmentPersistence, permissionDomain, logger),
		AccountBlockDomain:    account_block.NewAccountService(persistence.AccountBlockPersistance),
		PermissionDomain:      permissionDomain,
		AmountBasedAuthDomain: amount_based_auth_domain.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, logger),
		PortalCardDomain:      portalcard.NewPortalCardDomain(persistence.PortalCardPersistance),
		WalletDomain:          wallet_service.NewWalletService(persistence.WalletPersistance, minioClient, "wallets", cfg, logger),
		MiniAppDomain:         miniApp_domain.NewService("miniapps", minioClient, persistence.miniAppPersistance, cfg, keygenService, miniAppMerchantDomain, logger),
		FaydaDomain:           fayda_service.InitFaydaAccountDomain(persistence.FaydaPersistence, logger),
		EventDomain:           event_domain.NewEventService(persistence.EventPersistence, minioClient, "events", cfg, logger),
		CPSActionDomain:       cps_action_service.NewCPSActionService(persistence.CPSActionsPersistance, logger),
		MiniAppMerchantDomain: miniAppMerchantDomain,
		AccountLookup:         account_service.NewUserSearchService(persistence.AccounLookUp),

		BulkServiceDomain:   BulkServiceDomain.NewBulkService(persistence.BulkServicePersistence, logger),
		ServiceCheckDomain:  service_domain.NewServiceDomain(persistence.ServicePersistence, logger),
		NotificationService: notification_domain.NewNotificationService(persistence.NotificationPersisitence, logger),
		ProductCodeService:  productcode.NewService(persistence.ProductCodePersistenct, logger),
		KeyGenService:       keygenService,
	}
}
