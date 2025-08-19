package customer

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerService struct {
	repo   storage.UserRepository
	logger utils.Logger
}

func NewCustomerService(repo storage.UserRepository, logger utils.Logger) service.CustomerService {
	return &customerService{
		repo:   repo,
		logger: logger,
	}
}

func (c *customerService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	c.logger.Infof("Customer service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
