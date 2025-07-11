package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	ad_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/service"
	bank_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	budget_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/services"
	customer_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"
	department "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	fayda_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	feedback_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback"
	password_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	unlink_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	wallet_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	// ad_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/service"
	avatar_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
)

type Domain struct {
	AdDomain             ad_service.AdvertService
	AvatarDomian         avatar_domain.AvatarDomainService
	BankDomain           bank_service.BankService
	WalletDomain         wallet_service.WalletService
	FaydaDomain          *fayda_service.FaydaAccountDomain
	CustomerDomain       *customer_service.CustomerDomain
	FeedbackDomain       *feedback_service.FeedbackDomain
	UnlinkDomain         *unlink_service.Service
	BudgetDomain         *budget_service.BudgetService
	AccountDomain        account_validation.Service
	BulkServicesDomain   action.ServiceImpl
	CPSUserDomain        services.CPSUserService
	PasswordRuleDomain   password_service.PasswordRuleService
	PortalCardDomain     portalcard.PortaCardInterface
	SeviceDetailDomain   service.ServiceInterface
	DepartmentDomain     department.Service
	BudgetCategoryDomain budget_category.BudgetCategoryService
}

func InitDomain(minioClient config.MinioClientInterface, persistence Persitence, logger utils.Logger) Domain {
	return Domain{
		AdDomain:             ad_service.InitADDomian(persistence.advertPersistence, logger),
		AvatarDomian:         avatar_domain.InitAvatarDomain(persistence.avatarPersitence, minioClient, "avatars", logger),
		CustomerDomain:       customer_service.IntiCustomerDomain(persistence.CustomerPersistence, logger),
		FeedbackDomain:       feedback_service.InitFeedbackDomain(persistence.FeedBackPersistence, logger),
		UnlinkDomain:         unlink_service.NewUnlinkService(persistence.UnlinkPersistence),
		BudgetDomain:         budget_service.InitBudgetDomain(persistence.BudgetPersistence, logger),
		AccountDomain:        account_validation.NewService(persistence.AccountPersistence, persistence.BulkServicesPersistence, logger),
		BulkServicesDomain:   action.NewService(persistence.BulkServicesPersistence, logger),
		CPSUserDomain:        services.NewCPSUserService(persistence.CPSUserPersistence),
		PasswordRuleDomain:   password_service.NewPasswordRuleService(persistence.PasswordRulesPersistence),
		BankDomain:           bank_service.InitBankDomain(persistence.BankPersistance, minioClient, "banks", logger),
		DepartmentDomain:     *department.InitDepartmentDomain(persistence.DepartmentPersistence, persistence.DepartmentPersistence, logger),
		BudgetCategoryDomain: *budget_category.NewBudgetCategoryService(persistence.BudgetCategoryPersistence, logger),
	}
}
