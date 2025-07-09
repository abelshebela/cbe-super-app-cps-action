package initiator

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitMongo(mongo_uri string, logger utils.Logger) *mongo.Client {
	mongoClient, err := config.ConnectToMongoDB(mongo_uri)
	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}
	// defer func() {
	// 	if err := mongoClient.Disconnect(context.Background()); err != nil {
	// 		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	// 	}
	// }()

	return mongoClient
}
