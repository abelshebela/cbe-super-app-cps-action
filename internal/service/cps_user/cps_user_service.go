package cpsuser

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type cpsUserService struct {
	repo   storage.CpsUserRepository
	logger utils.Logger
}

func NewCPSUserService(repo storage.CpsUserRepository, logger utils.Logger) service.CPSUserService {
	return &cpsUserService{
		repo:   repo,
		logger: logger,
	}
}

func (c *cpsUserService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	c.logger.Infof("CPS User service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
