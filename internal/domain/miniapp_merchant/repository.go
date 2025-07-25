package miniappmerchant

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type MiniAppMerchantRepository interface {
	CreateMiniAppMerchant(ctx context.Context, marchant *MiniAppMerchant) (*MiniAppMerchant, error)
	UpdateMiniAppMerchant(ctx context.Context, marchant *MiniAppMerchant) (*MiniAppMerchant, error)
	ListMiniAppMerchant(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*MiniAppMerchant], error)
	DetailMiniAppByID(ctx context.Context, id string) (*MiniAppMerchant, error)
	EnableMiniAppMerchant(ctx context.Context, id string) (*MiniAppMerchant, error)
	DisableMiniAppMerchant(ctx context.Context, id string) (*MiniAppMerchant, error)
	DeleteMiniAppMerchant(ctx context.Context, id string) (*MiniAppMerchant, error)
	MiniAppMerchantInfoExists(ctx context.Context, data CheckMiniAppMerchant, opts *MiniAppMerchantExistOptions) (bool, error)
}
