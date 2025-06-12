package outbound

import (
    "context"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/entities"
)

type BranchOutboundPort interface {
    FilterSingleBranches(ctx context.Context, region, district string) ([]string, error)
    DisableSingleBranch(ctx context.Context, branchCode, cpsData string) (*entities.CPSAction, error)
    ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

    FilterMultipleBranches(ctx context.Context, region, district string) ([]string, error)
    DisableMultipleBranches(ctx context.Context, branchCodes []string) (*entities.CPSAction, error)
    ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error
}