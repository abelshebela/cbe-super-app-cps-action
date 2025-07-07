package initiator

import (
	accountblock_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_block"
	accountvalidation_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/ad"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/amount_based_auth_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bank"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/branch_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/budget_handler"
	bulkservices_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bulk_service"
	cpsmakerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/cps-maker_handler"
	customerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/customer_handler"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/department_handler"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/fayda_account"
	feedbackhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/feedback_handler"
	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/password_rule"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/permission_handler"
	avatar_adapter "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/avatar"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/portal_card"
	service_details_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service_details"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink_device_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/wallet"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	accountblock "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_block"
	inboundAccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_validation"
	inboundAD "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	inboundBank "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bank"
	inboundBudget "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget"
	inboundBulkServices "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bulk_services"
	inboundDepartment "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	inboundFeedback "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/feedback"
	inboundPermission "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/service_details"
	inboundUnlink "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	inboundWallet "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/wallet"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	inboundAvatar "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/avatar"
)

type Adapter struct {
	AvatarAdapter          inboundAvatar.AvatarInbound
	BankAdapter            inboundBank.BankAdapter
	AdAdapter              inboundAD.ADAdapter
	AmountBasedAuthAdapter inbound.AmountBasedAuthHandler
	WalletAdapter          inboundWallet.WalletAdapter
	FaydaAdapter           inbound.FaydaAccount
	BranchAdapter          inbound.BranchHandler
	CustomerAdapter        inbound.CustomerDetail
	FeedbackAdapter        inboundFeedback.Feedback
	DepartmentAdapter      inboundDepartment.DepartmentPortHandler
	PermissionAdapter      inboundPermission.PermissionPortHandler
	UnlinkAdapter          inboundUnlink.UnlinkPortHandler
	BudgetAdapter          inboundBudget.BudgetPortHandler
	AccountAdapter         inboundAccount.Inbound
	BulkServiceAdapter     inboundBulkServices.Inbound
	CPSUserAdapter         inbound.CPSUserMakerHandler
	PasswordRuleAdapter    inbound.PasswordRuleInbound
	PortalCardAdapter      inbound.PortalCardBound
	ServiceDetailAdapter   service_details.ServiceDetailsInbound
	AccountBlockAdapter   accountblock.AccountBlockHandler

}

func InitAdapter(application Application, logger utils.Logger) Adapter {
	return Adapter{
		AvatarAdapter: avatar_adapter.InitAvatarHTTPHandler(application.AvatarApplication, logger),
		BankAdapter:   bank.InitBankAdapter(application.BankApplication, logger),
		AdAdapter:     ad.InitADAdapter(application.AdApplication, logger),

		WalletAdapter:   wallet.InitWalletAdapter(application.WalletApplication, logger),
		FaydaAdapter:    faydaaccount.InitFaydaAdapter(application.FaydaApplication, logger),
		CustomerAdapter: customerhandler.NewCustomerHTTPHandler(application.CustomerApplication, logger),
		FeedbackAdapter: feedbackhandler.NewFeedbackHTTPHandler(application.FeedbackApplication, logger),

		UnlinkAdapter:        unlink_device_handler.NewHTTPUnlinkHandler(application.UnlinkApplication, logger),
		BudgetAdapter:        budget_handler.NewBudgetHTTPHandler(application.BudgetApplication, logger),
		AccountAdapter:       accountvalidation_inbound.NewHttpAccountValidation(application.AccountApplication, logger),
		BulkServiceAdapter:   bulkservices_inbound.NewHttpBulkService(application.BulkServicesApplication, logger),
		CPSUserAdapter:       cpsmakerhandler.InitCPSUserMakerHandler(application.CPSUserApplication, logger),
		PasswordRuleAdapter:  passwordrule.NewPasswordRuleHTTPHandler(application.PasswordRuleApplication),
		PortalCardAdapter:    portalcard.NewportalCardHandler(application.PortalCardApplication, logger),
		ServiceDetailAdapter: service_details_inbound.NewHttpServiceDetails(application.ServiceDetailApplication, logger),

		DepartmentAdapter: department_handler.NewDepartmentHTTPHandler(application.DepartmentApplication, logger),
		AccountBlockAdapter: accountblock_handler.NewAccountBlockHandler(application.AccountBlockApplication, logger),
	}
}
