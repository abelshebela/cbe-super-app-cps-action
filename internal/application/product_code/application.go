package productcode

import (
	"context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constan "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// Application defines the interface for product code application logic
type Application interface {
	FetchByID(ctx context.Context, id string) (*domain.ProductCode, error)
	FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*domain.ProductCode], error)
	Update(ctx context.Context, id string, request domain.UpdateProductCodeRequest, maker cps_entities.User) error
}

// ApplicationStore implements Application
type ApplicationStore struct {
	service    domain.Service
	cpsService cps_service.CPSActionService
	logger     shared.Logger
}

// NewApplication creates a new ApplicationStore
func NewApplication(service domain.Service, cpsService cps_service.CPSActionService, logger shared.Logger) Application {
	return &ApplicationStore{
		service:    service,
		cpsService: cpsService,
		logger:     logger,
	}
}

// handleCPSAction encapsulates the common CPS action logic
func (a *ApplicationStore) handleCPSAction(ctx context.Context, maker cps_entities.User, requestAction constant.RequestAction, curData, prevData interface{}, actionType constant.ActionType) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, cps_entities.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  constant.ActionPending,
		ActionType:    actionType,
	})

	_, err := a.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("[ProductCode.handleCPSAction] failed to create CPS action, action: %s, error: %v", requestAction, err)
		return err
	}
	return nil
}

// FetchByID fetches a product code by ID
func (a *ApplicationStore) FetchByID(ctx context.Context, id string) (*domain.ProductCode, error) {
	response, err := a.service.FetchByID(ctx, id)
	if err != nil {
		a.logger.Errorf("[ProductCode.FetchByID] failed to fetch product code, id: %s, error: %v", id, err)
		return nil, err
	}
	return response, nil
}

// FetchAll fetches all product codes with pagination and filtering
func (a *ApplicationStore) FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*domain.ProductCode], error) {
	response, err := a.service.FetchAll(ctx, filterParams)
	if err != nil {
		a.logger.Errorf("[ProductCode.FetchAll] failed to fetch product codes, error: %v", err)
		return nil, err
	}
	return response, nil
}

// Update handles product code updates
func (a *ApplicationStore) Update(ctx context.Context, id string, request domain.UpdateProductCodeRequest, maker cps_entities.User) error {
	request.ID = id
	curAction, prevAction, err := a.service.Update(ctx, request)
	if err != nil {
		a.logger.Errorf("[ProductCode.Update] failed to update product code, id: %s, error: %v", id, err)
		return err
	}
	return a.handleCPSAction(ctx, maker, constant.RequestUpdateProductCode, curAction, prevAction, constant.ActionUpdate)
}
