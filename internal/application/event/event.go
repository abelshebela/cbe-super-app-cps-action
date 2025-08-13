package event_application

import (
	"context"
	"fmt"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entitites "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	evententity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	merchant_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationAbstracts interface {
	CreateEvent(ctx context.Context, event dto.EventRequest, maker cps_entitites.User) error
	UpdateEvent(ctx context.Context, id string, event dto.EventRequest, maker cps_entitites.User) error
	DeleteEvent(ctx context.Context, id string, maker cps_entitites.User) error
	EnableDisableEvent(ctx context.Context, id string, maker cps_entitites.User, enable bool) error

	FetchEventByID(ctx context.Context, id string) (*evententity.Event, error)
	FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*evententity.Event], error)
}
type ApplicationStore struct {
	service         domain.EventService
	cpsService      cps_service.CPSActionService
	merchantService merchant_service.MiniAppMerchantService
	logger          utils.Logger
}

func NewEventApplication(service domain.EventService,
	cpsService cps_service.CPSActionService,
	merchantService merchant_service.MiniAppMerchantService,
	logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service:         service,
		cpsService:      cpsService,
		logger:          logger,
		merchantService: merchantService,
	}
}

// handleCPSAction encapsulates the common CPS action logic
func (a *ApplicationStore) handleCPSAction(ctx context.Context, maker cps_entitites.User, requestAction cps_const.RequestAction, curData, prevData interface{}, actionType cps_const.ActionType) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})
	userdata := entities.CheckCPSAction{
		UserCode:    maker.UserCode,
		Department:  maker.Department,
		FullName:    maker.FullName,
		PhoneNumber: maker.PhoneNumber,
	}
	existing, err := a.cpsService.CPSActionExists(ctx, userdata)
	if err != nil {
		return err
	}
	if existing {
		return fmt.Errorf("PENDING_REQUEST_EXISTS")
	}
	_, err = a.cpsService.CreateCPSAction(ctx, cpsAction)

	return err
}

// setMerchantDetails retrieves merchant details and sets them in the event request
func (a *ApplicationStore) setMerchantDetails(ctx context.Context, event *dto.EventRequest) error {
	if event.MerchantID == "" {
		return nil
	}

	merchant, err := a.merchantService.DetailMiniAppByID(ctx, event.MerchantID)
	if err != nil {
		if err.Error() == common_util.NotFound {
			return fmt.Errorf("MERCHANT_NOT_FOUND")
		}
		return err
	}

	event.MercahntName = merchant.MerchantName
	event.MerchantEmail = merchant.PhoneNumber
	event.MerchantPhoneNumber = merchant.PhoneNumber
	event.AccountNumber = merchant.BankAccountNumber
	return nil
}

func (a *ApplicationStore) CreateEvent(ctx context.Context, event dto.EventRequest, maker cps_entitites.User) error {
	if err := a.setMerchantDetails(ctx, &event); err != nil {
		return err
	}

	res, err := a.service.CreateEvent(ctx, event)
	if err != nil {
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestCreateEvent, res, nil, cps_const.ActionCreate)
}

func (a *ApplicationStore) UpdateEvent(ctx context.Context, id string, event dto.EventRequest, maker cps_entitites.User) error {
	if err := a.setMerchantDetails(ctx, &event); err != nil {
		return err
	}

	curAction, prevAction, err := a.service.UpdateEvent(ctx, id, event)
	if err != nil {
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestUpdateEvent, curAction, prevAction, cps_const.ActionUpdate)
}

func (a *ApplicationStore) DeleteEvent(ctx context.Context, id string, maker cps_entitites.User) error {
	curAction, prevAction, err := a.service.DeleteEvent(ctx, id)
	if err != nil {
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestDeleteEvent, curAction, prevAction, cps_const.ActionDelete)
}

func (a *ApplicationStore) EnableDisableEvent(ctx context.Context, id string, maker cps_entitites.User, enable bool) error {

	userdata := entities.CheckCPSAction{
		UserCode:    maker.UserCode,
		Department:  maker.Department,
		FullName:    maker.FullName,
		PhoneNumber: maker.PhoneNumber,
	}
	existing, err := a.cpsService.CPSActionExists(ctx, userdata)
	if err != nil {
		return err
	}
	if existing {
		return fmt.Errorf("PENDING_REQUEST_EXISTS")
	}

	curAction, prevAction, err := a.service.EnableDisableEvent(ctx, id, enable)
	if err != nil {
		return err
	}

	var action cps_const.RequestAction
	if enable {
		action = cps_const.RequestEnableEvent
	} else {
		action = cps_const.RequestDisableEvent
	}

	return a.handleCPSAction(ctx, maker, action, curAction, prevAction, cps_const.ActionUpdate)
}

func (a *ApplicationStore) FetchEventByID(ctx context.Context, id string) (*evententity.Event, error) {
	return a.service.FetchEventByID(ctx, id)
}

func (a *ApplicationStore) FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*evententity.Event], error) {
	return a.service.FetchEvent(ctx, filterParam)
}
