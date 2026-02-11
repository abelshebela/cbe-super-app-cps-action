package initiator

import (
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
)

func InitMongo(mongoUri string, logger utils.Logger) *mongo.Client {
	// Create client options with OpenTelemetry monitor and connection pool settings

	opts := options.Client().
		ApplyURI(mongoUri).
		SetMonitor(otelmongo.NewMonitor()).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute).
		SetConnectTimeout(10 * time.Second)

	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}

	// Verify the connection with a ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Ping(ctx, nil); err != nil {
		logger.Fatalf("failed to ping mongo: %v", err)
	} //14G

	return mongoClient
}
