package miniappmerchant

import (
	"context"
	"fmt"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	account_lookup_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared_util "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MiniAppMerchantApplication interface {
	CreateOne(ctx context.Context, req entities.CreateCPSAction) (*entities.CPSAction, error)
	UpdateOne(ctx context.Context, req entities.CreateCPSAction) (*entities.CPSAction, error)
	DeleteOne(ctx context.Context, id string, req entities.CreateCPSAction) (*entities.CPSAction, error)
	EnableOne(ctx context.Context, id string, req entities.CreateCPSAction) (*entities.CPSAction, error)
	DisableOne(ctx context.Context, id string, req entities.CreateCPSAction) (*entities.CPSAction, error)
	List(ctx context.Context, filter *shared_util.Filter) (*common_util.PaginatedResponse[[]*domain.MiniAppMerchant], error)
	Detail(ctx context.Context, id string) (*domain.MiniAppMerchant, error)
}

type MiniAppMerchantHandlerImpl struct {
	service              domain.MiniAppMerchantService
	cpsService           cps_service.CPSActionService
	accountLookUpService account_lookup_service.UserSearchService
	logger               utils.Logger
}

func NewMiniAppMerchantHandler(service domain.MiniAppMerchantService, cpsService cps_service.CPSActionService,
	accountLookUpService account_lookup_service.UserSearchService, logger utils.Logger) MiniAppMerchantApplication {
	return &MiniAppMerchantHandlerImpl{
		service:              service,
		cpsService:           cpsService,
		logger:               logger,
		accountLookUpService: accountLookUpService,
	}
}

func (h *MiniAppMerchantHandlerImpl) CreateOne(ctx context.Context, req entities.CreateCPSAction) (*entities.CPSAction, error) {
	return h.handleAction(ctx, func() (*entities.CPSAction, error) {
		fmt.Println("I get called") // Debug log — consider replacing with structured logging

		curData, err := h.service.GetCurrentData(&req)
		if err != nil {
			return nil, err
		}

		_, err = h.accountLookUpService.SearchUser(ctx, curData.BankAccountNumber)
		if err != nil {
			return nil, err
		}
		return h.service.CreateMiniAppMerchant(ctx, &req)
	}, "[MiniAppMerchant.Create]")
}

func (h *MiniAppMerchantHandlerImpl) UpdateOne(ctx context.Context, req entities.CreateCPSAction) (*entities.CPSAction, error) {
	return h.handleAction(ctx, func() (*entities.CPSAction, error) {
		return h.service.UpdateMiniAppMerchant(ctx, &req)
	}, "[MiniAppMerchant.Update]")
}

func (h *MiniAppMerchantHandlerImpl) DeleteOne(ctx context.Context, id string, req entities.CreateCPSAction) (*entities.CPSAction, error) {
	return h.handleAction(ctx, func() (*entities.CPSAction, error) {
		return h.service.DeleteMiniAppMerchant(ctx, id, &req)
	}, "[MiniAppMerchant.Delete]")
}

func (h *MiniAppMerchantHandlerImpl) EnableOne(ctx context.Context, id string, req entities.CreateCPSAction) (*entities.CPSAction, error) {
	return h.handleAction(ctx, func() (*entities.CPSAction, error) {
		return h.service.EnableOrDisableMerchant(ctx, id, cps_const.RequestEnableMiniAppMerchant, &req)
	}, "[MiniAppMerchant.Enable]")
}

func (h *MiniAppMerchantHandlerImpl) DisableOne(ctx context.Context, id string, req entities.CreateCPSAction) (*entities.CPSAction, error) {
	return h.handleAction(ctx, func() (*entities.CPSAction, error) {
		return h.service.EnableOrDisableMerchant(ctx, id, cps_const.RequestDisableMiniAppMerchant, &req)
	}, "[MiniAppMerchant.Disable]")
}

func (h *MiniAppMerchantHandlerImpl) List(ctx context.Context, filter *shared_util.Filter) (*common_util.PaginatedResponse[[]*domain.MiniAppMerchant], error) {
	return h.service.ListMiniAppMerchant(ctx, filter)
}

func (h *MiniAppMerchantHandlerImpl) Detail(ctx context.Context, id string) (*domain.MiniAppMerchant, error) {
	return h.service.DetailMiniAppByID(ctx, id)
}

func (h *MiniAppMerchantHandlerImpl) handleAction(ctx context.Context, buildAction func() (*entities.CPSAction, error), logPrefix string) (*entities.CPSAction, error) {

	action, err := buildAction()
	if err != nil {
		h.logger.Errorf("%s build action failed: %v", logPrefix, err)
		return nil, err
	}

	_, err = h.cpsService.CPSActionExists(ctx, entities.CheckCPSAction{
		UserCode:      action.MakerID,
		FullName:      action.MakerName,
		Department:    action.Department,
		PhoneNumber:   action.MakerPhoneNumber,
		RequestAction: string(action.RequestAction),
	})
	if err != nil {
		h.logger.Errorf("%s CPSActionExists check failed: %v", logPrefix, err)
		return nil, err
	}

	return h.cpsService.CreateCPSAction(ctx, action)
}
