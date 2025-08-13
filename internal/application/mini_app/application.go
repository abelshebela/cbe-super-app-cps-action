package miniapp_application

import (
	"context"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	merchant_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ApplicationAbstracts interface {
	CreateMiniApp(ctx context.Context, miniApp *dto.MiniAppCreateRequest, maker entities.User) error
	UpdateMiniApp(ctx context.Context, req *dto.MiniAppCreateRequest, maker entities.User) error
	DeleteMiniApp(ctx context.Context, maker entities.User, id string) error
	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*miniApp_domain.MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*miniApp_domain.MiniApp, error)
	EnableDisableMiniAppByID(ctx context.Context, id string, enabled bool, maker entities.User) error
}

type ApplicationStore struct {
	service         domain.MiniAppService
	Logger          utils.Logger
	cpsService      cps_service.CPSActionService
	merchantService merchant_service.MiniAppMerchantService
}

func NewApplicationService(service domain.MiniAppService,
	cpsService cps_service.CPSActionService,
	merchantService merchant_service.MiniAppMerchantService,
	logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service:         service,
		Logger:          logger,
		cpsService:      cpsService,
		merchantService: merchantService,
	}
}

// handleCPSAction encapsulates the common CPS action logic
func (a *ApplicationStore) handleCPSAction(ctx context.Context, maker entities.User, requestAction cps_const.RequestAction, curData, prevData interface{}, actionType cps_const.ActionType) error {
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

	cpsAction := a.cpsService.BuildCPSAction(ctx, entities.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})

	_, err = a.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		return err
	}

	return err
}

// setMerchantDetails retrieves merchant details and sets them in the mini app request
func (a *ApplicationStore) setMerchantDetails(ctx context.Context, miniApp *dto.MiniAppCreateRequest) error {
	if miniApp.MerchantID == "" {
		return nil
	}

	_, err := a.merchantService.DetailMiniAppByID(ctx, miniApp.MerchantID)
	if err != nil {
		if err.Error() == common_util.NotFound {
			return fmt.Errorf("MERCHANT_NOT_FOUND")
		}
		return err
	}

	return nil
}

func (a *ApplicationStore) CreateMiniApp(ctx context.Context, miniApp *dto.MiniAppCreateRequest, maker entities.User) error {
	if err := a.setMerchantDetails(ctx, miniApp); err != nil {
		return err
	}

	existing, err := a.service.GetMiniAppByName(ctx, miniApp.AppName)
	res, err := a.service.CreateMiniApp(ctx, *miniApp, maker)
	if err != nil {
		a.Logger.Errorf("[mini_app.CreateMiniApp] %v", err)
		return err
	}
	if existing != nil {
		return fmt.Errorf("APP_NAME_EXIST")
	}
	err = a.handleCPSAction(ctx, maker, cps_const.RequestCreateMiniApp, res, nil, cps_const.ActionCreate)
	if err != nil {
		a.Logger.Errorf("[mini_app.CreateMiniApp] CPS action failed: %v", err)
		return err
	}

	return nil
}

func (a *ApplicationStore) UpdateMiniApp(ctx context.Context, req *dto.MiniAppCreateRequest, maker entities.User) error {
	if err := a.setMerchantDetails(ctx, req); err != nil {
		return err
	}

	curAction, prevAction, err := a.service.UpdateMiniApp(ctx, *req, maker, req.ID)
	if err != nil {
		a.Logger.Errorf("[mini_app.UpdateMiniApp] %v", err)
		return err
	}

	err = a.handleCPSAction(ctx, maker, cps_const.RequestUpdateMiniApp, curAction, prevAction, cps_const.ActionUpdate)
	if err != nil {
		a.Logger.Errorf("[mini_app.UpdateMiniApp] CPS action failed: %v", err)
		return err
	}

	return nil
}

func (a *ApplicationStore) DeleteMiniApp(ctx context.Context, maker entities.User, id string) error {
	curAction, prevAction, err := a.service.DeleteMiniApp(ctx, maker, id)
	if err != nil {
		a.Logger.Errorf("[mini_app.DeleteMiniApp] %v", err)
		return err
	}

	err = a.handleCPSAction(ctx, maker, cps_const.RequestDeleteMiniApp, curAction, prevAction, cps_const.ActionDelete)
	if err != nil {
		a.Logger.Errorf("[mini_app.DeleteMiniApp] CPS action failed: %v", err)
		return err
	}

	return nil
}

func (a *ApplicationStore) ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*miniApp_domain.MiniApp], error) {
	list, err := a.service.ListMiniApp(ctx, filterParam)
	if err != nil {
		a.Logger.Errorf("[mini_app.ListMiniApp] %v", err)
		return nil, err
	}
	return list, nil
}

func (a *ApplicationStore) DetailMiniAppByID(ctx context.Context, id string) (*miniApp_domain.MiniApp, error) {
	detail, err := a.service.DetailMiniAppByID(ctx, id)
	if err != nil {
		a.Logger.Errorf("[mini_app.DetailMiniAppByID] %v", err)
		return nil, err
	}
	return detail, nil
}

func (a *ApplicationStore) EnableDisableMiniAppByID(ctx context.Context, id string, enabled bool, maker entities.User) error {
	curAction, prevAction, err := a.service.EnableDisableMiniApp(ctx, maker, id, enabled)
	if err != nil {
		a.Logger.Errorf("[mini_app.EnableDisableMiniAppByID] %v", err)
		return err
	}

	var action cps_const.RequestAction
	if enabled {
		action = cps_const.RequestEnableMiniApp
	} else {
		action = cps_const.RequestDisableMiniApp
	}

	err = a.handleCPSAction(ctx, maker, action, curAction, prevAction, cps_const.ActionUpdate)
	if err != nil {
		a.Logger.Errorf("[mini_app.EnableDisableMiniAppByID] CPS action failed: %v", err)
		return err
	}

	return nil
}
