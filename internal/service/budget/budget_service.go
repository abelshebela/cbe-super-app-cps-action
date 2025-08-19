package budget

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type budgetService struct {
	repo   storage.AmountBasedAuthRepository
	logger utils.Logger
}

func NewBudgetService(repo storage.AmountBasedAuthRepository, logger utils.Logger) service.BudgetService {
	return &budgetService{
		repo:   repo,
		logger: logger,
	}
}

func (b *budgetService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	b.logger.Infof("Budget service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
