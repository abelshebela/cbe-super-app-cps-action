package initiator

import (
	// Inbound section
	accountvalidationInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	"cbe-super-app-cps-action/internal/constants/interfaces/bank"
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	hqInbound "cbe-super-app-cps-action/internal/constants/interfaces/hq"
	miniAppInbound "cbe-super-app-cps-action/internal/constants/interfaces/mini_app"

	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	eventInbound "cbe-super-app-cps-action/internal/constants/interfaces/event"
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
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	eventhandler "cbe-super-app-cps-action/internal/handlers/rest/http/event"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	hqHandler "cbe-super-app-cps-action/internal/handlers/rest/http/hq"
	miniapphandler "cbe-super-app-cps-action/internal/handlers/rest/http/mini_app"

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
	HqHandler         hqInbound.HQAdapter
	MiniAPPHandler    miniAppInbound.MiniAppInbound
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
		HqHandler:         hqHandler.InitHQAdapter(serviceLayer.HQService, logger),
		MiniAPPHandler:    miniapphandler.InitMiniAppAdapter(serviceLayer.MiniAppService, logger),
	}
}
