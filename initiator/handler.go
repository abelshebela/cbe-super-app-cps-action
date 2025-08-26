package initiator

import (
	// Inbound section
	"cbe-super-app-cps-action/internal/constants/interfaces/bank"
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	portalCardInterface "cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"

	// Handler section
	inbound "cbe-super-app-cps-action/internal/handlers/rest"
	bankHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bank"
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	eventhandler "cbe-super-app-cps-action/internal/handlers/rest/http/event"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	portalcard "cbe-super-app-cps-action/internal/handlers/rest/http/portal_card"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"
	advertHandlerImpl "cbe-super-app-cps-action/internal/handlers/rest/http/ad"
	advertHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/ad"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler actionInbound.CPSActionAdapter
	UnlinkHandler    unlinkInbound.UnlinkAdapter
	EventHandler     inbound.EventHandler
	BpsHandler       bpsInbound.BPSUserHandler
	BankHandler      bank.BankHandler
	FeedbackHandler  feedbackinterface.FeedbackAdapter
	AdvertHandler advertHandlerInterface.ADAdapter
	PortalCardHander portalCardInterface.PortalCardAdapter
	AdvertHandler advertHandlerInterface.ADAdapter
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	return Handler{
		UnlinkHandler:    unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:       bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		BankHandler:      bankHandler.InitBankAdapter(serviceLayer.Bank, logger),
		CpsActionHandler: cpsactionhandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		EventHandler:     eventhandler.InitEventAdapter(serviceLayer.EventService, logger),
		FeedbackHandler:  feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
		PortalCardHander: portalcard.InitPortalCardAdapter(serviceLayer.PortalCard, logger),
		AdvertHandler: advertHandlerImpl.InitAdvertAdapter(serviceLayer.Advert, logger),
	}
}
