package cpsaction

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CPSActionApplication interface {
	ApproveCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	RejectCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	GetCPSActionsByDepartment(ctx context.Context, department string, status string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error)
	GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entities.CPSAction, error)
}

type cpsActionApplication struct {
	service    service.CPSActionService
	services   application.Domain
	dispatcher Dispatcher
	logger     utils.Logger
}

func NewCPSActionApplication(service service.CPSActionService, services application.Domain, dispatcher Dispatcher,
	logger utils.Logger,
) CPSActionApplication {
	return &cpsActionApplication{
		service:    service,
		services:   services,
		dispatcher: dispatcher,
		logger:     logger,
	}
}

func (a *cpsActionApplication) ApproveCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error) {
	var result *entities.CPSAction

	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := a.service.RunInTransaction(txCtx, func(txCtx context.Context) error {
		a.logger.Infof("Approving CPS action: %s", action.ActionCode)

		cpsAction, err := a.service.ApproveCPSAction(txCtx, action)
		if err != nil {
			a.logger.Errorf("failed to approve CPS action: %v", err)
			return err
		}

		cpsRes, err := a.dispatcher.Authorize(txCtx, cpsAction)
		if err != nil {
			a.logger.Errorf("failed to authorize CPS action: %v", err)
			return err
		}

		a.logger.Infof("Successfully approved and authorized CPS action: %s", cpsRes.ID)
		result = cpsRes
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *cpsActionApplication) RejectCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error) {
	return a.service.RejectCPSAction(ctx, action)
}

func (a *cpsActionApplication) GetCPSActionsByDepartment(ctx context.Context, department string, status string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.CPSAction], error) {
	return a.service.GetCPSActionsByDepartment(ctx, department, status, filterParams)
}

func (a *cpsActionApplication) GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error) {
	return a.service.GetCPSActionByID(ctx, id)
}

func (a *cpsActionApplication) GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entities.CPSAction, error) {
	return a.service.GetCPSActionByActionCode(ctx, uniqueID)
}
