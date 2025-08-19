package fayda

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type faydaService struct {
	repo   storage.LinkedAccountRepository
	logger utils.Logger
}

func NewFaydaService(repo storage.LinkedAccountRepository, logger utils.Logger) service.FaydaAccountService {
	return &faydaService{
		repo:   repo,
		logger: logger,
	}
}

func (f *faydaService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	f.logger.Infof("Fayda service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
