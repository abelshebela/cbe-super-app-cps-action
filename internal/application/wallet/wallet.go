package wallet

import (
	"context"

	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"

	cps_entitites "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	common_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type WalletHandlerAppllication interface {
	GetAllWallet(ctx context.Context, filterParams *constant.MongoFilter) (*common_utils.PaginatedResponse[[]*domain.Wallet], error)
	GetWallet(ctx context.Context, id string) (*domain.Wallet, error)

	CreateWallet(ctx context.Context, req domain.WalletRequest, maker cps_entitites.User) error
	UpdateWallet(ctx context.Context, id string, req domain.WalletRequest, maker cps_entitites.User) error
	DeleteWallet(ctx context.Context, id string, maker cps_entitites.User) error
	EnableOrDisableWallet(ctx context.Context, id string, enable bool, maker cps_entitites.User) error
}

type WalletHandler struct {
	walletDomain domain.WalletService
	logger       utils.Logger
	cpsService   cps_service.CPSActionService
}

func InitWalletApplication(walletDomain domain.WalletService,
	cpsService cps_service.CPSActionService,
	logger utils.Logger) WalletHandlerAppllication {
	return &WalletHandler{
		walletDomain: walletDomain,
		logger:       logger,
		cpsService:   cpsService,
	}
}

func (w *WalletHandler) CreateWallet(ctx context.Context, req domain.WalletRequest, maker cps_entitites.User) error {
	wallet, err := w.walletDomain.CreateWallet(ctx, req)
	if err != nil {
		return err
	}

	return w.handleCPSAction(ctx, maker, cps_const.RequestCreateWallet, wallet, nil, cps_const.ActionCreate)
}

func (w *WalletHandler) DeleteWallet(ctx context.Context, id string, maker cps_entitites.User) error {
	curAction, prevAction, err := w.walletDomain.DeleteWallet(ctx, id)
	if err != nil {
		return err
	}

	return w.handleCPSAction(ctx, maker, cps_const.RequestDeleteWallet, curAction, prevAction, cps_const.ActionDelete)

}

func (w *WalletHandler) GetAllWallet(ctx context.Context, filterParams *constant.MongoFilter) (*common_utils.PaginatedResponse[[]*domain.Wallet], error) {
	banks, err := w.walletDomain.FetchWallet(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return banks, nil

}

func (w *WalletHandler) GetWallet(ctx context.Context, id string) (*domain.Wallet, error) {
	bank, err := w.walletDomain.FetchWalletByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (w *WalletHandler) UpdateWallet(ctx context.Context, id string, req domain.WalletRequest, maker cps_entitites.User) error {
	curAction, prevAction, err := w.walletDomain.UpdateWallet(ctx, id, req)
	if err != nil {
		return err
	}

	return w.handleCPSAction(ctx, maker, cps_const.RequestUpdateWallet, curAction, prevAction, cps_const.ActionUpdate)
}

func (w *WalletHandler) EnableOrDisableWallet(ctx context.Context, id string,
	enable bool, maker cps_entitites.User) error {
	curAction, prevAction, err := w.walletDomain.EnableDisableWallet(ctx, id, enable)
	if err != nil {
		return err
	}

	var action cps_const.RequestAction
	if enable {
		action = cps_const.RequestEnableWallet
	} else {
		action = cps_const.RequestDisableWallet
	}

	return w.handleCPSAction(ctx, maker, action, curAction, prevAction, cps_const.ActionUpdate)
}

func (a *WalletHandler) handleCPSAction(ctx context.Context, maker cps_entitites.User, requestAction cps_const.RequestAction, curData, prevData interface{}, actionType cps_const.ActionType) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
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
