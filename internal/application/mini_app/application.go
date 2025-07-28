package miniapp_application

import (
	"context"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

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
	MakerCreateMiniApp(ctx context.Context, miniApp *dto.MiniAppCreateRequest, maker entities.User) (*entities.CPSAction, error)
	MakerUpdateMiniApp(ctx context.Context, req *dto.MiniAppCreateRequest, maker entities.User) (*entities.CPSAction, error)
	MakerDeleteMiniApp(ctx context.Context, maker entities.User, id string) (*entities.CPSAction, error)
	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*miniApp_domain.MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*miniApp_domain.MiniApp, error)
	EnableDisableMiniAppByID(ctx context.Context, id string, enabled bool, maker entities.User) (*entities.CPSAction, error)
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

func (a *ApplicationStore) MakerCreateMiniApp(ctx context.Context, miniApp *dto.MiniAppCreateRequest, maker entities.User) (*entities.CPSAction, error) {
	return a.cpsService.HandleMakerAction(ctx, func() (*entities.CPSAction, error) {
		_, err := a.merchantService.DetailMiniAppByID(ctx, miniApp.MerchantID)

		if err != nil {
			if err.Error() == common_util.NotFound {
				return nil, fmt.Errorf("MERCHANT_NOT_FOUND")
			}
			return nil, err
		}
		return a.service.CreateMiniAppAction(ctx, *miniApp, maker)
	}, "[mini_app.MakerCreateMiniApp]")
}

func (a *ApplicationStore) MakerUpdateMiniApp(ctx context.Context, req *dto.MiniAppCreateRequest, maker entities.User) (*entities.CPSAction, error) {
	return a.cpsService.HandleMakerAction(ctx, func() (*entities.CPSAction, error) {

		if req.MerchantID != "" {
			_, err := a.merchantService.DetailMiniAppByID(ctx, req.MerchantID)

			if err != nil {
				if err.Error() == common_util.NotFound {
					return nil, fmt.Errorf("MERCHANT_NOT_FOUND")
				}
				return nil, err
			}
		}
		return a.service.UpdateMiniAppAction(ctx, *req, maker, req.ID)
	}, "[mini_app.MakerUpdateMiniApp]")
}

func (a *ApplicationStore) MakerDeleteMiniApp(ctx context.Context, maker entities.User, id string) (*entities.CPSAction, error) {
	return a.cpsService.HandleMakerAction(ctx, func() (*entities.CPSAction, error) {
		return a.service.DeleteMiniAppAction(ctx, maker, id)
	}, "[mini_app.MakerDeleteMiniApp]")
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

func (a *ApplicationStore) EnableDisableMiniAppByID(ctx context.Context, id string, enabled bool, maker entities.User) (*entities.CPSAction, error) {

	return a.cpsService.HandleMakerAction(ctx, func() (*entities.CPSAction, error) {

		detail, err := a.service.EnableDisableMiniApp(ctx, maker, id, enabled)
		if err != nil {
			a.Logger.Errorf("[EnableDisableMiniAppByID] %v", err)
			return nil, err
		}
		return detail, nil

	}, "[mini_app.EnableDisableMiniAppByID]")
}

