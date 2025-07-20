package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	ad_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/service"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	avatar_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	bank_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	budget_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/services"
	customer_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"
	department "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	fayda_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	feedback_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback"
	hq "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	password_rule_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/repository"
	password_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	permission "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"

	// budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	cps_action_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	event_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	unlink_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	wallet_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitDomain(minioClient config.MinioClientInterface, persistence Persitence, logger utils.Logger) application.Domain {
	permissionDomain := permission.InitPermissionDomain(persistence.PermissionPersistence, persistence.PermissionPersistence, persistence.PermissionPersistence, logger)

	return application.Domain{
		AdDomain:              ad_service.InitADDomian("adverts", minioClient, persistence.advertPersistence, logger),
		AvatarDomian:          avatar_domain.InitAvatarDomain(persistence.avatarPersitence, minioClient, "avatars", logger),
		CustomerDomain:        customer_service.IntiCustomerDomain(persistence.CustomerPersistence, logger),
		FeedbackDomain:        feedback_service.InitFeedbackDomain(persistence.FeedBackPersistence, logger),
		UnlinkDomain:          unlink_service.NewUnlinkService(persistence.UnlinkPersistence),
		BudgetDomain:          budget_service.InitBudgetDomain(persistence.BudgetPersistence, logger),
		AccountDomain:         account_validation.NewAccountValidationService(persistence.AccountPersistence, persistence.BulkServicesPersistence, logger),
		CPSUserDomain:         services.NewCPSUserService(persistence.CPSUserPersistence, permissionDomain, logger),
		PasswordRuleDomain:    password_service.NewPasswordRuleService(persistence.PasswordRulesPersistence.(password_rule_repo.PasswordRuleRepository)),
		BankDomain:            bank_service.InitBankDomain(persistence.BankPersistance, minioClient, "banks", logger),
		DepartmentDomain:      *department.InitDepartmentDomain(persistence.DepartmentPersistence, persistence.DepartmentPersistence, logger),
		AccountBlockDomain:    account_block.NewAccountService(persistence.AccountBlockPersistance),
		PermissionDomain:      permissionDomain,
		HQDomain:              hq.NewService(persistence.HQPersistence, persistence.BulkServicesPersistence, logger),
		AmountBasedAuthDomain: amount_based_auth_domain.NewAmountBasedAuthService(persistence.AmountBasedAuthPersistence, logger),
		ServiceDomain:         service.NewServiceStore(persistence.ServiceDetailsStore, persistence.BulkServicesPersistence, logger),
		ActionDomain:          action.NewService(persistence.BulkServicesPersistence, logger),
		PortalCardDomain: portalcard.NewPortalCardDomain(persistence.PortalCardPersistance),
		WalletDomain:          wallet_service.InitWalletDomain(persistence.WalletPersistance, minioClient, "wallets", logger),
		MiniAppDomain:         miniApp_domain.NewService(persistence.miniAppPersistance, logger),
		FaydaDomain:           fayda_service.InitFaydaAccountDomain(persistence.FaydaPersistence, logger),
		EventDomain:           event_domain.NewEventService(persistence.EventPersistence),
		// BudgetCategoryDomain:  *budget_category.NewBudgetCategoryService(persistence.BudgetCategoryPersistence, logger),
		CPSActionDomain: cps_action_service.NewCPSActionService(persistence.CPSActionsPersistance, logger),
	}
}
