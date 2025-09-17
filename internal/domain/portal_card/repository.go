package portalcard

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type PortaCardInterface interface {
	GetAllPortalCard(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
}

type PortalCardRepository interface {
	GetAllPortalCard(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
}
