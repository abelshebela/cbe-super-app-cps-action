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
	"fmt"

	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type accountBlockService struct {
	repo       storage.AccountBlockRepository
	cpsService service.CPSActionService
}

func NewAccountService(repo storage.AccountBlockRepository, cpsService service.CPSActionService) service.AccountBlockService {
	return &accountBlockService{repo: repo, cpsService: cpsService}
}

func (s *accountBlockService) GetBranchByCode(ctx context.Context, branchCode string) (*model.AccountBlock, error) {
	return s.repo.GetBranchByCode(ctx, branchCode)
}

func (s *accountBlockService) GetRegionByCode(ctx context.Context, regionCode string) (*model.AccountBlock, error) {
	return s.repo.GetRegionByCode(ctx, regionCode)
}

func (s *accountBlockService) GetDistrictByCode(ctx context.Context, districtCode string) (*model.AccountBlock, error) {
	return s.repo.GetDistrictById(ctx, districtCode)
}

func (s *accountBlockService) GetCityByCode(ctx context.Context, cityCode string) (*model.AccountBlock, error) {
	return s.repo.GetCityByCode(ctx, cityCode)
}

func (s *accountBlockService) GetAllBranches(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	return s.repo.FindAllBranchesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllRegions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	return s.repo.FindAllRegionsWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllDistricts(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	return s.repo.FindAllDistrictsWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllCities(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	return s.repo.FindAllCitiesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) EnableOrDisableBranches(ctx context.Context, branchIds []string, reason string, enabled bool) error {
	var alreadyEnabled []string
	var alreadyDisabled []string

	for _, id := range branchIds {
		branch, err := s.repo.GetBranchByIds(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorBranchNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if branch.IsEnabled == enabled {
			if enabled {
				alreadyEnabled = append(alreadyEnabled, id)
				// return errors.New(localization.ErrorAlreadyEnabled.Code)
			}
			alreadyDisabled = append(alreadyDisabled, id)
			// return errors.New(localization.ErrorAlreadyDisabled.Code)
		}

		if enabled {
			district, err := s.repo.GetDistrictById(ctx, branch.DistrictID)
			if err != nil {
				return err
			}
			if !district.IsEnabled {
				return errors.New(localization.ErrorCannotEnableBranch.Code)
			}
		}

	}

	if len(alreadyEnabled) > 0 {
		return fmt.Errorf("these branches are already enabled: %s", alreadyEnabled)
	} else if len(alreadyDisabled) > 0 {
		return fmt.Errorf("these branches are already disabled: %s", alreadyDisabled)
	}

	var actionType constants.ActionType
	var requestActionType constants.RequestAction
	if enabled {
		actionType = constants.ActionType(cps_constants.ActionEnable)
		requestActionType = constants.RequestEnableBranches
	} else {
		actionType = constants.ActionType(cps_constants.ActionDisable)
		requestActionType = constants.RequestDisableBranches
	}

	cpsAction := core.GenerateCPSAction(ctx, "BRANCH", enabled, branchIds, reason, actionType, requestActionType)

	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *accountBlockService) EnableOrDisableRegions(ctx context.Context, regionIds []string, reason string, enabled bool) error {
	for _, id := range regionIds {
		region, err := s.repo.GetRegionByIds(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorRegionNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if region.IsEnabled == enabled {
			if enabled {
				return errors.New(localization.ErrorAlreadyEnabled.Code)
			}
			return errors.New(localization.ErrorAlreadyDisabled.Code)
		}
	}

	var actionType constants.ActionType
	var requestActionType constants.RequestAction
	if enabled {
		actionType = constants.ActionType(cps_constants.ActionEnable)
		requestActionType = constants.RequestEnableRegions
	} else {
		actionType = constants.ActionType(cps_constants.ActionDisable)
		requestActionType = constants.RequestDisableRegions
	}

	cpsAction := core.GenerateCPSAction(ctx, "REGION", enabled, regionIds, reason, actionType, requestActionType)

	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *accountBlockService) EnableOrDisableDistricts(ctx context.Context, districtIds []string, reason string, enabled bool) error {
	for _, id := range districtIds {
		district, err := s.repo.GetDistrictById(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorDistrictNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if district.IsEnabled == enabled {
			if enabled {
				return errors.New(localization.ErrorAlreadyEnabled.Code)
			}
			return errors.New(localization.ErrorAlreadyDisabled.Code)
		}

		if enabled {
			region, err := s.repo.GetRegionByIds(ctx, district.RegionID)
			if err != nil {
				return err
			}
			if !region.IsEnabled {
				return errors.New(localization.ErrorCannotEnableDistrict.Code)
			}
		}
	}

	var actionType constants.ActionType
	var requestActionType constants.RequestAction
	if enabled {
		actionType = constants.ActionType(cps_constants.ActionEnable)
		requestActionType = constants.RequestEnableDistricts
	} else {
		actionType = constants.ActionType(cps_constants.ActionDisable)
		requestActionType = constants.RequestDisableDistricts
	}

	cpsAction := core.GenerateCPSAction(ctx, "DISTRICT", enabled, districtIds, reason, actionType, requestActionType)

	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *accountBlockService) EnableOrDisableCities(ctx context.Context, citiesCode []string, reason string, enabled bool) error {
	for _, code := range citiesCode {
		city, err := s.repo.GetCityByCode(ctx, code)
		if err != nil {
			if err.Error() == localization.ErrorCityNotFound.Code {
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if city.IsEnabled == enabled {
			if enabled {
				return errors.New(localization.ErrorAlreadyEnabled.Code)
			}
			return errors.New(localization.ErrorAlreadyDisabled.Code)
		}
	}

	var actionType constants.ActionType
	var requestActionType constants.RequestAction
	if enabled {
		actionType = constants.ActionType(cps_constants.ActionEnable)
		requestActionType = constants.RequestEnableCities
	} else {
		actionType = constants.ActionType(cps_constants.ActionDisable)
		requestActionType = constants.RequestDisableCities
	}

	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, citiesCode, reason, actionType, requestActionType)

	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *accountBlockService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	switch constants.RequestAction(action.RequestAction) {
	case constants.RequestEnableBranches:

		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, branchCode := range action.Codes {
			err = s.repo.EnableOrDisableBranch(ctx, branchCode, action.Reason, true)

			if err != nil {
				return nil, err
			}
		}

	case constants.RequestDisableBranches:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, branchCode := range action.Codes {
			err = s.repo.EnableOrDisableBranch(ctx, branchCode, action.Reason, false)
			if err != nil {
				return nil, err
			}
		}

	case constants.RequestEnableRegions:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableRegion(ctx, id, action.Reason, true)
			if err != nil {
				return nil, err
			}
		}

	case constants.RequestDisableRegions:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableRegion(ctx, id, action.Reason, false)
			if err != nil {
				return nil, err
			}
		}

	case constants.RequestEnableDistricts:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableDistrict(ctx, id, action.Reason, true)
			if err != nil {
				return nil, err
			}
		}

	case constants.RequestDisableDistricts:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableDistrict(ctx, id, action.Reason, false)
			if err != nil {
				return nil, err
			}
		}

	case constants.RequestEnableCities:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, cityCode := range action.Codes {
			err = s.repo.EnableOrDisableCity(ctx, cityCode, action.Reason, true)
			if err != nil {
				return nil, err
			}
		}

	case constants.RequestDisableCities:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			return nil, err
		}

		for _, cityCode := range action.Codes {
			err = s.repo.EnableOrDisableCity(ctx, cityCode, action.Reason, false)
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
