package budget

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type Repository interface {
	CreateIconAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	FetchIcons(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Icon], error)
	UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	CheckColorExist(ctx context.Context, color string) (bool, error)
	CreateAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	CreateColor(ctx context.Context, color string, action entities.CPSAction) (*entities.CPSAction, error)
	ListAllColor(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Color], error)
	GetByIDColor(ctx context.Context, id string) (*entities.Color, error)
	UpdateColor(ctx context.Context, color entities.CPSAction) (*entities.CPSAction, error)
	ApproveAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	Authorize(ctx context.Context, cpsAction *cps_entities.CPSAction) (*cps_entities.CPSAction, error)
}
