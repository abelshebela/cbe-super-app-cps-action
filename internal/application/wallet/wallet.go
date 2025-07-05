package wallet

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/service"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletHandlerService interface {
	GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error)
	GetWallet(ctx context.Context, id string) (*entity.Wallet, error)
	CreateWallet(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error)
	UpdateWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	DeleteWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error)
	EnableOrDisableWallet(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
}

type WalletHandler struct {
	walletDomain service.WalletService
	logger       utils.Logger
}

func InitWalletHanlder(walletDomain service.WalletService, logger utils.Logger) WalletHandlerService {
	return &WalletHandler{
		walletDomain: walletDomain,
		logger:       logger,
	}
}

func (w *WalletHandler) CreateWallet(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error) {
	cpsRes, err := w.walletDomain.CreateWallet(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsRes, nil
}

func (w *WalletHandler) DeleteWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error) {
	cpsAction, err := w.walletDomain.DeleteWallet(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletHandler) GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error) {
	banks, err := w.walletDomain.GetAllWallet(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (w *WalletHandler) GetWallet(ctx context.Context, id string) (*entity.Wallet, error) {
	bank, err := w.walletDomain.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (w *WalletHandler) UpdateWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error) {
	cpsAction, err := w.walletDomain.UpdateWallet(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletHandler) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error) {
	cpsAction, err := w.walletDomain.Authorize(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletHandler) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {
	cpsAction, err := w.walletDomain.Reject(ctx, req)
	if err != nil {
		return nil, err
	}
	return cpsAction, nil
}

func (w *WalletHandler) EnableOrDisableWallet(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	cpsAction, err := w.walletDomain.EnableOrDisableWallet(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}
