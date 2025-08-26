package initiator

import (
	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/service"
	bankService "cbe-super-app-cps-action/internal/service/bank"
	bpsService "cbe-super-app-cps-action/internal/service/bps_user"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/service/event"
	mini_app_merchant "cbe-super-app-cps-action/internal/service/mini_app_merchant"
	portalcard "cbe-super-app-cps-action/internal/service/portal_card"
	"cbe-super-app-cps-action/internal/service/unlink"
	"cbe-super-app-cps-action/internal/storage/persistance"
	feedback "cbe-super-app-cps-action/internal/service/feedback"
	advert "cbe-super-app-cps-action/internal/service/ad"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceLayer struct {
	EventService service.EventService

	CPSAction service.CPSActionService
	Feedback  service.FeedbackService
	// Services  service.ServiceContainer
	Unlink     service.UnlinkService
	BpsUser    service.BPSUserService
	Bank       service.BankService

	Advert service.AdvertService

	PortalCard service.PortalCardService
}
var advertBucketName = "advert-bucket" // TODO: Add to config

func InitServiceLayer(mongoClient *mongo.Client, persistence persistance.Persistence, logger utils.Logger, sessionGRPCClient session.SessionServiceClient, cfg *config.VaultConfig, minioClient config.MinioClientInterface) ServiceLayer {

	// Create CPS action service with the dispatcher
	cpsActionService := cpsaction.NewCPSActionService(persistence.CPSAction, persistence, logger)
	feedbackService := feedback.NewFeedbackService(persistence.FeedbackPersistence, logger)
	portalCardService := portalcard.NewportalCardService(persistence.PortalCardPersistence, logger)
	merchantService := mini_app_merchant.NewMiniAppMerchantService(persistence.MiniAppMerchantPersistence, logger)
	eventService := event.NewEventService(persistence.EventPersistence, cpsActionService, merchantService, persistence.UserPersistence, minioClient, "events", cfg, logger)
	bank_service := bankService.NewBankService(logger, persistence.BankPersistence, cpsActionService, minioClient, cfg, "banks")

	return ServiceLayer{
		CPSAction:    cpsActionService,
		BpsUser:      bpsService.NewBPSUserService(persistence.BPSUserPersistence, cpsActionService, logger),
		Bank:         bank_service,
		EventService: eventService,
		Feedback:     feedbackService,
    Advert: advert.NewAdvertService(persistence.AdvertRepositoryPersistence, cpsActionService, minioClient,advertBucketName, cfg, logger),
		Unlink:       unlink.NewUnlinkService(mongoClient, persistence.UserPersistence, persistence.ArchivedUserPersistence, persistence.LinkedAccountPersistence, persistence.ArchivedLinkedAccountPersistence, cpsActionService, logger),
		PortalCard:   portalCardService,
	}
}
