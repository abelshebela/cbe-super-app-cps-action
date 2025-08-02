package initiator

import (
	"cbe-super-app-budget/internal/handlers/rest"
	"cbe-super-app-budget/internal/handlers/rest/spending"
	"cbe-super-app-budget/platform/logger"
)

type HandlerLayer struct {
	spending rest.SpendingHandler
}

func InitHandlerLayer(serviceLayer ServiceLayer, logger logger.Logger) HandlerLayer {
	return HandlerLayer{
		spending: spending.NewSpendingService(serviceLayer.spending, logger.Named("spending_service")),
	}
}
