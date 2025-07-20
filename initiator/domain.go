package initiator

import (
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

	budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	event_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	unlink_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	wallet_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Domain struct {
	AdDomain              ad_service.AdvertService
	AvatarDomian          avatar_domain.AvatarDomainService
	BankDomain            bank_service.BankService
	WalletDomain          wallet_service.WalletService
	FaydaDomain           *fayda_service.FaydaAccountDomain
	CustomerDomain        *customer_service.CustomerDomain
	FeedbackDomain        *feedback_service.FeedbackDomain
	UnlinkDomain          *unlink_service.Service
	BudgetDomain          *budget_service.BudgetService
	AccountDomain         account_validation.Service
	CPSUserDomain         services.CPSUserService
	PasswordRuleDomain    password_service.PasswordRuleService
	PortalCardDomain      portalcard.PortaCardInterface
	DepartmentDomain      department.Service
	HQDomain              hq.Service
	miniAppDomain         miniApp_domain.MiniAppService
	AccountBlockDomain    account_block.ApplicationServices
	PermissionDomain      permission.Service
	AmountBasedAuthDomain amount_based_auth_domain.AmountBasedAuthRepository
	ServiceDomain         service.ServiceInterface
	ActionDomain          action.ServiceInterface
	EventDomain           event_domain.EventService
	BudgetCategoryDomain  budget_category.BudgetCategoryService
}

func InitDomain(minioClient config.MinioClientInterface, persistence Persitence, logger utils.Logger) Domain {
	permissionDomain := permission.InitPermissionDomain(persistence.PermissionPersistence, persistence.PermissionPersistence, persistence.PermissionPersistence, logger)
	return Domain{
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
		miniAppDomain:         miniApp_domain.NewService(persistence.miniAppPersistance, logger),
		FaydaDomain:           fayda_service.InitFaydaAccountDomain(persistence.FaydaPersistence, logger),
		EventDomain:           event_domain.NewEventService(persistence.EventPersistence),
		BudgetCategoryDomain:  *budget_category.NewBudgetCategoryService(persistence.BudgetCategoryPersistence, logger),
	}
}
