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

	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type accountBlockService struct {
	repo       storage.AccountBlockRepository
	cpsService service.CPSActionService
}

func NewAccountService(repo storage.AccountBlockRepository, cpsService service.CPSActionService) service.AccountBlockService {
	return &accountBlockService{repo: repo, cpsService: cpsService}
}

func (s *accountBlockService) GetBranchByCode(ctx context.Context, branchCode string) (*model.Branch, error) {
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

func (s *accountBlockService) GetAllBranches(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Branch], error) {
	return s.repo.FindAllBranchesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllRegions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Region], error) {
	return s.repo.FindAllRegionsWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllDistricts(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.District], error) {
	return s.repo.FindAllDistrictsWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllCities(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.City], error) {
	return s.repo.FindAllCitiesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) EnableBranches(ctx context.Context, branchCodes []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "BRANCH", enabled, branchCodes, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableBranches)
	return s.enableOrDisable(ctx, "BRANCH", branchCodes, cpsAction, constants.RequestEnableBranches)
}

func (s *accountBlockService) DisableBranches(ctx context.Context, branchCodes []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "BRANCH", enabled, branchCodes, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableBranches)
	return s.enableOrDisable(ctx, "BRANCH", branchCodes, cpsAction, constants.RequestDisableBranches)
}

func (s *accountBlockService) EnableRegions(ctx context.Context, regionsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "REGION", enabled, regionsCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableRegion)
	return s.enableOrDisable(ctx, "REGION", regionsCode, cpsAction, constants.RequestEnableRegion)
}

func (s *accountBlockService) DisableRegions(ctx context.Context, regionsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "REGION", enabled, regionsCode, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableRegion)
	return s.enableOrDisable(ctx, "REGION", regionsCode, cpsAction, constants.RequestDisableRegion)
}

func (s *accountBlockService) EnableDistricts(ctx context.Context, districtsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "DISTRICT", enabled, districtsCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableDistrict)
	return s.enableOrDisable(ctx, "DISTRICT", districtsCode, cpsAction, constants.RequestEnableDistrict)
}

func (s *accountBlockService) DisableDistricts(ctx context.Context, districtsCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "DISTRICT", enabled, districtsCode, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableDistrict)
	return s.enableOrDisable(ctx, "DISTRICT", districtsCode, cpsAction, constants.RequestDisableDistrict)
}

func (s *accountBlockService) EnableCities(ctx context.Context, citiesCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, citiesCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableCity)
	return s.enableOrDisable(ctx, "CITY", citiesCode, cpsAction, constants.RequestEnableCity)
}

func (s *accountBlockService) DisableCities(ctx context.Context, citiesCode []string, enabled bool) error {
	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, citiesCode, constants.ActionType(cps_constants.ActionDisable), constants.RequestDisableCity)
	return s.enableOrDisable(ctx, "CITY", citiesCode, cpsAction, constants.RequestDisableCity)
}

func (s *accountBlockService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {

	switch constants.RequestAction(action.ActionType) {
	case constants.RequestEnableBranches:
		branches, err := local_util.JsonUnmarshal[[]model.Branch](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, branch := range *branches {
			err = s.repo.EnableOrDisableBranch(ctx, branch.ID.Hex(), true)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	case constants.RequestDisableBranches:
		branches, err := local_util.JsonUnmarshal[[]model.Branch](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, branch := range *branches {
			err = s.repo.EnableOrDisableBranch(ctx, branch.ID.Hex(), false)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	case constants.RequestEnableRegion:
		regions, err := local_util.JsonUnmarshal[[]model.Region](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, region := range *regions {
			err = s.repo.EnableOrDisableRegion(ctx, region.ID.Hex(), true)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	case constants.RequestDisableRegion:
		regions, err := local_util.JsonUnmarshal[[]model.Region](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, region := range *regions {
			err = s.repo.EnableOrDisableRegion(ctx, region.ID.Hex(), false)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	case constants.RequestEnableDistrict:
		districts, err := local_util.JsonUnmarshal[[]model.District](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, district := range *districts {
			err = s.repo.EnableOrDisableDistrict(ctx, district.ID.Hex(), true)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	case constants.RequestDisableDistrict:
		districts, err := local_util.JsonUnmarshal[[]model.District](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, district := range *districts {
			err = s.repo.EnableOrDisableDistrict(ctx, district.ID.Hex(), false)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	case constants.RequestEnableCity:
		cities, err := local_util.JsonUnmarshal[[]model.City](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, city := range *cities {
			err = s.repo.EnableOrDisableCity(ctx, city.ID.Hex(), true)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	case constants.RequestDisableCity:
		cities, err := local_util.JsonUnmarshal[[]model.City](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, city := range *cities {
			err = s.repo.EnableOrDisableCity(ctx, city.ID.Hex(), false)
			if err != nil {
				return nil, err
			}
		}
		return action, nil

	default:
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *accountBlockService) enableOrDisable(ctx context.Context, blockType string, codes []string, cpsAction model.CPSAction, requestType constants.RequestAction) error {
	// Check if duplicate action is requested
	switch blockType {
	case "BRANCH":
		var boolStatus bool
		switch requestType {
		case constants.RequestEnableBranches:
			boolStatus = true
		case constants.RequestDisableBranches:
			boolStatus = false
		}

		for _, code := range codes {
			branch, err := s.repo.GetBranchByCode(ctx, code)
			if err != nil {
				if err.Error() == localization.ErrorBranchNotFound.Code {
					continue
				}
				return errors.New(localization.ErrorUnexpectedError.Code)
			}

			if branch.Enabled == boolStatus {
				return errors.New(localization.ErrorDuplicateAction.Code)
			}
		}
	case "REGION":

		var boolStatus bool
		switch requestType {
		case constants.RequestEnableRegion:
			boolStatus = true
		case constants.RequestDisableRegion:
			boolStatus = false
		}

		for _, code := range codes {
			branch, err := s.repo.GetRegionByCode(ctx, code)
			if err != nil {
				if err.Error() == localization.ErrorRegionNotFound.Code {
					continue
				}
				return errors.New(localization.ErrorUnexpectedError.Code)
			}

			if branch.Enabled == boolStatus {
				return errors.New(localization.ErrorDuplicateAction.Code)
			}
		}
	case "DISTRICT":
		var boolStatus bool
		switch requestType {
		case constants.RequestEnableDistrict:
			boolStatus = true
		case constants.RequestDisableDistrict:
			boolStatus = false
		}

		for _, code := range codes {
			branch, err := s.repo.GetDistrictByCode(ctx, code)
			if err != nil {
				if err.Error() == localization.ErrorDistrictNotFound.Code {
					continue
				}
				return errors.New(localization.ErrorUnexpectedError.Code)
			}

			if branch.Enabled == boolStatus {
				return errors.New(localization.ErrorDuplicateAction.Code)
			}
		}
	case "CITY":
		var boolStatus bool
		switch requestType {
		case constants.RequestEnableCity:
			boolStatus = true
		case constants.RequestDisableCity:
			boolStatus = false
		}

		for _, code := range codes {
			branch, err := s.repo.GetCityByCode(ctx, code)
			if err != nil {
				if err.Error() == localization.ErrorCityNotFound.Code {
					continue
				}
				return errors.New(localization.ErrorUnexpectedError.Code)
			}

			if branch.Enabled == boolStatus {
				return errors.New(localization.ErrorDuplicateAction.Code)
			}
		}
	}

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		return err
	}

	return nil
}
