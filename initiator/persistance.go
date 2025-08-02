package initiator

import (
	"cbe-super-app-budget/internal/storage"
	"cbe-super-app-budget/internal/storage/persistance/spending"
	"cbe-super-app-budget/platform/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type PersistanceLayer struct {
	spending storage.SpendingRepository
}

func InitPersistanceLayer(db *mongo.Database, logger logger.Logger) PersistanceLayer {
	return PersistanceLayer{
		spending: spending.NewSpendingRepository(db, logger.Named("spending_persistance")),
	}
}
