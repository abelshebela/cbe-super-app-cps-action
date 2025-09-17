package miniapp

import (
	"context"

	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type MiniRepository interface {
	CreateMiniApp(context.Context, *miniApp_domain.MiniApp) (*miniApp_domain.MiniApp, error)
	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*miniApp_domain.MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*miniApp_domain.MiniApp, error)
	UpdateMinApp(ctx context.Context, action *miniApp_domain.MiniApp) (*miniApp_domain.MiniApp, error)
	DeleteMiniAppAction(ctx context.Context, action *miniApp_domain.MiniApp) (*miniApp_domain.MiniApp, error)
	EnableDisableMiniApp(ctx context.Context, id string, enabled bool) (*miniApp_domain.MiniApp, error)
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
