package miniapp

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	// miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type Outbound interface {
	CreateMiniAppAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error)
	DeleteMiniAppAction(ctx context.Context, action model.CPSAction, id string) (model.CPSAction, error)

	CreateMiniApp(ctx context.Context, miniapp model.MiniApp) (model.CPSAction, error)
	GetMiniAppActionId(ctx context.Context, action_id string) (model.CPSAction, error)
	UpdateCpsAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error)
	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*model.MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
}
