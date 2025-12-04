package initiator

import (
	redis_store "cbe-super-app-cps-action/internal/storage/redis"

	"github.com/redis/go-redis/v9"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// InitRedisStorageLayer initializes the Redis storage layer
func InitRedisStorageLayer(client *redis.Client, log utils.Logger) *redis_store.RedisStorageFactory {
	redisFactory := redis_store.NewRedisStorageFactory(client, log)
	redisFactory.Initialize()

	log.Infof("Redis storage layer initialized successfully")
	return redisFactory
}
