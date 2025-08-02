package initiator

import (
	"cbe-super-app-budget/internal/service"
	"cbe-super-app-budget/internal/service/spending"
	"cbe-super-app-budget/platform/logger"
)

type ServiceLayer struct {
	spending service.SpendingService
}

func InitServiceLayer(persistnace PersistanceLayer, logger logger.Logger) ServiceLayer {
	return ServiceLayer{
		spending: spending.NewSpendingService(persistnace.spending, logger.Named("spending_service")),
	}
}
