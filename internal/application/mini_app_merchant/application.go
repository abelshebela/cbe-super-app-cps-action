package miniapp_merchant_application

import (
	"context"
	"fmt"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	account_lookup_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MiniAppMerchantApplication interface {
	CreateMerchant(ctx context.Context, merchant dto.MiniAppMerchantRequest, maker cps_entities.User) error
	UpdateMerchant(ctx context.Context, id string, merchant dto.MiniAppMerchantRequest, maker cps_entities.User) error
	DeleteMerchant(ctx context.Context, id string, maker cps_entities.User) error
	EnableDisableMerchant(ctx context.Context, id string, maker cps_entities.User, enable bool) error
	FetchMerchantByID(ctx context.Context, id string) (*domain.MiniAppMerchant, error)
	FetchMerchants(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*domain.MiniAppMerchant], error)
}

type MiniAppMerchantApplicationStore struct {
	service              domain.MiniAppMerchantService
	cpsService           cps_service.CPSActionService
	accountLookUpService account_lookup_service.UserSearchService
	logger               utils.Logger
}

func NewMiniAppMerchantApplication(
	service domain.MiniAppMerchantService,
	cpsService cps_service.CPSActionService,
	accountLookUpService account_lookup_service.UserSearchService,
	logger utils.Logger,
) MiniAppMerchantApplication {
	return &MiniAppMerchantApplicationStore{
		service:              service,
		cpsService:           cpsService,
		accountLookUpService: accountLookUpService,
		logger:               logger,
	}
}

// handleCPSAction encapsulates the common CPS action logic
func (a *MiniAppMerchantApplicationStore) handleCPSAction(
	ctx context.Context,
	maker cps_entities.User,
	requestAction cps_const.RequestAction,
	curData, prevData interface{},
	actionType cps_const.ActionType,
) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, cps_entities.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})

	_, err := a.cpsService.CreateCPSAction(ctx, cpsAction)
	return err
}

// validateAccountDetails performs account lookup validation
func (a *MiniAppMerchantApplicationStore) validateAccountDetails(ctx context.Context, merchant *dto.MiniAppMerchantRequest) error {
	if merchant.AccountNumber == "" {
		return nil
	}

	_, err := a.accountLookUpService.SearchUser(ctx, merchant.AccountNumber)
	if err != nil {
		if err.Error() == common_util.NotFound {
			return fmt.Errorf("ACCOUNT_NOT_FOUND")
		}
		return err
	}
	return nil
}

func (a *MiniAppMerchantApplicationStore) CreateMerchant(ctx context.Context, merchant dto.MiniAppMerchantRequest, maker cps_entities.User) error {
	if err := a.validateAccountDetails(ctx, &merchant); err != nil {
		return err
	}

	res, err := a.service.CreateMiniAppMerchant(ctx, &merchant)
	if err != nil {
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestCreateMiniAppMerchant, res, nil, cps_const.ActionCreate)
}

func (a *MiniAppMerchantApplicationStore) UpdateMerchant(ctx context.Context, id string, merchant dto.MiniAppMerchantRequest, maker cps_entities.User) error {
	if err := a.validateAccountDetails(ctx, &merchant); err != nil {
		return err
	}

	curAction, prevAction, err := a.service.UpdateMiniAppMerchant(ctx, id, &merchant)
	if err != nil {
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestUpdateMiniAppMerchant, curAction, prevAction, cps_const.ActionUpdate)
}

func (a *MiniAppMerchantApplicationStore) DeleteMerchant(ctx context.Context, id string, maker cps_entities.User) error {
	curAction, prevAction, err := a.service.DeleteMiniAppMerchant(ctx, id)
	if err != nil {
		return err
	}

	return a.handleCPSAction(ctx, maker, cps_const.RequestDeleteMiniAppMerchant, curAction, prevAction, cps_const.ActionDelete)
}

func (a *MiniAppMerchantApplicationStore) EnableDisableMerchant(ctx context.Context, id string, maker cps_entities.User, enable bool) error {
	curAction, prevAction, err := a.service.EnableOrDisableMerchant(ctx, id, enable)
	if err != nil {
		return err
	}

	var action cps_const.RequestAction
	if enable {
		action = cps_const.RequestEnableMiniAppMerchant
	} else {
		action = cps_const.RequestDisableMiniAppMerchant
	}

	return a.handleCPSAction(ctx, maker, action, curAction, prevAction, cps_const.ActionUpdate)
}

func (a *MiniAppMerchantApplicationStore) FetchMerchantByID(ctx context.Context, id string) (*domain.MiniAppMerchant, error) {
	return a.service.DetailMiniAppByID(ctx, id)
}

func (a *MiniAppMerchantApplicationStore) FetchMerchants(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*domain.MiniAppMerchant], error) {
	return a.service.ListMiniAppMerchant(ctx, filterParam)
}