package initiator

import (
	"context"

	"github.com/redis/go-redis/v9"
	sharedConfig "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	sharedRedis "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config/redis"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitRedis(ctx context.Context, cfg *sharedConfig.VaultConfig, logger utils.Logger) (*sharedRedis.RedisClient, *redis.Client, error) {
	logger.Infof("Initializing Redis connection")

	redisClient, err := sharedRedis.ConnectToRedis(ctx, cfg)
	if err != nil {
		logger.Errorf("Failed to initialize Redis connection: %v", err)
		return nil, nil, err
	}
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// 	DB:   0,
	// })

	client := sharedRedis.NewRedisClient(redisClient, *cfg)

	logger.Infof("Redis connection established successfully")
	return client, redisClient, nil
}

func CloseRedis(ctx context.Context, client *sharedRedis.RedisClient, logger utils.Logger) {
	if client == nil {
		return
	}
	if err := client.Close(ctx); err != nil {
		logger.Errorf("Failed to close Redis connection: %v", err)
	} else {
		logger.Infof("Redis connection closed successfully")
	}
}
