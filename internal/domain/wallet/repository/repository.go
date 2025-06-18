package repository

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/wallet/entity"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type Repository interface {
	CreateWallet(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	UpdateWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	DeleteWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	GetWallet(ctx context.Context, id string) (*entity.Wallet, error)
	GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error)
	EnableOrDisableWallet(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
}
