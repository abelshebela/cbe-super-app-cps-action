package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// ConnectRedis establishes connection to Redis
func ConnectRedis(ctx context.Context, cfg *config.VaultConfig, logger utils.Logger) (*redis.Client, error) {
	tr := otel.Tracer("redis-client")
	_, span := tr.Start(ctx, "redis.connect")
	defer span.End()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURI,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		span.SetAttributes(attribute.String("redis.error", err.Error()))
		log.Fatalf("failed to connect to redis %v", err)
	}
	return client, nil
}

// DisconnectRedis gracefully disconnects from Redis
func DisconnectRedis(ctx context.Context, client *redis.Client, logger interface{}) {
	if err := client.Close(); err != nil {
		if log, ok := logger.(interface {
			Error(context.Context, string, ...zap.Field)
		}); ok {
			log.Error(ctx, "Error disconnecting Redis", zap.Error(err))
		}
	} else {
		if log, ok := logger.(interface {
			Info(context.Context, string, ...zap.Field)
		}); ok {
			log.Info(ctx, "Disconnected Redis successfully")
		}
	}
}
