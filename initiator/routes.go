package initiator

import (
	accountblock_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_block"
	accountvalidation_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/avatar"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink"

	account_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_lookup"
	amount_based_auth "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/amount_based_auth_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/budget_handler"
	cpsmakerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/cps-user"
	cps_action_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/cps_action"
	mini_app_merchant_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/mini_app_merchant"

	customerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/customer_handler"
	department_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/department_handler"
	eventhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/event"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/fayda_account"
	feedbackhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/feedback_handler"
	hq_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/hq"
	miniapp_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/mini_app_handler"
	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/password_rule"
	permission_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/permission_handler"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/portal_card"

	//service_details_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service"
	bulk_service_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bulk_service"

	service_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service"

	// service_details_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service_details"
	bps_userhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bps_user"
	notificationhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/notification"
	productcodehandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/product_code"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink_device_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/wallet"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitRoutes(r chi.Router, adapter Adapter, secretKey, key, iv string, cpsService cps_service.CPSActionService, logger utils.Logger) {
	authMiddleware := middleware.InitAuthMiddleware(secretKey, key, iv, logger)
	cpsGuard := middleware.NewCPSActionMiddlewareFactory(cpsService)

	r.Route("/api/v1/cbesuperapp/cps_action", func(sub chi.Router) {
		ad.InitADRoutes(sub, adapter.AdAdapter, authMiddleware, cpsGuard)
		avatar.InitAvatarRoutes(sub, adapter.AvatarAdapter, authMiddleware, cpsGuard)
		bank.InitBankRoutes(sub, adapter.BankAdapter, authMiddleware)
		bps_userhandler.RegisterBPSUserMakerRoutes(sub, adapter.BPSUserAdapter, authMiddleware)

		wallet.InitWalletRoutes(sub, adapter.WalletAdapter, authMiddleware)
		faydaaccount.InitFaydaRoutes(sub, adapter.FaydaAdapter, authMiddleware)
		customerhandler.InitCustomerRoutes(sub, adapter.CustomerAdapter, authMiddleware)
		feedbackhandler.InitFeedbackRoutes(sub, adapter.FeedbackAdapter)
		department_handler.InitDepartmentRoutes(sub, adapter.DepartmentAdapter, authMiddleware)
		permission_handler.InitPermissionRoutes(sub, adapter.PermissionAdapter, authMiddleware)
		unlink.InitUnlinkHanldler(sub, adapter.UnlinkAdapter, authMiddleware)
		budget_handler.InitBudgetRoutes(sub, adapter.BudgetAdapter, authMiddleware)
		accountvalidation_inbound.InitAccountValidationHandlerMaker(sub, adapter.AccountAdapter, authMiddleware)
		cpsmakerhandler.RegisterCPSUserMakerRoutes(sub, adapter.CPSUserAdapter, authMiddleware)
		passwordrule.RegisterPasswordRuleRoutes(sub, adapter.PasswordRuleAdapter, authMiddleware)
		portalcard.InitPortalCardRoutes(sub, adapter.PortalCardAdapter, authMiddleware)
		accountblock_handler.RegisterAccountBlockRoutes(sub, adapter.AccountBlockAdapter, authMiddleware)
		hq_handler.InitHQRoutes(sub, adapter.HQAdapter, authMiddleware)
		amount_based_auth.InitAmountBasedAuthHandler(sub, adapter.AmountBasedAuth, authMiddleware)
		miniapp_handler.InitMiniAppHandlerMaker(sub, adapter.MiniAppAdapter, authMiddleware, cpsGuard)
		eventhandler.InitEventsHandlerMaker(sub, adapter.EventAdapter, authMiddleware, cpsGuard)
		cps_action_inbound.InitCPSActionsRoutes(sub, adapter.CPSActionAdapter, authMiddleware)
		mini_app_merchant_inbound.InitMiniAppMerchantHandlerMaker(sub, adapter.MiniAppMerchantAdapter, authMiddleware, cpsGuard)
		account_inbound.InitAccountLookUpRoutes(sub, adapter.AccountLookUp, authMiddleware)
		bulk_service_inbound.RegisterBulkServiceRoutes(sub, adapter.BulkServiceAdapter, authMiddleware)
		service_handler.InteServiceRoute(sub, adapter.ServiceCheckAdapter, authMiddleware)
		notificationhandler.InitNotificationsHandlerRoutes(sub, adapter.NotificationAdapter, authMiddleware, cpsGuard)
		productcodehandler.InitProductCodeHandlerRoutes(sub, adapter.ProductCodeAdapter, authMiddleware, cpsGuard)
	})
}
