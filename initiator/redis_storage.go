package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/redis"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"
	"github.com/redis/go-redis/v9"
)

// InitRedisStorageLayer initializes the Redis storage layer
func InitRedisStorageLayer(client *redis.Client, log logger.Logger) *redis.RedisStorageFactory {
	redisFactory := redis.NewRedisStorageFactory(client, log)
	redisFactory.Initialize()

	log.Info(nil, "Redis storage layer initialized successfully")
	return redisFactory
}
