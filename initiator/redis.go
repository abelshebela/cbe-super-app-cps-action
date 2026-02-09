package initiator

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitRedis(cfg *config.VaultConfig, log utils.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURI,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	log.Debugf("Initialized Redis client with URI: %s", cfg.RedisURI)

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Enable OpenTelemetry tracing for Redis
	if err := redisotel.InstrumentTracing(client); err != nil {
		log.Warnf("Failed to instrument Redis with OpenTelemetry: %v", err)
	}

	return client
}
