package account_block

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type ApplicationServices interface {
	FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)
	EnableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)

	FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	EnableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

	EnableRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	UpdateRegion(ctx context.Context, region action.Region) error

	EnableDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)

	EnableCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)
	GetAllBranches(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)

	GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error)
	BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}
