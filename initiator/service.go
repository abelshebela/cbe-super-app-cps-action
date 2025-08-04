package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/config"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/user"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/redis"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"
)

type ServiceLayer struct {
	UserService service.UserService
}

func InitServiceLayer(persistence Persistence, redisStorage *redis.RedisStorageFactory, logger logger.Logger, cfg *config.VaultConfig) ServiceLayer {
	return ServiceLayer{
		UserService: user.NewUserService(persistence.UserPersistence,
			persistence.OTPPersistence,
			persistence.HQPersistence,
			redisStorage,
			logger,
			cfg,
		),
	}
}
