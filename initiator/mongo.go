package initiator

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitMongo(mongoUri string, logger utils.Logger) *mongo.Client {
	mongoClient, err := config.ConnectToMongoDB(mongoUri)
	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}

	return mongoClient
}
