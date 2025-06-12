package services

import (
	"context"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/repository"
)

type branchService struct {
	repo repository.BulkCustomerRepo
}

func NewBranchService(repo repository.BulkCustomerRepo) BranchServices {
	return &branchService{repo: repo}
}

func (s *branchService) FilterSingleBranches(ctx context.Context, region, district string) ([]entities.Branch, error) {
	if region == "" || district == "" {
		return nil, fmt.Errorf("region and district are required")
	}
	return s.repo.FilterSingleBranches(ctx, region, district)
}



func (s *branchService) DisableSingleBranch(ctx context.Context, branchCode string, cpsData entities.CPSAction) error {
    if branchCode == "" {
        return fmt.Errorf("branchCode is required")
    }
    return s.repo.DisableSingleBranch(ctx, branchCode, cpsData)
}


func (s *branchService) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (s *branchService) FilterMultipleBranches(ctx context.Context, region, district string) ([]string, error) {
	if region == "" || district == "" {
		return nil, fmt.Errorf("region and district are required")
	}
	return s.repo.FilterMultipleBranches(ctx, region, district)
}

func (s *branchService) DisableMultipleBranches(ctx context.Context, branchCodes []string) (*entities.CPSAction, error) {
	if len(branchCodes) == 0 {
		return nil, fmt.Errorf("branchCodes and cpsData are required")
	}
	return s.repo.DisableMultipleBranches(ctx, branchCodes)
}

func (s *branchService) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
}
