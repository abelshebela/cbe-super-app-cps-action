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

type AccountBlockRepo interface {
	FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	EnableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)
	// AuthorizeSingleBranchDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	// AuthorizeSingleBranchEnable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)
	EnableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	// AuthorizeBulkBranchesDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	// AuthorizeBulkBranchesEnable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	UpdateRegion(ctx context.Context, region action.Region) error

	EnableRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	// AuthorizeRegionBlock(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	// AuthorizeRegionEnable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	EnableDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)
	// AuthorizeBlockDistrict(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	// AuthorizeEnableDistrict(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	EnableCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	// AuthorizeBlockCity(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	// AuthorizeEnableCity(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error)
	BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error)
	AuthorizeBlockUser(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	// AuthorizeEnableUser(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)
	GetAllBranches(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)

	// Newly added
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
