package hq

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type hqService struct {
	repo   storage.HQRepository
	logger utils.Logger
}

func NewHQService(repo storage.HQRepository, logger utils.Logger) service.HQService {
	return &hqService{
		repo:   repo,
		logger: logger,
	}
}

func (h *hqService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	h.logger.Infof("HQ service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
