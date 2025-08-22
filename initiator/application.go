package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_block"

	accountvalidation_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_validation"
	ad "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	avatar_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bank"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget"
	cpsusermaker "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/customer"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/department"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/fayda_account"
	feedback "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/feedback"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/hq"
	miniApp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/password_rule"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/permission"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/portal_card"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/wallet"

	account_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_lookup"
	amount_based_auth_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/amount_based_auth_app"
	budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget_category"
	cps_actions_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_action"
	mini_app_merchant_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app_merchant"

	bulk_service_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bulk_service"
	event_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/event"
	notification_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/notification"
	service_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service"
	unlink_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/unlink"

	bps_user "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bps_user"
	donation_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/donation"
	productcode "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/product_code"
	file "github.com/CBE-Super-App/cbe-super-app-cps-action/utils/file"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Application struct {
	AvatarApplication          avatar_app.AvatarApplicationService
	BankApplication            bank.BankHandlerService
	AdApplication              ad.ADHandlers
	WalletApplication          wallet.WalletHandlerAppllication
	FaydaApplication           faydaaccount.ApplicationService
	CustomerApplication        customer.ApplicationService
	FeedbackApplication        feedback.FeedbackService
	UnlinkApplication          unlink_application.UnlinkAccount
	BudgetApplication          budget.BudgetService
	AccountApplication         accountvalidation_app.ApplicationAbstracts
	CPSUserApplication         cpsusermaker.ApplicationService
	BPSUserApplication         bps_user.ApplicationService
	PasswordRuleApplication    *passwordrule.PasswordRuleHandler
	PortalCardApplication      portalcard.PortalCardApplication
	DepartmentApplication      department.DepartmentService
	AccountBlockApplication    account_block.ApplicationService
	HQApplication              hq.ApplicationAbstracts
	PermissionApplication      permission.PermissionService
	AmountBasedAuthApplication amount_based_auth_app.ApplicationService
	MiniAppApplication         miniApp_application.ApplicationAbstracts
	EventApplication           event_application.ApplicationAbstracts
	CPSActionApplication       cps_actions_application.CPSActionApplication
	DispatcherApplication      cps_actions_application.Dispatcher

	BudgetCategoryApplication  budget_category.BudgetCategoryApplicationService
	FileService                file.FileService
	MiniAppMerchantApplication mini_app_merchant_application.MiniAppMerchantApplication
	ServicCheckeApplication    service_application.ApplicationService
	AccountLookupApplication   account_application.UserSearchService

	BulkServiceApplication  bulk_service_app.BulkServiceApplication
	NotificationApplication notification_application.NotificationApplicationAbstracts
	ProductCodeApplication  productcode.Application
	DonationApplication     donation_application.DonationAbstract
}

// InitApplication initializes the application layer with the provided domain, minio client, and logger.
func InitApplication(domain application.Domain, minioClient config.MinioClientInterface, logger utils.Logger, cfg *config.VaultConfig) Application {
	dispatcher := cps_actions_application.NewDispatcher(domain)
	return Application{
		BankApplication:     bank.InitBankHanlder(domain.BankDomain, logger),
		AdApplication:       ad.InitADHandler(domain.AdDomain, domain.CPSActionDomain, logger),
		AvatarApplication:   avatar_app.InitAvatarAPP(domain.AvatarDomian, domain.CPSActionDomain, logger),
		WalletApplication:   wallet.InitWalletApplication(domain.WalletDomain, domain.CPSActionDomain, logger),
		FaydaApplication:    faydaaccount.InitFaydaHandler(domain.FaydaDomain, logger),
		CustomerApplication: customer.InitCustomerHandler(domain.CustomerDomain, logger),
		FeedbackApplication: feedback.InitFeedbackHandler(domain.FeedbackDomain, logger),

		UnlinkApplication: unlink_application.NewUnlinkApplication(domain.UnlinkDomain, logger),

		BudgetApplication:       budget.InitBudgetHandler(domain.BudgetDomain, minioClient, "icons", logger, cfg),
		AccountApplication:      accountvalidation_app.NewApplication(domain.AccountDomain, logger),
		CPSUserApplication:      cpsusermaker.NewApplicationHandler(domain.CPSUserDomain, logger),
		BPSUserApplication:      bps_user.NewApplicationHandler(domain.BPSUserDomain, logger),
		PasswordRuleApplication: passwordrule.InitPasswordRuleHandler(domain.PasswordRuleDomain, logger),
		PortalCardApplication:   portalcard.NewPortalCardApp(domain.PortalCardDomain, logger),

		DepartmentApplication: department.InitDepartmentHandler(domain.DepartmentDomain, domain.CPSActionDomain, domain.PortalCardDomain, logger),

		AccountBlockApplication:    account_block.NewApplicationHandler(domain.AccountBlockDomain),
		HQApplication:              hq.NewApplication(domain.HQDomain, logger),
		PermissionApplication:      permission.InitPermissionHandler(&domain.PermissionDomain, logger),
		AmountBasedAuthApplication: amount_based_auth_app.ApplicationService(domain.AmountBasedAuthDomain),
		MiniAppApplication:         miniApp_application.NewApplicationService(domain.MiniAppDomain, domain.CPSActionDomain, domain.MiniAppMerchantDomain, logger),
		EventApplication:           event_application.NewEventApplication(domain.EventDomain, domain.CPSActionDomain, domain.MiniAppMerchantDomain, logger),
		BudgetCategoryApplication:  budget_category.InitBudgetCategoryHandler(domain.BudgetCategoryDomain, logger),
		DispatcherApplication:      *dispatcher,
		CPSActionApplication:       cps_actions_application.NewCPSActionApplication(domain.CPSActionDomain, domain, *dispatcher),
		ServicCheckeApplication:    service_application.NewServiceApplication(domain.ServiceCheckDomain, logger),
		MiniAppMerchantApplication: mini_app_merchant_application.NewMiniAppMerchantApplication(domain.MiniAppMerchantDomain, domain.CPSActionDomain, *domain.AccountLookup, logger),
		AccountLookupApplication:   *account_application.NewUserSearchService(domain.AccountLookup),

		BulkServiceApplication:  *bulk_service_app.NewApplicationHandler(domain.BulkServiceDomain, logger),
		NotificationApplication: notification_application.NewNotificationApplication(domain.NotificationService, domain.CPSActionDomain, logger),
		ProductCodeApplication:  productcode.NewApplication(domain.ProductCodeService, domain.CPSActionDomain, logger),
		DonationApplication:     donation_application.NewDonationApplication(domain.DonationDomain, domain.CPSActionDomain, logger),
	}
}
