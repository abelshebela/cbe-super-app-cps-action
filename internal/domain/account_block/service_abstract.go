package account_block

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ApplicationServices interface {
	GetBranch(ctx context.Context, branchCode string, filterParams *constant.Filter) (*model.Branch, error)
	GetAllBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)
	EnableBranches(ctx context.Context, branchCodes []string, enabled bool) error
	DisableBranches(ctx context.Context, branchCodes []string, enabled bool) error
	EnableRegion(ctx context.Context, regionsCode []string, enabled bool) error
	DisableRegion(ctx context.Context, regionsCode []string, enabled bool) error
	EnableDistrict(ctx context.Context, districtsCode []string, enabled bool) error
	DisableDistrict(ctx context.Context, districtsCode []string, enabled bool) error
	EnableCity(ctx context.Context, citiesCode []string, enabled bool) error
	DisableCity(ctx context.Context, citiesCode []string, enabled bool) error
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}
