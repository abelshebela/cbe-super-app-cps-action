package miniapp

import (
	"context"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

)

type MiniRepository interface {
	CreateMiniApp(context.Context, *entities.CPSAction) (*miniApp_domain.MiniApp, error)
	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*miniApp_domain.MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*miniApp_domain.MiniApp, error)
	UpdateMinApp(ctx context.Context, action *entities.CPSAction) (*miniApp_domain.MiniApp, error)
	DeleteMiniAppAction(ctx context.Context, action *entities.CPSAction) (*miniApp_domain.MiniApp, error)
	// ExtractActionData(currentAction interface{}) (*miniApp_domain.MiniApp, error)
}
