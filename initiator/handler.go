package initiator

import (
	// Inbound section
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	"cbe-super-app-cps-action/internal/constants/interfaces/productcode"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"

	// Handler section
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	cpsActionHandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	productCodeHandler "cbe-super-app-cps-action/internal/handlers/rest/http/product_code"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler   actionInbound.CPSActionAdapter
	UnlinkHandler      unlinkInbound.UnlinkAdapter
	BpsHandler         bpsInbound.BPSUserHandler
	FeedbackHandler    feedbackinterface.FeedbackAdapter
	ProductCodeHandler productcode.ProductCodeAdapter
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	pcs := serviceLayer.ProductCode
	return Handler{
		CpsActionHandler:   cpsActionHandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		UnlinkHandler:      unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:         bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		FeedbackHandler:    feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
		// it accepts productcode.Productcodeservice
		ProductCodeHandler: productCodeHandler.InitProductcodeAdapter(pcs, logger),
	}
}
