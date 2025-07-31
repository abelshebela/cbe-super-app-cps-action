package ad

import (
	"context"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

)

// ADHandlers defines the interface for advert application logic
type ADHandlers interface {
	CreateAdvert(ctx context.Context, request entity.AdvertRequest, maker cps_entities.User) error
	UpdateAdvert(ctx context.Context, id string, request entity.AdvertRequest, maker cps_entities.User) error
	DeleteAdvert(ctx context.Context, id string, maker cps_entities.User) error
	EnableDisableAdvert(ctx context.Context, id string, maker cps_entities.User, enable bool) error
	FetchAdvertByID(ctx context.Context, id string) (*entity.Advert, error)
	FetchAdverts(ctx context.Context, filterParams *util_constant.Filter) (*utils.PaginatedResponse[[]*entity.Advert], error)
}

// ADHandler implements ADHandlers
type ADHandler struct {
	service    entity.AdvertService
	cpsService cps_service.CPSActionService
	logger     shared.Logger
}

// InitADHandler initializes a new ADHandler
func InitADHandler(service entity.AdvertService, cpsService cps_service.CPSActionService, logger shared.Logger) ADHandlers {
	return &ADHandler{
		service:    service,
		cpsService: cpsService,
		logger:     logger,
	}
}

// handleCPSAction encapsulates the common CPS action logic
func (a *ADHandler) handleCPSAction(ctx context.Context, maker cps_entities.User, requestAction constant.RequestAction, curData, prevData interface{}, actionType constant.ActionType) error {
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
		a.logger.Errorf("[event.handleCPSAction] failed to create CPS action, action: %s, error: %v", requestAction, err)
		return err
	}
	return nil
}

// CreateAdvert handles advert creation
func (a *ADHandler) CreateAdvert(ctx context.Context, request entity.AdvertRequest, maker cps_entities.User) error {
	res, err := a.service.CreateAdvert(ctx, request)
	if err != nil {
		a.logger.Errorf("[event.CreateAdvert] failed to create advert, error: %v", err)
		return err
	}

	return a.handleCPSAction(ctx, maker, constant.RequestCreateAdvert, res, nil, constant.ActionCreate)
}

// UpdateAdvert handles advert updates
func (a *ADHandler) UpdateAdvert(ctx context.Context, id string, request entity.AdvertRequest, maker cps_entities.User) error {
	curAction, prevAction, err := a.service.UpdateAdvert(ctx, id, request)
	if err != nil {
		a.logger.Errorf("[event.UpdateAdvert] failed to update advert, id: %s, error: %v", id, err)
		return err
	}

	return a.handleCPSAction(ctx, maker, constant.RequestUpdateAdvert, curAction, prevAction, constant.ActionUpdate)
}

// DeleteAdvert handles advert deletion
func (a *ADHandler) DeleteAdvert(ctx context.Context, id string, maker cps_entities.User) error {
	curAction, prevAction, err := a.service.DeleteAdvert(ctx, id)
	if err != nil {
		a.logger.Errorf("[event.DeleteAdvert] failed to delete advert, id: %s, error: %v", id, err)
		return err
	}

	return a.handleCPSAction(ctx, maker, constant.RequestDeleteAdvert, curAction, prevAction, constant.ActionDelete)
}

// EnableDisableAdvert handles enabling or disabling an advert
func (a *ADHandler) EnableDisableAdvert(ctx context.Context, id string, maker cps_entities.User, enable bool) error {
	curAction, prevAction, err := a.service.EnableDisableAdvert(ctx, id, enable)
	if err != nil {
		a.logger.Errorf("[event.EnableDisableAdvert] failed to enable/disable advert, id: %s, enable: %v, error: %v", id, enable, err)
		return err
	}

	var action constant.RequestAction
	if enable {
		action = constant.RequestEnableAdvert
	} else {
		action = constant.RequestDisableAdvert
	}

	return a.handleCPSAction(ctx, maker, action, curAction, prevAction, constant.ActionUpdate)
}

// FetchAdvertByID fetches an advert by ID
func (a *ADHandler) FetchAdvertByID(ctx context.Context, id string) (*entity.Advert, error) {
	advert, err := a.service.FetchAdvertByID(ctx, id)
	if err != nil {
		a.logger.Errorf("[event.FetchAdvertByID] failed to fetch advert, id: %s, error: %v", id, err)
		return nil, err
	}
	return advert, nil
}

// FetchAdverts fetches adverts with pagination and filtering
func (a *ADHandler) FetchAdverts(ctx context.Context, filterParams *util_constant.Filter) (*utils.PaginatedResponse[[]*entity.Advert], error) {
	response, err := a.service.FetchAdverts(ctx, filterParams)
	if err != nil {
		a.logger.Errorf("[event.FetchAdverts] failed to fetch adverts, error: %v", err)
		return nil, err
	}

	

	return response, nil
}