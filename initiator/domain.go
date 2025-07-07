package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	ad_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/service"
	avatar_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	bank_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	budget_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user_maker/services"
	customer_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"
	department "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	permission "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	fayda_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	feedback_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback"
	password_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	unlink_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	wallet_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Domain struct {
	AdDomain           ad_service.AdvertService
	AvatarDomian       avatar_domain.AvatarDomainService
	BankDomain         bank_service.BankService
	WalletDomain       wallet_service.WalletService
	FaydaDomain        *fayda_service.FaydaAccountDomain
	CustomerDomain     *customer_service.CustomerDomain
	FeedbackDomain     *feedback_service.FeedbackDomain
	UnlinkDomain       *unlink_service.Service
	BudgetDomain       *budget_service.BudgetService
	AccountDomain      account_validation.Service
	BulkServicesDomain action.ServiceImpl
	CPSUserDomain      services.CPSUserService
	PasswordRuleDomain password_service.PasswordRuleService
	PortalCardDomain   portalcard.PortaCardInterface
	SeviceDetailDomain service.ServiceInterface
	DepartmentDomain   department.Service

	AccountBlockDomain account_block.ApplicationServices

	PermissionDomain   permission.Service
}

func InitDomain(minioClient config.MinioClientInterface, persitence Persitence, logger utils.Logger) Domain {
	return Domain{
		AdDomain:           ad_service.InitADDomian(persitence.advertPersistence, logger),
		AvatarDomian:       avatar_domain.InitAvatarDomain(persitence.avatarPersitence, minioClient, "avatars", logger),
		CustomerDomain:     customer_service.IntiCustomerDomain(persitence.CustomerPersistence, logger),
		FeedbackDomain:     feedback_service.InitFeedbackDomain(persitence.FeedBackPersistence, logger),
		UnlinkDomain:       unlink_service.NewUnlinkService(persitence.UnlinkPersistence),
		BudgetDomain:       budget_service.InitBudgetDomain(persitence.BudgetPersistence, logger),
		AccountDomain:      account_validation.NewService(persitence.AccountPersistence, persitence.BulkServicesPersistence, logger),
		BulkServicesDomain: action.NewService(persitence.BulkServicesPersistence, logger),
		CPSUserDomain:      services.NewCPSUserService(persitence.CPSUserPersistence),
		PasswordRuleDomain: password_service.NewPasswordRuleService(persitence.PasswordRulesPersistence),
		BankDomain:         bank_service.InitBankDomain(persitence.BankPersistance, minioClient, "banks", logger),
		DepartmentDomain:   *department.InitDepartmentDomain(persitence.DepartmentPersistence, persitence.DepartmentPersistence, logger),
		AccountBlockDomain: account_block.NewAccountService(persitence.AccountBlockPersistance),
		PermissionDomain: *permission.InitPermissionDomain(persitence.PermissionPersistence, persitence.PermissionPersistence,persitence.PermissionPersistence,logger),
	}
}
