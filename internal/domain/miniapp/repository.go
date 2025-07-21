package miniapp

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

)

type MiniRepository interface {
	CreateMiniAppAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error)
	DeleteMiniAppAction(ctx context.Context, action *model.CPSAction, id string) (*model.CPSAction, error)

	CreateMiniApp(ctx context.Context, miniapp *entities.CPSAction) (*model.MiniApp, error)

	GetMiniAppActionId(ctx context.Context, action_id string) (model.CPSAction, error)
	UpdateCpsAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error)

	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*model.MiniApp], error)

	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
}
