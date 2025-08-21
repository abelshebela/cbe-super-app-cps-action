package initiator

import (
	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	cpsActionHandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler actionInbound.CPSActionAdapter
	UnlinkHandler    unlinkInbound.UnlinkAdapter
	FeedbackHandler  feedbackinterface.FeedbackAdapter
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	return Handler{
		CpsActionHandler: cpsActionHandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		UnlinkHandler:    unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		FeedbackHandler:  feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
	}
}
