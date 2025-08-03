package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/storage"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/users"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/spending"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/redis"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ServiceLayer struct {
	UserService storage.UserRepository
}

func InitService(persistence Persistence, redisStorage *redis.RedisStorageFactory, logger utils.Logger) ServiceLayer {
	return ServiceLayer{
		UserService: users.NewUserService()
	}
}