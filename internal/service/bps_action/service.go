package bps_action_service

import (
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BPSActionService struct {
	BPSActionPersistence storage.BPSActionRepository
	logger               utils.Logger
}

func NewBPSActionService(bpsActionRepo storage.BPSActionRepository, logger utils.Logger, cfg *config.VaultConfig) service.BPSActionService {
	return &BPSActionService{
		BPSActionPersistence: bpsActionRepo,
		logger:               logger,
	}
}
