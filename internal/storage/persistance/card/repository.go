package card

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository = storage.GenericRepository[model.Card]

func NewCardRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) {
	// return generic.NewMongoGenericRepository[model.Card](client, dbName, collection, logger)
}
