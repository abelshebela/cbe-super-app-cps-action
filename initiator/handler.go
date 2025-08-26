package initiator

import (
	// Inbound section
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	budget "cbe-super-app-cps-action/internal/constants/interfaces/budget"
	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"

	// Handler section
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	budgetHandler "cbe-super-app-cps-action/internal/handlers/rest/http/budget"
	cpsActionHandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler actionInbound.CPSActionAdapter
	UnlinkHandler    unlinkInbound.UnlinkAdapter
	BpsHandler       bpsInbound.BPSUserHandler
	FeedbackHandler  feedbackinterface.FeedbackAdapter
	BudgetHandler    budget.BudgetPortHandler
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	return Handler{
		CpsActionHandler: cpsActionHandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		UnlinkHandler:    unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:       bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		FeedbackHandler:  feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
		BudgetHandler:    budgetHandler.InitBudgetAdapter(serviceLayer.Budget, logger),
	}
}
