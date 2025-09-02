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

func (s *accountBlockService) EnableOrDisableBranches(ctx context.Context, branchCodes []string, enabled bool) error {

	for _, code := range branchCodes {
		branch, err := s.repo.GetBranchByCode(ctx, code)
		if err != nil {
			if err.Error() == localization.ErrorBranchNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if branch.Enabled == enabled {
			if enabled {
				return errors.New(localization.ErrorBranchAlreadyEnabled.Code)
			} else {
				return errors.New(localization.ErrorBranchAlreadyDisabled.Code)
			}
		}
	}

	cpsAction := core.GenerateCPSAction(ctx, "BRANCH", enabled, branchCodes, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableBranches)

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		return err
	}

	return nil
}

func (s *accountBlockService) EnableOrDisableRegions(ctx context.Context, regionsCode []string, enabled bool) error {
	for _, code := range regionsCode {
		branch, err := s.repo.GetRegionByCode(ctx, code)
		if err != nil {
			if err.Error() == localization.ErrorRegionNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if branch.Enabled == enabled {
			if enabled {
				return errors.New(localization.ErrorRegionAlreadyEnabled.Code)
			} else {
				return errors.New(localization.ErrorRegionAlreadyDisabled.Code)
			}
		}
	}

	cpsAction := core.GenerateCPSAction(ctx, "REGION", enabled, regionsCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableRegions)

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		return err
	}

	return nil
}

func (s *accountBlockService) EnableOrDisableDistricts(ctx context.Context, districtsCode []string, enabled bool) error {
	for _, code := range districtsCode {
		branch, err := s.repo.GetDistrictByCode(ctx, code)
		if err != nil {
			if err.Error() == localization.ErrorDistrictNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if branch.Enabled == enabled {
			if enabled {
				return errors.New(localization.ErrorDistrictAlreadyEnabled.Code)
			} else {
				return errors.New(localization.ErrorDistrictAlreadyDisabled.Code)
			}
		}
	}

	cpsAction := core.GenerateCPSAction(ctx, "DISTRICT", enabled, districtsCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableDistricts)

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		return err
	}

	return nil
}

func (s *accountBlockService) EnableOrDisableCities(ctx context.Context, citiesCode []string, enabled bool) error {

	for _, code := range citiesCode {
		branch, err := s.repo.GetCityByCode(ctx, code)
		if err != nil {
			if err.Error() == localization.ErrorCityNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if branch.Enabled == enabled {
			if enabled {
				return errors.New(localization.ErrorCityAlreadyEnabled.Code)
			} else {
				return errors.New(localization.ErrorCityAlreadyDisabled.Code)
			}
		}
	}

	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, citiesCode, constants.ActionType(cps_constants.ActionEnable), constants.RequestEnableCities)

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		return err
	}

	return nil
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

	case constants.RequestEnableRegions:
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

	case constants.RequestDisableRegions:
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

	case constants.RequestEnableDistricts:
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

	case constants.RequestDisableDistricts:
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

	case constants.RequestEnableCities:
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

	case constants.RequestDisableCities:
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

	default:
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	action.ActionStatus = string(constants.ActionApproved)
	return action, nil
}
