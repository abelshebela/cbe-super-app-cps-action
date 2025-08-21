// Package repository provides interfaces for wallet data persistence and operations.
package wallet

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet Wallet) (*Wallet, error)
	FetchWalletByID(ctx context.Context, id string) (*Wallet, error)
	FetchWallet(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*Wallet], error)
	UpdateWallet(ctx context.Context, wallet Wallet) (*Wallet, error)
	DeleteWallet(ctx context.Context, id string) (*Wallet, error)
	EnableDisableWallet(ctx context.Context, id string, enable bool) (*Wallet, error)
	WalletNameExists(ctx context.Context, name string, code string, id *string) (bool, error)
}
