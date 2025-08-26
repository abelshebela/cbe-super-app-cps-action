package accountblock

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/account_block/core"
	cps_constants "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
)

type accountBlockService struct {
	repo storage.AccountBlockRepository
}

func NewAccountService(repo storage.AccountBlockRepository) service.AccountBlockService {
	return &accountBlockService{repo: repo}
}

func (s *accountBlockService) GetBranch(ctx context.Context, branchCode string, filterParams *types.Filter) (*model.Branch, error) {
	return s.repo.GetBranchByCode(ctx, branchCode)
}

func (s *accountBlockService) GetRegionByCode(ctx context.Context, regionCode string) (*model.Region, error) {
	return s.repo.GetRegionByCode(ctx, regionCode)
}

func (s *accountBlockService) GetDistrictByCode(ctx context.Context, districtCode string) (*model.District, error) {
	return s.repo.GetDistrictByCode(ctx, districtCode)
}

func (s *accountBlockService) GetCityByCode(ctx context.Context, cityCode string) (*model.City, error) {
	return s.repo.GetCityByCode(ctx, cityCode)
}

func (s *accountBlockService) GetAllBranches(ctx context.Context, region, district string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Branch], error) {
	// NEEDS FIX
	return s.repo.FindAllBranchesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllCities(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.City], error) {
	return s.repo.FindAllCitiesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllDistricts(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.District], error) {
	return s.repo.FindAllDistrictsWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllRegions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Region], error) {
	return s.repo.FindAllRegionsWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) EnableBranches(ctx context.Context, branchCodes []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "BRANCH", enabled, branchCodes, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableBranches)
	return s.repo.EnableOrDisable(ctx, "BRANCH", branchCodes, cpsAction, constants.RequestEnableBranches)
}

func (s *accountBlockService) DisableBranches(ctx context.Context, branchCodes []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "BRANCH", enabled, branchCodes, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableBranches)
	return s.repo.EnableOrDisable(ctx, "BRANCH", branchCodes, cpsAction, constants.RequestDisableBranches)
}

func (s *accountBlockService) EnableRegion(ctx context.Context, regionsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "REGION", enabled, regionsCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableRegion)
	return s.repo.EnableOrDisable(ctx, "REGION", regionsCode, cpsAction, constants.RequestEnableRegion)
}

func (s *accountBlockService) DisableRegion(ctx context.Context, regionsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "REGION", enabled, regionsCode, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableRegion)
	return s.repo.EnableOrDisable(ctx, "REGION", regionsCode, cpsAction, constants.RequestDisableRegion)
}

func (s *accountBlockService) EnableDistrict(ctx context.Context, districtsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "DISTRICT", enabled, districtsCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableDistrict)
	return s.repo.EnableOrDisable(ctx, "DISTRICT", districtsCode, cpsAction, constants.RequestEnableDistrict)
}

func (s *accountBlockService) DisableDistrict(ctx context.Context, districtsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "DISTRICT", enabled, districtsCode, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableDistrict)
	return s.repo.EnableOrDisable(ctx, "DISTRICT", districtsCode, cpsAction, constants.RequestDisableDistrict)
}

func (s *accountBlockService) EnableCity(ctx context.Context, citiesCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, citiesCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableCity)
	return s.repo.EnableOrDisable(ctx, "CITY", citiesCode, cpsAction, constants.RequestEnableCity)
}

func (s *accountBlockService) DisableCity(ctx context.Context, citiesCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, citiesCode, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableCity)
	return s.repo.EnableOrDisable(ctx, "CITY", citiesCode, cpsAction, constants.RequestDisableCity)
}

func (s *accountBlockService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {

	switch constants.RequestAction(action.ActionType) {
	case constants.RequestEnableBranches:
		return s.repo.AuthorizeEnableBranches(ctx, action)
	case constants.RequestDisableBranches:
		return s.repo.AuthorizeDisableBranches(ctx, action)
	case constants.RequestEnableRegion:
		return s.repo.AuthorizeEnableRegions(ctx, action)
	case constants.RequestDisableRegion:
		return s.repo.AuthorizeDisableRegions(ctx, action)
	case constants.RequestEnableDistrict:
		return s.repo.AuthorizeEnableDistricts(ctx, action)
	case constants.RequestDisableDistrict:
		return s.repo.AuthorizeDisableDistrict(ctx, action)
	case constants.RequestEnableCity:
		return s.repo.AuthorizeEnableCities(ctx, action)
	case constants.RequestDisableCity:
		return s.repo.AuthorizeDisableCities(ctx, action)
	default:
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}
