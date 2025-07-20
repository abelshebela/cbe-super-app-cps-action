package account_block

import (
	"context"
	"fmt"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type AccountService struct {
	repo AccountBlockRepo
}

func NewAccountService(repo AccountBlockRepo) ApplicationServices {
	return &AccountService{repo: repo}
}

func (s *AccountService) FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	if region == "" || district == "" {
		return nil, fmt.Errorf("region and district are required")
	}
	return s.repo.FilterSingleBranches(ctx, region, district)
}

func (s *AccountService) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	if branch.BranchCode == "" {
		return "", fmt.Errorf("branchCode is required")
	}
	return s.repo.DisableSingleBranch(ctx, branch, maker)
}
func (s *AccountService) AuthorizeSingleBranchDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	if cpsAction == nil {
		return nil, fmt.Errorf("actionID is required")
	}
	return s.repo.AuthorizeSingleBranchDisable(ctx, cpsAction)
}

func (s *AccountService) FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}
	return s.repo.FilterMultipleBranches(ctx, region, district)
}
func (s *AccountService) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	if len(branches) == 0 {
		return "", fmt.Errorf("branches list is empty")
	}
	return s.repo.DisableMultipleBranches(ctx, branches, maker)
}

func (s *AccountService) AuthorizeBulkBranchesDisable(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	if cpsAction == nil {
		return nil, fmt.Errorf("actionID is required")
	}
	return s.repo.AuthorizeBulkBranchesDisable(ctx, cpsAction)
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

func (s *AccountService) AuthorizeRegionBlock(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	if cpsAction == nil {
		return nil, fmt.Errorf("actionID is required")
	}
	return s.repo.AuthorizeRegionBlock(ctx, cpsAction)
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
func (s *AccountService) AuthorizeBlockDistrict(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	if cpsAction == nil {
		return nil, fmt.Errorf("actionID is required")
	}
	return s.repo.AuthorizeBlockDistrict(ctx, cpsAction)
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
func (s *AccountService) AuthorizeBlockCity(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	if cpsAction == nil {
		return nil, fmt.Errorf("actionID is required")
	}
	return s.repo.AuthorizeBlockCity(ctx, cpsAction)
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

func (s *AccountService) AuthorizeBlockUser(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	if cpsAction == nil {
		return nil, fmt.Errorf("the acction code must not be empty")
	}
	return s.repo.AuthorizeBlockUser(ctx, cpsAction)

}
