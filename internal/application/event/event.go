package event_application

import (
	"context"
	"encoding/json"
	"fmt"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entitites "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

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

func (a *ApplicationStore) CreateEvent(ctx context.Context, event dto.EventRequest, maker cps_entitites.User) error {
	_, err := a.cpsService.CPSActionExists(ctx, cps_entitites.CheckCPSAction{
		UserCode:      maker.UserCode,
		FullName:      maker.FullName,
		Department:    maker.Department,
		PhoneNumber:   maker.PhoneNumber,
		RequestAction: string(cps_const.RequestCreateEvent),
	})
	if err != nil {
		return err
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

	res, err := a.service.CreateEvent(ctx, event)
	if err != nil {
		return err
	}

	cpsAction := a.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
		User:          maker,
		CurData:       res,
		PrevData:      nil,
		RequestAction: cps_const.RequestCreateEvent,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    cps_const.ActionCreate,
	})

	_, err = a.cpsService.CreateCPSAction(ctx, cpsAction)

	return err
}

func (a *ApplicationStore) UpdateEvent(ctx context.Context, id string, event dto.EventRequest, maker cps_entitites.User) error {

	_, err := a.cpsService.CPSActionExists(ctx, cps_entitites.CheckCPSAction{
		UserCode:      maker.UserCode,
		FullName:      maker.FullName,
		Department:    maker.Department,
		PhoneNumber:   maker.PhoneNumber,
		RequestAction: string(cps_const.RequestUpdateEvent),
	})
	if err != nil {
		return err
	}

	if event.MerchantID != "" {
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
	}

	curAction, prevAction, err := a.service.UpdateEvent(ctx, id, event)
	if err != nil {
		return err
	}

	cpsAction := a.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
		User:          maker,
		CurData:       curAction,
		PrevData:      prevAction,
		RequestAction: cps_const.RequestUpdateEvent,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    cps_const.ActionUpdate,
	})

	_, err = a.cpsService.CreateCPSAction(ctx, cpsAction)

	return err
}
func (a *ApplicationStore) DeleteEvent(ctx context.Context, id string, maker cps_entitites.User) error {

	_, err := a.cpsService.CPSActionExists(ctx, cps_entitites.CheckCPSAction{
		UserCode:      maker.UserCode,
		FullName:      maker.FullName,
		Department:    maker.Department,
		PhoneNumber:   maker.PhoneNumber,
		RequestAction: string(cps_const.RequestDeleteEvent),
	})
	if err != nil {
		return err
	}
	curAction, prevAction, err := a.service.DeleteEvent(ctx, id)
	if err != nil {
		return err
	}

	cpsAction := a.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
		User:          maker,
		CurData:       curAction,
		PrevData:      prevAction,
		RequestAction: cps_const.RequestDeleteEvent,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    cps_const.ActionDelete,
	})

	_, err = a.cpsService.CreateCPSAction(ctx, cpsAction)

	return err
}

func PrettyPrintJSON(data interface{}) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON))
	return nil
}

func (a *ApplicationStore) EnableDisableEvent(ctx context.Context, id string, maker cps_entitites.User, enable bool) error {

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

	_, err = a.cpsService.CPSActionExists(ctx, cps_entitites.CheckCPSAction{
		UserCode:      maker.UserCode,
		FullName:      maker.FullName,
		Department:    maker.Department,
		PhoneNumber:   maker.PhoneNumber,
		RequestAction: string(action),
	})
	if err != nil {
		return err
	}

	cpsAction := a.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
		User:          maker,
		CurData:       curAction,
		PrevData:      prevAction,
		RequestAction: action,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    cps_const.ActionUpdate,
	})
	_, err = a.cpsService.CreateCPSAction(ctx, cpsAction)

	return err
}
func (a *ApplicationStore) FetchEventByID(ctx context.Context, id string) (*evententity.Event, error) {
	return a.service.FetchEventByID(ctx, id)
}
func (a *ApplicationStore) FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*evententity.Event], error) {
	return a.service.FetchEvent(ctx, filterParam)
}
