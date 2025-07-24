package miniappmerchant

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type MiniAppMerchantRepository interface {
	CreateMiniAppMerchant(ctx context.Context, marchant *entities.MiniAppMerchant) (*entities.MiniAppMerchant, error)
	UpdateMiniAppMerchant(ctx context.Context, marchant *entities.MiniAppMerchant) (*entities.MiniAppMerchant, error)
	ListMiniAppMerchant(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*entities.MiniAppMerchant], error)
	DetailMiniAppByID(ctx context.Context, id string) (*entities.MiniAppMerchant, error)
	EnableMiniAppMerchant(ctx context.Context, id string) (*entities.MiniAppMerchant, error)
	DisableMiniAppMerchant(ctx context.Context, id string) (*entities.MiniAppMerchant, error)
	DeleteMiniAppMerchant(ctx context.Context, id string) (*entities.MiniAppMerchant, error)
	MiniAppMerchantInfoExists(ctx context.Context, data entities.CheckMiniAppMerchant) (bool, error)
}
