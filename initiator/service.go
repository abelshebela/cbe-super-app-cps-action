package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/service"
	bpsService "cbe-super-app-cps-action/internal/service/bps_user"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/service/unlink"
	"cbe-super-app-cps-action/internal/storage/persistance"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceLayer struct {
	CPSAction service.CPSActionService
	// Services  service.ServiceContainer
	Unlink service.UnlinkService
	BpsUser service.BPSUserService
}

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {

	// Create CPS action service with the dispatcher
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)

	return ServiceLayer{
		CPSAction: cpsActionService,
		BpsUser: bpsService.NewBPSUserService(persistence.BPSUserPersistence,cpsActionService,logger),
		// Services:  services,
		Unlink: unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger),
	}
}
