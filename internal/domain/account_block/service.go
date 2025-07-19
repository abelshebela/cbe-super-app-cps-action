package account_block

import (
	"context"
	"fmt"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type AccountService struct {
	repo AccountBlockRepo
}

func NewAccountService(repo AccountBlockRepo) ApplicationServices {
	return &AccountService{repo: repo}
}

func (s *AccountService) FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	if region == "" || district == "" {
		return nil, fmt.Errorf("region and district are required")
	}
	return s.repo.FilterSingleBranches(ctx, region, district, filterParams)
}

func (s *AccountService) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	if branch.BranchCode == "" {
		return "", fmt.Errorf("branchCode is required")
	}
	return s.repo.DisableSingleBranch(ctx, branch, maker)
}
func (s *AccountService) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (s *AccountService) FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}
	return s.repo.FilterMultipleBranches(ctx, region, district, filterParams)
}
func (s *AccountService) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	if len(branches) == 0 {
		return "", fmt.Errorf("branches list is empty")
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
func (s *AccountService) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
	if regionCode == "" {
		return action.Region{}, fmt.Errorf("regionCode is required")
	}
	return s.repo.GetRegionByCode(ctx, regionCode)
}
func (s *AccountService) BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error) {
	if regionCode == "" {
		return "", fmt.Errorf("regionCode is required")
	}
	actionCode, err := s.repo.BlockRegion(ctx, regionCode, maker)
	if err != nil {
		return "", err
	}
	return actionCode, nil
}

func (s *AccountService) ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveRegionBlock(ctx, actionID, approve, reason, checker)
}
func (s *AccountService) UpdateRegion(ctx context.Context, region action.Region) error {
	if region.RegionCode == "" {
		return fmt.Errorf("regionCode is required")
	}
	return s.repo.UpdateRegion(ctx, region)
}
func (s *AccountService) BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error) {
	if districtCode == "" {
		return "", fmt.Errorf("districtCode is required")
	}
	actionCode, err := s.repo.BlockDistrict(ctx, districtCode, maker)
	if err != nil {
		return "", err
	}
	return actionCode, nil
}
func (s *AccountService) GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error) {
	if districtCode == "" {
		return action.District{}, fmt.Errorf("districtCode is required")
	}
	return s.repo.GetDistrictByCode(ctx, districtCode)
}
func (s *AccountService) ApproveBlockDistrict(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	if actionID == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveBlockDistrict(ctx, actionID, approve, reason, checker)
}
func (s *AccountService) BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error) {
	if cityCode == "" {
		return "", fmt.Errorf("cityCode is required")
	}
	actionCode, err := s.repo.BlockCity(ctx, cityCode, maker)
	if err != nil {
		return "", err
	}
	return actionCode, nil
}

func (s *AccountService) GetCityByCode(ctx context.Context, cityCode string) (action.City, error) {
	if cityCode == "" {
		return action.City{}, fmt.Errorf("cityCode is required")
	}
	return s.repo.GetCityByCode(ctx, cityCode)
}
func (s *AccountService) ApproveBlockCity(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	if strings.TrimSpace(actionID) == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveBlockCity(ctx, actionID, approve, reason, checker)
}
func (s *AccountService) GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error) {
	if strings.TrimSpace(phoneNumber) == "" {
		return member.User{}, fmt.Errorf("phoneNumber is required")
	}
	return s.repo.GetUserByPhone(ctx, phoneNumber, maker)
}

func (s *AccountService) BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("userID is required")
	}
	return s.repo.BlockUser(ctx, userID, maker)
}

func (s *AccountService) ApproveBlockUser(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	if strings.TrimSpace(actionID) == "" {
		return fmt.Errorf("actionID is required")
	}
	return s.repo.ApproveBlockUser(ctx, actionID, approve, reason, checker)
}
