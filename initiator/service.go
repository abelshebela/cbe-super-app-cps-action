package initiator

import (
	"cbe-super-app-budget/internal/service"
	"cbe-super-app-budget/internal/service/budget"
	"cbe-super-app-budget/internal/service/spending"
	"cbe-super-app-budget/platform/logger"
)

type ServiceLayer struct {
	budget         service.BudgetService
	spending       service.SpendingService
	budgetCategory service.BudgetCategoryService
}

func InitServiceLayer(persistnace PersistanceLayer, logger logger.Logger) ServiceLayer {
	return ServiceLayer{
		budget:         budget.NewBudgetService(persistnace.budget, logger.Named("budget_service")),
		budgetCategory: budget.NewBudgetCategoryService(persistnace.budgetCategory, logger.Named("budget_category_service")),
		spending:       spending.NewSpendingService(persistnace.spending, logger.Named("spending_service")),
	}
}
