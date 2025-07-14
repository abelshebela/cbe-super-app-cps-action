package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_block"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service"

	accountvalidation_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_validation"
	ad "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	avatar_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget"
	bulkservices_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bulk_services"
	cpsusermaker "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/customer"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/department"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/fayda_account"
	feedback "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/feedback"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/hq"
	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/password_rule"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/permission"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/portal_card"
	service_details_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service_details"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/unlink"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/wallet"

	amount_based_auth_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/amount_based_auth_app"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Application struct {
	AvatarApplication        avatar_app.AvatarApplicationService
	BankApplication          bank.BankHandlerService
	AdApplication            ad.ADHandlers
	WalletApplication        wallet.WalletHandlerAppllication
	FaydaApplication         faydaaccount.ApplicationService
	CustomerApplication      customer.ApplicationService
	FeedbackApplication      feedback.FeedbackService
	UnlinkApplication        unlink.ApplicationService
	BudgetApplication        budget.BudgetService
	AccountApplication       accountvalidation_app.ApplicationAbstracts
	BulkServicesApplication  bulkservices_application.ApplicationAbstracts
	CPSUserApplication       cpsusermaker.ApplicationService
	PasswordRuleApplication  *passwordrule.PasswordRuleHandler
	PortalCardApplication    portalcard.PortalCardApplication
	ServiceDetailApplication service_details_app.ApplicationAbstracts
	DepartmentApplication    department.DepartmentService
	AccountBlockApplication  account_block.ApplicationService
	HQApplication            hq.ApplicationAbstracts

	ServiceApplication         service.ServiceApplication
	PermissionApplication      permission.PermissionService
	AmountBasedAuthApplication amount_based_auth_app.ApplicationService
}

// InitApplication initializes the application layer with the provided domain, minio client, and logger.
func InitApplication(domain Domain, minioClient config.MinioClientInterface, logger utils.Logger) Application {
	return Application{
		BankApplication:            bank.InitBankHanlder(domain.BankDomain, logger),
		AdApplication:              ad.InitADHandler(domain.AdDomain, minioClient, "adverts", logger),
		AvatarApplication:          avatar_app.InitAvatarAPP(domain.AvatarDomian, logger),
		WalletApplication:          wallet.InitWalletApplication(domain.WalletDomain, logger),
		FaydaApplication:           faydaaccount.InitFaydaHandler(domain.FaydaDomain, logger),
		CustomerApplication:        customer.InitCustomerHandler(domain.CustomerDomain, logger),
		FeedbackApplication:        feedback.InitFeedbackHandler(domain.FeedbackDomain, logger),
		UnlinkApplication:          unlink.NewUnlinkHandler(domain.UnlinkDomain),
		BudgetApplication:          budget.InitBudgetHandler(domain.BudgetDomain, minioClient, "icons", logger),
		AccountApplication:         accountvalidation_app.NewApplication(domain.AccountDomain, logger),
		BulkServicesApplication:    bulkservices_application.NewAttachDetachChecker(domain.ActionDomain, logger),
		CPSUserApplication:         cpsusermaker.NewApplicationHandler(domain.CPSUserDomain, logger),
		PasswordRuleApplication:    passwordrule.InitPasswordRuleHandler(domain.PasswordRuleDomain, logger),
		PortalCardApplication:      portalcard.NewPortalCardApp(domain.PortalCardDomain, logger),
		ServiceDetailApplication:   service_details_app.NewApplication(domain.ServiceDomain, logger),
		DepartmentApplication:      department.InitDepartmentHandler(&domain.DepartmentDomain, logger),
		AccountBlockApplication:    account_block.NewApplicationHandler(domain.AccountBlockDomain),
		HQApplication:              hq.NewApplication(domain.HQDomain, logger),
		PermissionApplication:      permission.InitPermissionHandler(&domain.PermissionDomain, logger),
		AmountBasedAuthApplication: amount_based_auth_app.ApplicationService(domain.AmountBasedAuthDomain),
		ServiceApplication:         service.NewServiceApp(domain.ServiceDomain, domain.ActionDomain, logger),
	}
}
