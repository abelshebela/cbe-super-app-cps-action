package initiator

import (
	accountvalidation_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget"
	bulkservices_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bulk_services"
	cpsusermaker "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_user_maker"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/customer"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/fayda_account"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/feedback"
	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/password_rule"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/portal_card"
	service_details_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service_details"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/unlink"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/wallet"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Application struct {
	BankApplication          bank.BankHandlerService
	AdApplication            ad.ADHandlers
	WalletApplication        wallet.WalletHandlerService
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
}

func InitApplication(domain Domain, minioClient config.MinioClientInterface, logger utils.Logger) Application {
	return Application{
		BankApplication:          bank.InitBankHanlder(domain.BankDomain, logger),
		AdApplication:            ad.InitADHandler(domain.AdDomain, minioClient, "adverts", logger),
		WalletApplication:        wallet.InitWalletHanlder(domain.WalletDomain, logger),
		FaydaApplication:         faydaaccount.InitFaydaHandler(domain.FaydaDomain, logger),
		CustomerApplication:      customer.InitCustomerHandler(domain.CustomerDomain, logger),
		FeedbackApplication:      feedback.InitFeedbackHandler(domain.FeedbackDomain, logger),
		UnlinkApplication:        unlink.NewUnlinkHandler(domain.UnlinkDomain),
		BudgetApplication:        budget.InitBudgetHandler(domain.BudgetDomain, minioClient, "icons", logger),
		AccountApplication:       accountvalidation_app.NewApplication(domain.AccountDomain, logger),
		BulkServicesApplication:  bulkservices_application.NewAttachDetachChecker(domain.BulkServicesDomain, logger),
		CPSUserApplication:       cpsusermaker.NewApplicationHandler(domain.CPSUserDomain),
		PasswordRuleApplication:  passwordrule.InitPasswordRuleHandler(domain.PasswordRuleDomain, logger),
		PortalCardApplication:    portalcard.NewPortalCardApp(domain.PortalCardDomain, logger),
		ServiceDetailApplication: service_details_app.NewApplication(domain.SeviceDetailDomain, logger),
	}
}
