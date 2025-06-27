package account_block

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type AccountBlockOutboundPort interface {
	FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

	BlockRegion(ctx context.Context, regionID string, maker action.CPSAction) error
	UpdateRegion(ctx context.Context, region action.Region) error
	ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string) error
}
