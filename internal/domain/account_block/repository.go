package account_block

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AccountBlockRepo interface {
	GetBranch(ctx context.Context, branchCode string, filterParams *constant.Filter) (*model.Branch, error)
	GetAllBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)
	EnableOrDisable(ctx context.Context, blockType string, codes []string, cpsAction model.CPSAction, requestType model.RequestAction) error
	AuthorizeEnableBranches(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeDisableBranches(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeEnableRegions(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeDisableRegions(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeEnableDistricts(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeDisableDistrict(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeEnableCities(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeDisableCities(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}
