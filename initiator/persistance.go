package initiator

import (
	"cbe-super-app-budget/internal/storage"
	"cbe-super-app-budget/internal/storage/persistance/budget"
	"cbe-super-app-budget/internal/storage/persistance/spending"
	"cbe-super-app-budget/platform/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type PersistanceLayer struct {
	budget         storage.BudgetRepository
	spending       storage.SpendingRepository
	budgetCategory storage.BudgetCategoryRepository
}

func InitPersistanceLayer(db *mongo.Database, logger logger.Logger) PersistanceLayer {
	return PersistanceLayer{
		budget:         budget.NewBudgetRepository(db, logger.Named("budget_persistance")),
		budgetCategory: budget.NewBudgetCategoryRepository(db, logger.Named("budger_category_persistance")),
		spending:       spending.NewSpendingRepository(db, logger.Named("spending_persistance")),
	}
}
