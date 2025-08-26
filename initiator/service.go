package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/service"
	bpsService "cbe-super-app-cps-action/internal/service/bps_user"
	"cbe-super-app-cps-action/internal/service/budget"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	feedback "cbe-super-app-cps-action/internal/service/feedback"
	"cbe-super-app-cps-action/internal/service/unlink"
	"cbe-super-app-cps-action/internal/storage/persistance"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceLayer struct {
	CPSAction service.CPSActionService
	Feedback  service.FeedbackService
	// Services  service.ServiceContainer
	Unlink  service.UnlinkService
	BpsUser service.BPSUserService
	Budget  service.BudgetService
}

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {

	// Create CPS action service with the dispatcher
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)
	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, logger)

	return ServiceLayer{
		CPSAction: cpsActionService,
		BpsUser:   bpsService.NewBPSUserService(persistence.BPSUserPersistence, cpsActionService, logger),
		Feedback:  feedbackService,
		// Services:  services,
		Unlink: unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger),
		Budget: budget.NewBudgetService(persistence.BudgetPersistence, cpsActionService, "budget", minioClient, cfg, logger),
	}
}
