package account_block

import (
	"context"
	"fmt"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

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

func (s *AccountService) EnableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	if branch.BranchCode == "" {
		return "", fmt.Errorf("branchCode is required")
	}
	return s.repo.EnableSingleBranch(ctx, branch, maker)
}

func (s *AccountService) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	if branch.BranchCode == "" {
		return "", fmt.Errorf("branchCode is required")
	}
	return s.repo.DisableSingleBranch(ctx, branch, maker)
}

func (s *AccountService) FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}
	return s.repo.FilterMultipleBranches(ctx, region, district, filterParams)
}
func (s *AccountService) EnableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	if len(branches) == 0 {
		return "", fmt.Errorf("branches list is empty")
	}
	return s.repo.EnableMultipleBranches(ctx, branches, maker)
}

func (s *AccountService) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	if len(branches) == 0 {
		return "", fmt.Errorf("branches list is empty")
	}
	return s.repo.DisableMultipleBranches(ctx, branches, maker)
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

func (s *AccountService) EnableRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error) {
	if regionCode == "" {
		return "", fmt.Errorf("regionCode is required")
	}
	actionCode, err := s.repo.EnableRegion(ctx, regionCode, maker)
	if err != nil {
		return "", err
	}
	return actionCode, nil
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

func (s *AccountService) UpdateRegion(ctx context.Context, region action.Region) error {
	if region.RegionCode == "" {
		return fmt.Errorf("regionCode is required")
	}
	return s.repo.UpdateRegion(ctx, region)
}

func (s *AccountService) EnableDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error) {
	if districtCode == "" {
		return "", fmt.Errorf("districtCode is required")
	}
	actionCode, err := s.repo.EnableDistrict(ctx, districtCode, maker)
	if err != nil {
		return "", err
	}
	return actionCode, nil
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

func (s *AccountService) EnableCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error) {
	if cityCode == "" {
		return "", fmt.Errorf("cityCode is required")
	}
	actionCode, err := s.repo.EnableCity(ctx, cityCode, maker)
	if err != nil {
		return "", err
	}
	return actionCode, nil
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

func (s *AccountService) GetAllCities(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	return s.repo.GetAllCities(ctx, filterParams)
}

func (s *AccountService) GetAllDistricts(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	return s.repo.GetAllDistricts(ctx, filterParams)
}

func (s *AccountService) GetAllRegions(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	return s.repo.GetAllRegions(ctx, filterParams)
}

func (s *AccountService) GetAllBranches(ctx context.Context, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}
	return s.repo.GetAllBranches(ctx, filterParams)
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

func (s AccountService) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	reqAction := action.RequestAction

	switch reqAction {
	case cps_const.RequestBlockCity:
		return s.repo.AuthorizeBlockCity(ctx, action)
	case cps_const.RequestEnableCity:
		return s.repo.AuthorizeEnableCity(ctx, action)

	case cps_const.RequestBlockDistrict:
		return s.repo.AuthorizeBlockDistrict(ctx, action)
	case cps_const.RequestEnableDistrict:
		return s.repo.AuthorizeEnableDistrict(ctx, action)

	case cps_const.RequestBlockRegion:
		return s.repo.AuthorizeRegionBlock(ctx, action)
	case cps_const.RequestEnableRegion:
		return s.repo.AuthorizeRegionEnable(ctx, action)

	case cps_const.RequestDisableMultiBranches:
		return s.repo.AuthorizeBulkBranchesDisable(ctx, action)
	case cps_const.RequestEnableMultiBranches:
		return s.repo.AuthorizeBulkBranchesEnable(ctx, action)

	case cps_const.RequestDisableSingleBranch:
		return s.repo.AuthorizeSingleBranchDisable(ctx, action)
	case cps_const.RequestEnableSingleBranch:
		return s.repo.AuthorizeSingleBranchEnable(ctx, action)

	case cps_const.RequestBlockUser:
		return s.repo.AuthorizeBlockUser(ctx, action)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
}
