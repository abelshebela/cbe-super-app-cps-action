package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/spending"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/redis"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"
)

type ServiceLayer struct {
	spending service.SpendingService
	redis    *redis.RedisStorageFactory
}

func InitServiceLayer(persistnace PersistanceLayer, redisStorage *redis.RedisStorageFactory, logger logger.Logger) ServiceLayer {
	return ServiceLayer{
		spending: spending.NewSpendingService(persistnace.spending, logger.Named("spending_service")),
		redis:    redisStorage,
	}
}
