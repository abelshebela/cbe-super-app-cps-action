package branch

import (
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type branchService struct {
	repo   storage.BranchRepository
	logger utils.Logger
}

func NewBranchService(repo storage.BranchRepository, logger utils.Logger) service.BranchService {
	return &branchService{
		repo:   repo,
		logger: logger,
	}
}

func (b *branchService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	b.logger.Infof("Branch service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
