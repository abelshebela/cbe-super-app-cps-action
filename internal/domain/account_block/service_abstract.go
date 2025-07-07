package account_block

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type ApplicationServices interface {
	FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error
	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

	BlockRegion(ctx context.Context,  region action.Region, maker action.CPSAction) error
	UpdateRegion(ctx context.Context, region action.Region) error
	ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
	GetRegionByID(ctx context.Context, regionID string) (action.Region, error)

	BlockDistrict(ctx context.Context, districtID string, maker action.CPSAction) error
	GetDistrictByID(ctx context.Context, districtID string) (action.District, error)
	ApproveBlockDistrict(ctx context.Context, districtID string, checker action.CPSAction) error

	BlockCity(ctx context.Context, cityID string, maker action.CPSAction) error
	GetCityByID(ctx context.Context, cityID string) (action.City, error)
	ApproveBlockCity(ctx context.Context, cityID string, checker action.CPSAction) error

	BlockUser(ctx context.Context, userID string, maker action.CPSAction) error
	GetUserByID(ctx context.Context, userID string, maker action.CPSAction) (member.User, error)
	ApproveBlockUser(ctx context.Context, userID string, checker action.CPSAction) error
}
