package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulkcustomer/entities"
)

type BranchOutboundPort interface {
	GetBranch(ctx context.Context, region, district string) ([]string, error)
	DisableSingleBranch(ctx context.Context, branchCode, cpsData string) (*entities.CPSAction, error)
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	GetAllBranches(ctx context.Context, region, district string) ([]entities.Branch, error)
	DisableMultipleBranches(ctx context.Context, branches []entities.Branch, maker entities.User) (*entities.CPSAction, error)
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error
}
