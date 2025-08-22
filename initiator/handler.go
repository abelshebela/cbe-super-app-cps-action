package initiator

import (
	// Inbound section
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"

	// Handler section
	inbound "cbe-super-app-cps-action/internal/handlers/rest"
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	eventhandler "cbe-super-app-cps-action/internal/handlers/rest/http/event"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler actionInbound.CPSActionAdapter
	UnlinkHandler    unlinkInbound.UnlinkAdapter
	EventHandler     inbound.EventHandler
	BpsHandler       bpsInbound.BPSUserHandler
	FeedbackHandler  feedbackinterface.FeedbackAdapter
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	return Handler{
		UnlinkHandler:    unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:       bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		CpsActionHandler: cpsactionhandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		EventHandler:     eventhandler.InitEventAdapter(serviceLayer.EventService, logger),
		FeedbackHandler:  feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
	}
}
