package account_block

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type AccountBlockRepo interface {
	FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	UpdateRegion(ctx context.Context, region action.Region) error
	BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

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
