package initiator

import (
	accountvalidation_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/avatar"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/amount_based_auth_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/budget_handler"
	bulkservices_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bulk_service"
	cpsmakerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/cps-maker_handler"
	customerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/customer_handler"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/fayda_account"
	feedbackhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/feedback_handler"
	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/password_rule"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/portal_card"
	service_details_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service_details"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink_device_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/wallet"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitRoutes(r chi.Router, adapter Adapter, secretKey, key, iv string, logger utils.Logger) {
	authMiddleware := middleware.InitAuthMiddleware(secretKey, key, iv, logger)

	r.Route("/api/v1/cbesuperapp/cps_action", func(sub chi.Router) {
		ad.InitADRoutes(sub, adapter.AdAdapter, authMiddleware)
		avatar.InitAvatarRoutes(sub, adapter.AvatarAdapter, authMiddleware)
		bank.InitBankRoutes(sub, adapter.BankAdapter, authMiddleware)

		wallet.InitWalletRoutes(sub, adapter.WalletAdapter, authMiddleware)
		faydaaccount.InitFaydaRoutes(sub, adapter.FaydaAdapter, authMiddleware)
		// branch_handler.RegisterBranchRoutes(sub, adapter.BranchAdapter, authMiddleware)
		customerhandler.InitCustomerRoutes(sub, adapter.CustomerAdapter, authMiddleware)
		feedbackhandler.InitFeedbackRoutes(sub, adapter.FeedbackAdapter)
		// department_handler.InitDepartmentRoutes(sub, adapter.DepartmentAdapter, authMiddleware)
		// permission_handler.InitPermissionRoutes(sub, adapter.PermissionAdapter, authMiddleware)
		unlink_device_handler.RegisterHTTPUnlinkRoutes(sub, adapter.UnlinkAdapter, authMiddleware)
		budget_handler.InitBudgetRoutes(sub, adapter.BudgetAdapter, authMiddleware)
		bulkservices_inbound.InitServiceHandlerMaker(sub, adapter.BulkServiceAdapter, authMiddleware)
		accountvalidation_inbound.InitAccountValidationHandlerMaker(sub, adapter.AccountAdapter, authMiddleware)
		cpsmakerhandler.RegisterCPSUserMakerRoutes(sub, adapter.CPSUserAdapter, authMiddleware)
		passwordrule.RegisterPasswordRuleRoutes(sub, adapter.PasswordRuleAdapter, authMiddleware)
		portalcard.InitPortalCardRoutes(sub, adapter.PortalCardAdapter, authMiddleware)
		service_details_inbound.InitServiceDetailsRoutes(sub, adapter.ServiceDetailAdapter, authMiddleware)
	})
}
