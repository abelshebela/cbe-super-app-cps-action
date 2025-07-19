package account_block

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type ApplicationServices interface {
	FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error
	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

	BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	UpdateRegion(ctx context.Context, region action.Region) error

	BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)

	ApproveBlockDistrict(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
	BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	ApproveBlockCity(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

	GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error)
	BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error)
	ApproveBlockUser(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
}
