package repository

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/entities"
)

type BulkCustomerRepo interface {
	FilterSingleBranches(ctx context.Context, region, district string) ([]entities.Branch, error)
	DisableSingleBranch(ctx context.Context, branchCode string, cpsData entities.CPSAction) error
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string) ([]string, error)
	DisableMultipleBranches(ctx context.Context, branchCodes []string) (*entities.CPSAction, error)
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error
}
