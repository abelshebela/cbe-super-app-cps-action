package config

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Config struct {
}

func ConnectMongo(logger logger.Logger, env *VaultConfig) (*mongo.Client, *mongo.Database, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()

	options := options.Client().
		ApplyURI(env.MongoDBURI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	logger.Info(ctx, "Successfully Connected to MongoDB!")
	return client, client.Database(env.MongoDBDatabase), nil
}

func DisconnectMongo(ctx context.Context, client *mongo.Client, logger logger.Logger) {
	if err := client.Disconnect(ctx); err != nil {
		logger.Error(ctx, "Error disconnecting MongoDB", zap.Error(err))
	} else {
		logger.Info(ctx, "Disconnected MongoDB successfully")
	}
}
