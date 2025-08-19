package miniapp

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type miniAppService struct {
	repo   storage.MiniAppRepository
	logger utils.Logger
}

func NewMiniAppService(repo storage.MiniAppRepository, logger utils.Logger) service.MiniAppService {
	return &miniAppService{
		repo:   repo,
		logger: logger,
	}
}

func (h *miniAppService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	h.logger.Infof("HQ service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
