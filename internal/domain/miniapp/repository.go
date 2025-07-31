package miniapp

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type MiniRepository interface {
	CreateMiniApp(ctx context.Context, miniapp *MiniApp) (*MiniApp, error)
	UpdateMinApp(ctx context.Context, action *MiniApp) (*MiniApp, error)
	DeleteMiniAppAction(ctx context.Context, action *MiniApp) (*MiniApp, error)

	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*MiniApp, error)
	EnableDisableMiniApp(ctx context.Context, id string, enabled bool) (*MiniApp, error)
}
