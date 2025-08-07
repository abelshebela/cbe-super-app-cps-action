package initiator

import (
	session "cbe-super-app-member-auth/grpc"
	"cbe-super-app-member-auth/internal/service"
	"cbe-super-app-member-auth/internal/service/user"
	"cbe-super-app-member-auth/internal/token"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// "cbe-super-app-member-auth/platform/logger"
)

type ServiceLayer struct {
	UserService service.UserService
}

func InitServiceLayer(persistence Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {
	return ServiceLayer{
		UserService: user.NewUserService(persistence.UserPersistence,
			persistence.SMSSenderApi,
			persistence.OTPPersistence,
			persistence.HQPersistence,
			persistence.ResetSessionPersistence,
			*token.NewTokenService(cfg),
			minioClient,
			sessionGRPCClient,
			logger,
			*cfg,
		),
	}
}
