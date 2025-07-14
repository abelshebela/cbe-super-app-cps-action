// Package repository provides interfaces for wallet data persistence and operations.
package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	UpdateWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	DeleteWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	GetWallet(ctx context.Context, id string) (*entity.Wallet, error)
	GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error)
	EnableOrDisableWallet(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
}
