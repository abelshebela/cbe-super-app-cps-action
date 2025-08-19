package bulk

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bulkService struct {
	repo   storage.AccountBlockRepository
	logger utils.Logger
}

func NewBulkService(repo storage.AccountBlockRepository, logger utils.Logger) service.BulkService {
	return &bulkService{
		repo:   repo,
		logger: logger,
	}
}

func (b *bulkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	b.logger.Infof("Bulk service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
