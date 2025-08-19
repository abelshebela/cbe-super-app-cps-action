package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/service"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage/persistance"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ServiceLayer struct {
	CPSAction service.CPSActionService
	// Services  service.ServiceContainer
}

func InitServiceLayer(persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {

	// Create CPS action service with the dispatcher
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)

	return ServiceLayer{
		CPSAction: cpsActionService,
		// Services:  services,
	}
}
