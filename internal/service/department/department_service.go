package department

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type departmentService struct {
	repo   storage.AppAccessListRepository
	logger utils.Logger
}

func NewDepartmentService(repo storage.AppAccessListRepository, logger utils.Logger) service.DepartmentService {
	return &departmentService{
		repo:   repo,
		logger: logger,
	}
}

func (d *departmentService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	d.logger.Infof("Department service authorizing action: %s", cpsAction.ActionCode)

	// For now, return the action as approved
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}
