package initiator

import (
	// Inbound section
	accountvalidationInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	"cbe-super-app-cps-action/internal/constants/interfaces/bank"
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	bulk_service_inbound "cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"
	hqInbound "cbe-super-app-cps-action/internal/constants/interfaces/hq"

	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	customerInbound "cbe-super-app-cps-action/internal/constants/interfaces/customer"
	eventInbound "cbe-super-app-cps-action/internal/constants/interfaces/event"
	FaydaInbound "cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	passwordInbound "cbe-super-app-cps-action/internal/constants/interfaces/password_rule"
	portalCardInterface "cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	walletInbound "cbe-super-app-cps-action/internal/constants/interfaces/wallet"

	// Handler section

	advertHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/ad"
	accountValidation "cbe-super-app-cps-action/internal/handlers/rest/http/account_validation"
	advertHandlerImpl "cbe-super-app-cps-action/internal/handlers/rest/http/ad"
	bankHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bank"
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	bulkServiceHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bulk_service"
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	CustomerHandler "cbe-super-app-cps-action/internal/handlers/rest/http/customer"
	"cbe-super-app-cps-action/internal/handlers/rest/http/department"
	eventhandler "cbe-super-app-cps-action/internal/handlers/rest/http/event"
	faydaHandler "cbe-super-app-cps-action/internal/handlers/rest/http/fayda"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	hqHandler "cbe-super-app-cps-action/internal/handlers/rest/http/hq"

	passwordHandler "cbe-super-app-cps-action/internal/handlers/rest/http/password_rule"
	portalcard "cbe-super-app-cps-action/internal/handlers/rest/http/portal_card"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"
	walletHandler "cbe-super-app-cps-action/internal/handlers/rest/http/wallet"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler  actionInbound.CPSActionAdapter
	UnlinkHandler     unlinkInbound.UnlinkAdapter
	EventHandler      eventInbound.EventAdapter
	WalletHandler     walletInbound.WalletAdapter
	PasswordHandler   passwordInbound.PasswordRule
	BpsHandler        bpsInbound.BPSUserHandler
	BankHandler       bank.BankHandler
	FeedbackHandler   feedbackinterface.FeedbackAdapter
	AdvertHandler     advertHandlerInterface.ADAdapter
	PortalCardHander  portalCardInterface.PortalCardAdapter
	AccountValidation accountvalidationInterface.AccountValidation
	DepartmentHandler department.DepartmentHandler

	HqHandler          hqInbound.HQAdapter
	FaydaHandler       FaydaInbound.FaydaAccount
	bulkServiceHandler bulk_service_inbound.BulkServiceHandler
	customerHandler    customerInbound.CustomerDetail
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {

	return Handler{
		UnlinkHandler:     unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:        bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		BankHandler:       bankHandler.InitBankAdapter(serviceLayer.Bank, logger),
		CpsActionHandler:  cpsactionhandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		EventHandler:      eventhandler.InitEventAdapter(serviceLayer.EventService, logger),
		FeedbackHandler:   feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
		PortalCardHander:  portalcard.InitPortalCardAdapter(serviceLayer.PortalCard, logger),
		AdvertHandler:     advertHandlerImpl.InitAdvertAdapter(serviceLayer.Advert, logger),
		WalletHandler:     walletHandler.InitWalletAdapter(serviceLayer.Wallet, logger),
		PasswordHandler:   passwordHandler.InitPasswordRuleHandler(serviceLayer.PasswordRule, logger),
		AccountValidation: accountValidation.NewHttpAccountValidation(serviceLayer.ValidationService, logger),
		DepartmentHandler: department.NewDepartmentHandler(serviceLayer.Department, logger),
		HqHandler:         hqHandler.InitHQAdapter(serviceLayer.HQService, logger),

		bulkServiceHandler: bulkServiceHandler.InitBulkServiceAdapter(serviceLayer.BulkService, logger),
		customerHandler:    CustomerHandler.InitCustomerAdapter(serviceLayer.CustomerService, logger),
		FaydaHandler:       faydaHandler.InitFaydaHandler(serviceLayer.Fayda, logger),
	}
}
