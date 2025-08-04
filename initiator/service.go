package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/user"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/token"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// "github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"
)

type ServiceLayer struct {
	UserService service.UserService
}

func InitServiceLayer(persistence Persistence, logger utils.Logger, cfg *config.VaultConfig) ServiceLayer {
	return ServiceLayer{
		UserService: user.NewUserService(persistence.UserPersistence,
			persistence.SMSSenderApi,
			persistence.OTPPersistence,
			persistence.HQPersistence,
			persistence.ResetSessionPersistence,
			*token.NewTokenService(cfg),
			logger,
			*cfg,
		),
	}
}
