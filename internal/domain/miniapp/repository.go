package miniapp

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type MiniRepository interface {
	CreateMiniApp(ctx context.Context, miniapp *entities.CPSAction) (*MiniApp, error)
	UpdateMinApp(ctx context.Context, action *entities.CPSAction) (*MiniApp, error)
	DeleteMiniAppAction(ctx context.Context, action *entities.CPSAction) (*MiniApp, error)

	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*MiniApp, error)
}
