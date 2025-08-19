package initiator

import (
	inbound "cbe-super-app-cps-action/internal/handlers/rest"
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler inbound.CPSActionAdapter
}

func InitHandler(serviceLayer ServiceLayer, logger utils.Logger) Handler {
	return Handler{
		CpsActionHandler: cpsactionhandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
	}
}
