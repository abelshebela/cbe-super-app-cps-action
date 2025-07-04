package services

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulkcustomer/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulkcustomer/repository"
)

type BranchService struct {
	repo repository.BulkCustomerRepo
}

func NewBranchService(repo repository.BulkCustomerRepo) *BranchService {
	return &BranchService{repo: repo}
}

func (s *BranchService) FilterSingleBranches(ctx context.Context, region, district string) ([]entities.Branch, error) {
	if region == "" || district == "" {
		return nil, fmt.Errorf("region and district are required")
	}
	return s.repo.FilterSingleBranches(ctx, region, district)
}

func (s *BranchService) DisableSingleBranch(ctx context.Context, branch entities.Branch, maker entities.User) error {
	return s.repo.DisableSingleBranch(ctx, branch, maker)
}

func (s *BranchService) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (s *BranchService) FilterMultipleBranches(ctx context.Context, region, district string) ([]entities.Branch, error) {
	if region == "" || district == "" {
		return nil, fmt.Errorf("region and district are required")
	}
	return s.repo.FilterMultipleBranches(ctx, region, district)
}
func (s *BranchService) DisableMultipleBranches(ctx context.Context, branches []entities.Branch, maker entities.User) (*entities.CPSAction, error) {
	if len(branches) == 0 {
		return nil, fmt.Errorf("branches are required")
	}
	err := s.repo.DisableMultipleBranches(ctx, branches, maker)
	if err != nil {
		return nil, err
	}
	action := &entities.CPSAction{}
	return action, nil
}
func (s *BranchService) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
}

func (s *BranchService) GetBranchByCode(ctx context.Context, branchCode string) (entities.Branch, error) {
	return s.repo.GetBranchByCode(ctx, branchCode)
}
