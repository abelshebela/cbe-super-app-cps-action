package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	// "go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.uber.org/zap"
)

type Config struct {
}

func ConnectMongo(logger utils.Logger, env *config.VaultConfig) (*mongo.Client, *mongo.Database, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()

	tr := otel.Tracer("mongo-client")
	ctx, span := tr.Start(ctx, "mongo.connect")
	defer span.End()

	options := options.Client().
		ApplyURI(env.MongoDBURI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(options)
	if err != nil {
		span.SetAttributes(attribute.String("mongo.error", err.Error()))
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		span.SetAttributes(attribute.String("mongo.ping.error", err.Error()))
		return nil, nil, err
	}

	logger.Infof("Successfully Connected to MongoDB!")
	return client, client.Database(env.MongoDBDatabase), nil
}

func DisconnectMongo(ctx context.Context, client *mongo.Client, logger utils.Logger) {
	if err := client.Disconnect(ctx); err != nil {
		logger.Errorf("Error disconnecting MongoDB", zap.Error(err))
	} else {
		logger.Infof("Disconnected MongoDB successfully")
	}
}
