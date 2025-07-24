package initiator

import (
	"cbe-super-app-budget/internal/handlers/rest"
	"cbe-super-app-budget/internal/handlers/rest/budget"
	"cbe-super-app-budget/internal/handlers/rest/spending"
	"cbe-super-app-budget/platform/logger"
)

type HandlerLayer struct {
	budget         rest.BudgetHandler
	budgetCategory rest.BudgetCategoryHandler
	spending       rest.SpendingHandler
}

func InitHandlerLayer(serviceLayer ServiceLayer, logger logger.Logger) HandlerLayer {
	return HandlerLayer{
		budget:         budget.NewBudgetService(serviceLayer.budget, logger.Named("budget_service")),
		budgetCategory: budget.NewBudgetCategoryService(serviceLayer.budgetCategory, logger.Named("budget_category_service")),
		spending:       spending.NewSpendingService(serviceLayer.spending, logger.Named("spending_service")),
	}
}
