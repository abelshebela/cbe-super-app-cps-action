package wallet

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error)
	FetchWalletByID(ctx context.Context, id string) (*entity.Wallet, error)
	FetchWallet(ctx context.Context, filterParam *constant.MongoFilter) (*common_util.PaginatedResponse[[]*entity.Wallet], error)
	UpdateWallet(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error)
	DeleteWallet(ctx context.Context, id string) (*entity.Wallet, error)
	EnableDisableWallet(ctx context.Context, id string, enable bool) (*entity.Wallet, error)
	WalletNameExists(ctx context.Context, name string, code string, id *string) (bool, error)
}
