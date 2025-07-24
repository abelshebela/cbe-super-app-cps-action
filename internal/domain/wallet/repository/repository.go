// Package repository provides interfaces for wallet data persistence and operations.
package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	UpdateWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	DeleteWallet(ctx context.Context, id string, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	GetWallet(ctx context.Context, id string) (*entity.Wallet, error)
	GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*error_codes.PaginatedResponse[[]*entity.Wallet], error)
	AuthorizeCreate(ctx context.Context, cpsAction *entities.CPSAction, action entity.Wallet) (*entities.CPSAction, error)
	AuthorizeUpdate(ctx context.Context, cpsAction *entities.CPSAction, action, prev entity.Wallet) (*entities.CPSAction, error)
	AuthorizeDelete(ctx context.Context, cpsAction *entities.CPSAction, action entity.Wallet) (*entities.CPSAction, error)
	EnableOrDisableWallet(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
	ExtractActionData(cpsAction *entities.CPSAction) (action entity.Wallet, prev entity.Wallet, err error)
	CheckWalletExists(ctx context.Context, wallet entity.CheckWallet) (bool, error)
	CheckIsEnabled(ctx context.Context, id string) (bool, error)
}
