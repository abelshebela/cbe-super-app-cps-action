package account_block

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_block"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type ApplicationService interface {
	FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)
	BlockRegion(ctx context.Context, region action.Region, maker action.CPSAction) error
	UpdateRegion(ctx context.Context, region action.Region) error
    ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
	GetRegionByID(ctx context.Context, regionID string) (action.Region, error)
}

type Handler struct {
	repo account_block.AccountBlockRepo
}

func NewApplicationHandler(repo account_block.AccountBlockRepo) ApplicationService {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	return h.repo.FilterSingleBranches(ctx, region, district)
}

func (h *Handler) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error {
	return h.repo.DisableSingleBranch(ctx, branch, maker)
}

func (h *Handler) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	return h.repo.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (h *Handler) FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	return h.repo.FilterMultipleBranches(ctx, region, district)
}

func (h *Handler) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error {

    return h.repo.DisableMultipleBranches(ctx, branches, maker)
}
func (h *Handler) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	return h.repo.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
}

func (h *Handler) GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error) {
	return h.repo.GetBranchByCode(ctx, branchCode)
}

func (h *Handler) BlockRegion(ctx context.Context, region action.Region, maker action.CPSAction) error {

	region.Enabled = false
	region.UpdatedAt = time.Now()
	if err := h.repo.UpdateRegion(ctx, region); err != nil {
		return err
	}
	return h.repo.BlockRegion(ctx, region.ID, maker)
}

func (h *Handler) UpdateRegion(ctx context.Context, region action.Region) error {
	region.UpdatedAt = time.Now()
	return h.repo.UpdateRegion(ctx, region)
}

func (h *Handler) ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	return h.repo.ApproveRegionBlock(ctx, actionID, approve, reason, checker)
}
func (h *Handler) GetRegionByID(ctx context.Context, regionID string) (action.Region, error) {
    if regionID == "" {
        return action.Region{}, fmt.Errorf("regionID is required")
    }

    region, err := h.repo.GetRegionByID(ctx, regionID)
    if err != nil {
        return action.Region{}, fmt.Errorf("failed to get region by ID: %w", err)
    }

    return region, nil
}