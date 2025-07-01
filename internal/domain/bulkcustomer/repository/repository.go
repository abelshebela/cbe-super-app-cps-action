package repository

import (
	"context"

	"cbe-super-app-cps-action/internal/domain/bulkcustomer/entities"
)

type BulkCustomerRepo interface {
	FilterSingleBranches(ctx context.Context, region, district string) ([]entities.Branch, error)
	DisableSingleBranch(ctx context.Context, branch entities.Branch, maker entities.User) error
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string) ([]entities.Branch, error)
	DisableMultipleBranches(ctx context.Context, branches []entities.Branch, maker entities.User) error
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error
	GetBranchByCode(ctx context.Context, branchCode string) (entities.Branch, error)
}
