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
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	passwordInbound "cbe-super-app-cps-action/internal/constants/interfaces/password_rule"
	portalCardInterface "cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	walletInbound "cbe-super-app-cps-action/internal/constants/interfaces/wallet"

	advertHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/ad"
	accountBlockHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"

	// Handler section
	accountValidation "cbe-super-app-cps-action/internal/handlers/rest/http/account_validation"
	advertHandlerImpl "cbe-super-app-cps-action/internal/handlers/rest/http/ad"
	bankHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bank"
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	bulkServiceHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bulk_service"
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	CustomerHandler "cbe-super-app-cps-action/internal/handlers/rest/http/customer"
	eventhandler "cbe-super-app-cps-action/internal/handlers/rest/http/event"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	hqHandler "cbe-super-app-cps-action/internal/handlers/rest/http/hq"

	passwordHandler "cbe-super-app-cps-action/internal/handlers/rest/http/password_rule"
	portalcard "cbe-super-app-cps-action/internal/handlers/rest/http/portal_card"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"
	walletHandler "cbe-super-app-cps-action/internal/handlers/rest/http/wallet"
	accountBlockHandler "cbe-super-app-cps-action/internal/handlers/rest/http/account_block"

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
	AccountBlockHandler accountBlockHandlerInterface.AccountBlockAdapter
	HqHandler         hqInbound.HQAdapter
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
		AccountBlockHandler: accountBlockHandler.InitAccountBlockAdapter(serviceLayer.AccountBlock, logger),

		bulkServiceHandler: bulkServiceHandler.InitBulkServiceAdapter(serviceLayer.BulkService, logger),
		customerHandler:    CustomerHandler.InitCustomerAdapter(serviceLayer.CustomerService, logger),
		HqHandler:         hqHandler.InitHQAdapter(serviceLayer.HQService, logger),
	}
}
