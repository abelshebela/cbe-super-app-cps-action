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

type AccountBlockOutboundPort interface {
	FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)
	AuthorizeSingleBranchDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeSingleBranchEnable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	AuthorizeBulkBranchesDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeBulkBranchesEnable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

	BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	AuthorizeRegionBlock(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	UpdateRegion(ctx context.Context, region action.Region) error
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)

	GetAllDistricts(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)
	AuthorizeBlockDistrict(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	AuthorizeBlockCity(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	GetAllBranches(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)

	BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error)
	GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error)
	AuthorizeBlockUser(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
}
