package initiator

import (
	// Inbound section
	accountvalidationInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	"cbe-super-app-cps-action/internal/constants/interfaces/bank"
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	budget "cbe-super-app-cps-action/internal/constants/interfaces/budget"
	bulk_service_inbound "cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"
	cpsUserInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_user"
	hqInbound "cbe-super-app-cps-action/internal/constants/interfaces/hq"
	miniAppInbound "cbe-super-app-cps-action/internal/constants/interfaces/mini_app"
	permissionInbound "cbe-super-app-cps-action/internal/constants/interfaces/permission"

	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	customerInbound "cbe-super-app-cps-action/internal/constants/interfaces/customer"
	eventInbound "cbe-super-app-cps-action/internal/constants/interfaces/event"
	FaydaInbound "cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	miniAppMerchantInterface "cbe-super-app-cps-action/internal/constants/interfaces/mini_app_merchant"
	passwordInbound "cbe-super-app-cps-action/internal/constants/interfaces/password_rule"
	portalCardInterface "cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	walletInbound "cbe-super-app-cps-action/internal/constants/interfaces/wallet"

	service_details "cbe-super-app-cps-action/internal/constants/interfaces/service_details"

	accountBlockHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	advertHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/ad"
	avatarHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/avatar"

	// Handler section
	accountValidation "cbe-super-app-cps-action/internal/handlers/rest/http/account_validation"
	advertHandlerImpl "cbe-super-app-cps-action/internal/handlers/rest/http/ad"
	avatarHandlerImpl "cbe-super-app-cps-action/internal/handlers/rest/http/avatar"
	bankHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bank"
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	budgetHandler "cbe-super-app-cps-action/internal/handlers/rest/http/budget"
	bulkServiceHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bulk_service"
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	cpsUserHandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_user"
	CustomerHandler "cbe-super-app-cps-action/internal/handlers/rest/http/customer"
	"cbe-super-app-cps-action/internal/handlers/rest/http/department"
	eventhandler "cbe-super-app-cps-action/internal/handlers/rest/http/event"
	faydaHandler "cbe-super-app-cps-action/internal/handlers/rest/http/fayda"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	productCodeHandler "cbe-super-app-cps-action/internal/handlers/rest/http/product_code"

	hqHandler "cbe-super-app-cps-action/internal/handlers/rest/http/hq"
	miniapphandler "cbe-super-app-cps-action/internal/handlers/rest/http/mini_app"
	miniAppMerchantHandler "cbe-super-app-cps-action/internal/handlers/rest/http/mini_app_merchant"
	permissionHandler "cbe-super-app-cps-action/internal/handlers/rest/http/permission"

	accountBlockHandler "cbe-super-app-cps-action/internal/handlers/rest/http/account_block"
	amountBasedAuthHandler "cbe-super-app-cps-action/internal/handlers/rest/http/amount_based_auth"
	passwordHandler "cbe-super-app-cps-action/internal/handlers/rest/http/password_rule"
	portalcard "cbe-super-app-cps-action/internal/handlers/rest/http/portal_card"
	serviceDetailsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/service_details"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"
	walletHandler "cbe-super-app-cps-action/internal/handlers/rest/http/wallet"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler       actionInbound.CPSActionAdapter
	UnlinkHandler          unlinkInbound.UnlinkAdapter
	EventHandler           eventInbound.EventAdapter
	WalletHandler          walletInbound.WalletAdapter
	PasswordHandler        passwordInbound.PasswordRule
	BpsHandler             bpsInbound.BPSUserHandler
	BankHandler            bank.BankHandler
	FeedbackHandler        feedbackinterface.FeedbackAdapter
	BudgetHandler          budget.BudgetPortHandler
	AdvertHandler          advertHandlerInterface.ADAdapter
	PortalCardHander       portalCardInterface.PortalCardAdapter
	AccountValidation      accountvalidationInterface.AccountValidation
	HqHandler              hqInbound.HQAdapter
	MiniAPPHandler         miniAppInbound.MiniAppInbound
	MiniAppMerchantHandler miniAppMerchantInterface.MiniAppMerchant
	AccountBlockHandler    accountBlockHandlerInterface.AccountBlockAdapter
	ServiceDetailsHandler  service_details.ServiceAdapter
	// AmountBasedAuthHandler amountBasedAuthInbound.AmountBasedAuthAdapter

	DepartmentHandler  department.DepartmentHandler
	AvatarHandler      avatarHandlerInterface.AvatarInbound
	FaydaHandler       FaydaInbound.FaydaAccount
	bulkServiceHandler bulk_service_inbound.BulkServiceHandler
	customerHandler    customerInbound.CustomerDetail
	Permission         permissionInbound.PermissionHandler
	CPSUser            cpsUserInbound.CPSUserHandler
	ProductCodeHandler productCodeHandler.ProductCodeAdapter
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	pcs := serviceLayer.ProductCode
	return Handler{
		UnlinkHandler:       unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:          bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		BankHandler:         bankHandler.InitBankAdapter(serviceLayer.Bank, logger),
		CpsActionHandler:    cpsactionhandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		EventHandler:        eventhandler.InitEventAdapter(serviceLayer.EventService, logger),
		FeedbackHandler:     feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
		BudgetHandler:       budgetHandler.InitBudgetAdapter(serviceLayer.Budget, logger),
		PortalCardHander:    portalcard.InitPortalCardAdapter(serviceLayer.PortalCard, logger),
		AdvertHandler:       advertHandlerImpl.InitAdvertAdapter(serviceLayer.Advert, logger),
		AvatarHandler:       avatarHandlerImpl.InitAvatarAdapter(serviceLayer.Avatar, logger),
		WalletHandler:       walletHandler.InitWalletAdapter(serviceLayer.Wallet, logger),
		PasswordHandler:     passwordHandler.InitPasswordRuleHandler(serviceLayer.PasswordRule, logger),
		AccountValidation:   accountValidation.NewHttpAccountValidation(serviceLayer.ValidationService, logger),
		AccountBlockHandler: accountBlockHandler.InitAccountBlockAdapter(serviceLayer.AccountBlock, logger),

		DepartmentHandler:  department.NewDepartmentHandler(serviceLayer.Department, logger),
		HqHandler:          hqHandler.InitHQAdapter(serviceLayer.HQService, logger),
		MiniAPPHandler:     miniapphandler.InitMiniAppAdapter(serviceLayer.MiniAppService, logger),
		bulkServiceHandler: bulkServiceHandler.InitBulkServiceAdapter(serviceLayer.BulkService, logger),
		customerHandler:    CustomerHandler.InitCustomerAdapter(serviceLayer.CustomerService, logger),
		FaydaHandler:       faydaHandler.InitFaydaHandler(serviceLayer.Fayda, logger),
		Permission:         permissionHandler.InitPermissionHandler(serviceLayer.Permission, logger),
		CPSUser:            cpsUserHandler.InitCPSUserHandler(serviceLayer.CPSUser, logger),

		// AmountBasedAuthHandler: amountBasedAuthHandler.NewAmountBasedAuthHandler(serviceLayer.AmountBasedAuth, logger),

		MiniAppMerchantHandler: miniAppMerchantHandler.NewMiniAppMerchantAdapter(serviceLayer.MiniAppMerchant, logger),

		ServiceDetailsHandler: serviceDetailsHandler.InitServiceAdapter(serviceLayer.ServiceDetails, logger),

		ProductCodeHandler: productCodeHandler.InitProductcodeAdapter(pcs, logger),
	}
}
