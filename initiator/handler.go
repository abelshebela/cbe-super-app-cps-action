package initiator

import (
	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	cpsActionHandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler actionInbound.CPSActionAdapter
	UnlinkHandler    unlinkInbound.UnlinkAdapter
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	return Handler{
		CpsActionHandler: cpsActionHandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		UnlinkHandler:    unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
	}
}
