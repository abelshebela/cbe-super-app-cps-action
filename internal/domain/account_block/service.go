package account_block

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type AccountService struct {
	repo AccountBlockRepo
}

func NewAccountService(repo AccountBlockRepo) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	if region == "" || district == "" {
		return nil, fmt.Errorf("region and district are required")
	}
	return s.repo.FilterSingleBranches(ctx, region, district)
}

func (s *AccountService) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error {
	return s.repo.DisableSingleBranch(ctx, branch, maker)
}

func (s *AccountService) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (s *AccountService) FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}
	return s.repo.FilterMultipleBranches(ctx, region, district)
}

func (s *AccountService) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error {
	if len(branches) == 0 {
		return fmt.Errorf("branches list is empty")
	}
	return s.repo.DisableMultipleBranches(ctx, branches, maker)
}

func (s *AccountService) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
}

func (s *AccountService) GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error) {
	if branchCode == "" {
		return action.Branch{}, fmt.Errorf("branchCode is required")
	}
	return s.repo.GetBranchByCode(ctx, branchCode)
}

func (s *AccountService) BlockRegion(ctx context.Context, region action.Region, maker action.CPSAction) error {
	if region.ID == "" || region.RegionName == "" {
		return fmt.Errorf("region ID and name are required")
	}

	region.Enabled = false
	region.UpdatedAt = time.Now()

	if err := s.repo.UpdateRegion(ctx, region); err != nil {
		return fmt.Errorf("failed to disable region: %w", err)
	}

	branches, err := s.repo.FilterMultipleBranches(ctx, region.RegionName, "")
	if err != nil {
		return fmt.Errorf("failed to fetch branches under region: %w", err)
	}

	for i := range branches {
		branches[i].Enabled = false
		branches[i].UpdatedAt = time.Now()
	}

	if err := s.repo.DisableMultipleBranches(ctx, branches, maker.Maker); err != nil {
		return fmt.Errorf("failed to disable branches: %w", err)
	}

	if err := s.repo.BlockRegion(ctx, region.ID, maker); err != nil {
		return fmt.Errorf("failed to log region block action: %w", err)
	}

	return nil
}

func (s *AccountService) ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveRegionBlock(ctx, actionID, approve, reason)
}

func (s *AccountService) UpdateRegion(ctx context.Context, region action.Region) error {
	if region.ID == "" || region.RegionName == "" {
		return fmt.Errorf("region ID and name are required")
	}
	region.UpdatedAt = time.Now()
	return s.repo.UpdateRegion(ctx, region)
}
