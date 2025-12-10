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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accountBlockService struct {
	repo       storage.AccountBlockRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewAccountService(repo storage.AccountBlockRepository, cpsService service.CPSActionService, logger utils.Logger) service.AccountBlockService {
	return &accountBlockService{repo: repo, cpsService: cpsService, logger: logger}
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
	s.logger.Infof("[EnableOrDisableBranches] processing %d branches, enabled: %v", len(branchIds), enabled)
	var alreadyEnabled []string
	var alreadyDisabled []string

	for _, id := range branchIds {
		branch, err := s.repo.GetBranchByIds(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorBranchNotFound.Code {
				s.logger.Errorf("[EnableOrDisableBranches] branch not found: %s", id)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			s.logger.Errorf("[EnableOrDisableBranches] failed to get branch %s: %v", id, err)
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

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[EnableOrDisableBranches] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableBranches] CPS action created successfully for %d branches", len(branchIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableRegions(ctx context.Context, regionIds []string, reason string, enabled bool) error {
	s.logger.Infof("[EnableOrDisableRegions] processing %d regions, enabled: %v", len(regionIds), enabled)
	for _, id := range regionIds {
		region, err := s.repo.GetRegionByIds(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorRegionNotFound.Code {
				s.logger.Errorf("[EnableOrDisableRegions] region not found: %s", id)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			s.logger.Errorf("[EnableOrDisableRegions] failed to get region %s: %v", id, err)
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

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[EnableOrDisableRegions] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableRegions] CPS action created successfully for %d regions", len(regionIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableDistricts(ctx context.Context, districtIds []string, reason string, enabled bool) error {
	s.logger.Infof("[EnableOrDisableDistricts] processing %d districts, enabled: %v", len(districtIds), enabled)
	for _, id := range districtIds {
		district, err := s.repo.GetDistrictById(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorDistrictNotFound.Code {
				s.logger.Errorf("[EnableOrDisableDistricts] district not found: %s", id)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			s.logger.Errorf("[EnableOrDisableDistricts] failed to get district %s: %v", id, err)
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

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[EnableOrDisableDistricts] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableDistricts] CPS action created successfully for %d districts", len(districtIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableCities(ctx context.Context, citiesCode []string, reason string, enabled bool) error {
	s.logger.Infof("[EnableOrDisableCities] processing %d cities, enabled: %v", len(citiesCode), enabled)
	for _, code := range citiesCode {
		city, err := s.repo.FindCityByID(ctx, code)
		if err != nil {
			if err.Error() == localization.ErrorCityNotFound.Code {
				s.logger.Errorf("[EnableOrDisableCities] city not found: %s", code)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			s.logger.Errorf("[EnableOrDisableCities] failed to get city %s: %v", code, err)
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

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[EnableOrDisableCities] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableCities] CPS action created successfully for %d cities", len(citiesCode))
	return nil
}

func (s *accountBlockService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[Authorize] authorizing account block action: %s", action.RequestAction)
	switch constants.RequestAction(action.RequestAction) {
	case constants.RequestEnableBranches:

		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal enable branches action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableBranch(ctx, id, action.Reason, true)

			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable branch %s: %v", id, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d branches", len(action.Codes))

	case constants.RequestDisableBranches:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal disable branches action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableBranch(ctx, id, action.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable branch %s: %v", id, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d branches", len(action.Codes))

	case constants.RequestEnableRegions:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal enable regions action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableRegion(ctx, id, action.Reason, true)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable region %s: %v", id, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d regions", len(action.Codes))

	case constants.RequestDisableRegions:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal disable regions action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableRegion(ctx, id, action.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable region %s: %v", id, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d regions", len(action.Codes))

	case constants.RequestEnableDistricts:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal enable districts action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableDistrict(ctx, id, action.Reason, true)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable district %s: %v", id, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d districts", len(action.Codes))

	case constants.RequestDisableDistricts:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal disable districts action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, id := range action.Codes {
			err = s.repo.EnableOrDisableDistrict(ctx, id, action.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable district %s: %v", id, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d districts", len(action.Codes))

	case constants.RequestEnableCities:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal enable cities action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, cityCode := range action.Codes {
			err = s.repo.EnableOrDisableCity(ctx, cityCode, action.Reason, true)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable city %s: %v", cityCode, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d cities", len(action.Codes))

	case constants.RequestDisableCities:
		action, err := local_util.JsonUnmarshal[model.EnableDisableAction](action.CurrentAction)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to unmarshal disable cities action: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		for _, cityCode := range action.Codes {
			err = s.repo.EnableOrDisableCity(ctx, cityCode, action.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable city %s: %v", cityCode, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d cities", len(action.Codes))

	default:
		s.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	action.ActionStatus = string(constants.ActionApproved)
	s.logger.Infof("[Authorize] account block action authorized successfully: %s", action.RequestAction)
	return action, nil
}
