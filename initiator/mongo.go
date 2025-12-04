package initiator

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
)

func InitMongo(mongoUri string, logger utils.Logger) *mongo.Client {
	// Create client options with OpenTelemetry monitor
	opts := options.Client().ApplyURI(mongoUri).SetMonitor(otelmongo.NewMonitor())

	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}

	return mongoClient
}
