package accountblock

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/account_block/core"
	cps_constants "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type accountBlockService struct {
	repo       storage.AccountBlockRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewAccountService(repo storage.AccountBlockRepository, cpsService service.CPSActionService, logger utils.Logger) service.AccountBlockService {
	return &accountBlockService{repo: repo, cpsService: cpsService, logger: logger}
}

func (s *accountBlockService) GetBranchById(ctx context.Context, id string) (*model.AccountBlock, error) {
	return s.repo.GetBranchById(ctx, id)
}

func (s *accountBlockService) GetRegionById(ctx context.Context, id string) (*model.AccountBlock, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetRegionById", "BlockAccount", "GetRegionById")
	defer span.End()

	region, err := s.repo.GetRegionById(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			span.AddEvent("Region not found", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorRegionNotFound.Code)
		}
		span.AddEvent("Failed to fetch region", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
		return nil, err
	}
	return region, nil
}

func (s *accountBlockService) GetDistrictById(ctx context.Context, districtCode string) (*model.AccountBlock, error) {
	return s.repo.GetDistrictById(ctx, districtCode)
}

func (s *accountBlockService) GetCityById(ctx context.Context, cityCode string) (*model.AccountBlock, error) {
	return s.repo.GetCityById(ctx, cityCode)
}

func (s *accountBlockService) GetAllBranches(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	return s.repo.FindAllBranchesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllRegions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllRegions", "BlockAccount", "GetAllRegions")
	defer span.End()

	regions, err := s.repo.FindAllRegionsWithPagination(ctx, *filterParams)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			span.AddEvent("Region not found", trace.WithAttributes(attribute.Int("page", filterParams.Page), attribute.Int("per_page", filterParams.PerPage), attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorRegionNotFound.Code)
		}
		span.AddEvent("Failed to fetch regions", trace.WithAttributes(attribute.Int("page", filterParams.Page), attribute.Int("per_page", filterParams.PerPage), attribute.String("error", err.Error())))
		return nil, err
	}
	return regions, nil
}

func (s *accountBlockService) GetAllDistricts(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	return s.repo.FindAllDistrictsWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) GetAllCities(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error) {
	return s.repo.FindAllCitiesWithPagination(ctx, *filterParams)
}

func (s *accountBlockService) EnableOrDisableBranches(ctx context.Context, branchIds []string, reason string, enabled bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableBranches", "BlockAccount", "EnableOrDisableBranches")
	defer span.End()

	s.logger.Infof("[EnableOrDisableBranches] processing %d branches, enabled: %v", len(branchIds), enabled)
	var alreadyEnabled []string
	var alreadyDisabled []string
	var previousAction []types.EnableDisableAction
	var currentAction []types.EnableDisableAction

	for _, id := range branchIds {
		branch, err := s.repo.GetBranchById(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorBranchNotFound.Code {
				span.AddEvent("Branch not found", trace.WithAttributes(attribute.String("branch_id", id), attribute.String("error", err.Error())))
				s.logger.Errorf("[EnableOrDisableBranches] branch not found: %s", id)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			span.AddEvent("Failed to get branch", trace.WithAttributes(attribute.String("branch_id", id), attribute.String("error", err.Error())))
			s.logger.Errorf("[EnableOrDisableBranches] failed to get branch %s: %v", id, err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if enabled {
			if branch.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, id)
			}
		} else {
			if !branch.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, id)
			}
		}

		if enabled {
			city, err := s.repo.GetCityById(ctx, branch.CityID.Hex())
			if err != nil {
				span.AddEvent("Failed to fetch city for branch", trace.WithAttributes(attribute.String("branch_id", id), attribute.String("city_id", branch.CityID.Hex()), attribute.String("error", err.Error())))
				return err
			}
			if !city.IsEnabled {
				return errors.New(localization.ErrorCannotEnableBranch.Code)
			}
		}

		previousAction = append(previousAction, types.EnableDisableAction{
			ID:      branch.ID.Hex(),
			Name:    branch.Name,
			Enabled: branch.IsEnabled,
		})
		currentAction = append(currentAction, types.EnableDisableAction{
			ID:      branch.ID.Hex(),
			Name:    branch.Name,
			Enabled: enabled,
		})
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

	cpsAction := core.GenerateCPSAction(ctx, "BRANCH", enabled, previousAction, currentAction, reason, actionType, requestActionType)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableBranches] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableBranches] CPS action created successfully for %d branches", len(branchIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableRegions(ctx context.Context, regionIds []string, reason string, enabled bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableRegions", "BlockAccount", "EnableOrDisableRegions")
	defer span.End()

	s.logger.Infof("[EnableOrDisableRegions] processing %d regions, enabled: %v", len(regionIds), enabled)

	var alreadyEnabled []string
	var alreadyDisabled []string
	var previousAction []types.EnableDisableAction
	var currentAction []types.EnableDisableAction
	for _, id := range regionIds {
		region, err := s.repo.GetRegionById(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorRegionNotFound.Code {
				span.AddEvent("Region not found", trace.WithAttributes(attribute.String("region_id", id), attribute.String("error", err.Error())))
				s.logger.Errorf("[EnableOrDisableRegions] region not found: %s", id)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			span.AddEvent("Failed to get region", trace.WithAttributes(attribute.String("region_id", id), attribute.String("error", err.Error())))
			s.logger.Errorf("[EnableOrDisableRegions] failed to get region %s: %v", id, err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if enabled {
			if region.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, id)
			}
		} else {
			if !region.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, id)
			}
		}

		previousAction = append(previousAction, types.EnableDisableAction{
			ID:      region.ID.Hex(),
			Name:    region.Name,
			Enabled: region.IsEnabled,
		})
		currentAction = append(currentAction, types.EnableDisableAction{
			ID:      region.ID.Hex(),
			Name:    region.Name,
			Enabled: enabled,
		})
	}

	if len(alreadyEnabled) > 0 {
		return fmt.Errorf("these regions are already enabled: %s", alreadyEnabled)
	} else if len(alreadyDisabled) > 0 {
		return fmt.Errorf("these regions are already disabled: %s", alreadyDisabled)
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

	cpsAction := core.GenerateCPSAction(ctx, "REGION", enabled, previousAction, currentAction, reason, actionType, requestActionType)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableRegions] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableRegions] CPS action created successfully for %d regions", len(regionIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableDistricts(ctx context.Context, districtIds []string, reason string, enabled bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableDistricts", "BlockAccount", "EnableOrDisableDistricts")
	defer span.End()

	s.logger.Infof("[EnableOrDisableDistricts] processing %d districts, enabled: %v", len(districtIds), enabled)

	var alreadyEnabled []string
	var alreadyDisabled []string
	var previousAction []types.EnableDisableAction
	var currentAction []types.EnableDisableAction

	for _, id := range districtIds {
		district, err := s.repo.GetDistrictById(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorDistrictNotFound.Code {
				span.AddEvent("District not found", trace.WithAttributes(attribute.String("district_id", id), attribute.String("error", err.Error())))
				s.logger.Errorf("[EnableOrDisableDistricts] district not found: %s", id)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			span.AddEvent("Failed to get district", trace.WithAttributes(attribute.String("district_id", id), attribute.String("error", err.Error())))
			s.logger.Errorf("[EnableOrDisableDistricts] failed to get district %s: %v", id, err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if enabled {
			if district.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, id)
			}
		} else {
			if !district.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, id)
			}
		}

		if enabled {
			region, err := s.repo.GetRegionById(ctx, district.RegionID.Hex())
			if err != nil {
				span.AddEvent("Failed to fetch region for district", trace.WithAttributes(attribute.String("district_id", id), attribute.String("region_id", district.RegionID.Hex()), attribute.String("error", err.Error())))
				return err
			}
			if !region.IsEnabled {
				return errors.New(localization.ErrorCannotEnableDistrict.Code)
			}
		}

		previousAction = append(previousAction, types.EnableDisableAction{
			ID:      district.ID.Hex(),
			Name:    district.Name,
			Enabled: district.IsEnabled,
		})
		currentAction = append(currentAction, types.EnableDisableAction{
			ID:      district.ID.Hex(),
			Name:    district.Name,
			Enabled: enabled,
		})
	}

	if len(alreadyEnabled) > 0 {
		return fmt.Errorf("these districts are already enabled: %s", alreadyEnabled)
	} else if len(alreadyDisabled) > 0 {
		return fmt.Errorf("these districts are already disabled: %s", alreadyDisabled)
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

	cpsAction := core.GenerateCPSAction(ctx, "DISTRICT", enabled, previousAction, currentAction, reason, actionType, requestActionType)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableDistricts] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableDistricts] CPS action created successfully for %d districts", len(districtIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableCities(ctx context.Context, ids []string, reason string, enabled bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableCities", "BlockAccount", "EnableOrDisableCities")
	defer span.End()

	s.logger.Infof("[EnableOrDisableCities] processing %d cities, enabled: %v", len(ids), enabled)

	var alreadyEnabled []string
	var alreadyDisabled []string
	var previousAction []types.EnableDisableAction
	var currentAction []types.EnableDisableAction
	for _, id := range ids {
		city, err := s.repo.FindCityByID(ctx, id)
		if err != nil {
			if err.Error() == localization.ErrorCityNotFound.Code {
				span.AddEvent("City not found", trace.WithAttributes(attribute.String("city_id", id), attribute.String("error", err.Error())))
				s.logger.Errorf("[EnableOrDisableCities] city not found: %s", id)
				return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
			}
			span.AddEvent("Failed to get city", trace.WithAttributes(attribute.String("city_id", id), attribute.String("error", err.Error())))
			s.logger.Errorf("[EnableOrDisableCities] failed to get city %s: %v", id, err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}

		if enabled {
			if city.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, id)
			}
		} else {
			if !city.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, id)
			}
		}

		if enabled {
			district, err := s.repo.GetDistrictById(ctx, city.DistrictID.Hex())
			if err != nil {
				span.AddEvent("Failed to fetch district for city", trace.WithAttributes(attribute.String("city_id", id), attribute.String("region_id", city.RegionID.Hex()), attribute.String("error", err.Error())))
				return err
			}
			if !district.IsEnabled {
				return errors.New(localization.ErrorCannotEnableDistrict.Code)
			}
		}

		previousAction = append(previousAction, types.EnableDisableAction{
			ID:      city.ID.Hex(),
			Name:    city.Name,
			Enabled: city.IsEnabled,
		})
		currentAction = append(currentAction, types.EnableDisableAction{
			ID:      city.ID.Hex(),
			Name:    city.Name,
			Enabled: enabled,
		})

	}

	if len(alreadyEnabled) > 0 {
		return fmt.Errorf("these cities are already enabled: %s", alreadyEnabled)
	} else if len(alreadyDisabled) > 0 {
		return fmt.Errorf("these cities are already disabled: %s", alreadyDisabled)
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

	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, previousAction, currentAction, reason, actionType, requestActionType)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableCities] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableOrDisableCities] CPS action created successfully for %d cities", len(ids))
	return nil
}

func (s *accountBlockService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[Authorize] authorizing account block action: %s", action.RequestAction)

	actions, err := local_util.JsonUnmarshal[[]types.EnableDisableAction](action.CurrentAction)
	if err != nil {
		s.logger.Errorf("[Authorize] failed to unmarshal enable branches action: %v", err)
		return nil, err
	}

	switch constants.RequestAction(action.RequestAction) {
	case constants.RequestEnableBranches:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableBranch(ctx, act.ID, act.Reason, true)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable branch %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d branches", len(*actions))

	case constants.RequestDisableBranches:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableBranch(ctx, act.ID, act.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable branch %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d branches", len(*actions))

	case constants.RequestEnableCities:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableCity(ctx, act.ID, act.Reason, true)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable city %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d cities", len(*actions))

	case constants.RequestDisableCities:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableCity(ctx, act.ID, act.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable city %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d cities", len(*actions))

	case constants.RequestEnableDistricts:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableDistrict(ctx, act.ID, act.Reason, true)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable district %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d districts", len(*actions))

	case constants.RequestDisableDistricts:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableDistrict(ctx, act.ID, act.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable district %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d districts", len(*actions))

	case constants.RequestEnableRegions:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableRegion(ctx, act.ID, act.Reason, true)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to enable region %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully enabled %d regions", len(*actions))

	case constants.RequestDisableRegions:
		for _, act := range *actions {
			err = s.repo.EnableOrDisableRegion(ctx, act.ID, act.Reason, false)
			if err != nil {
				s.logger.Errorf("[Authorize] failed to disable region %s: %v", act.ID, err)
				return nil, err
			}
		}
		s.logger.Infof("[Authorize] successfully disabled %d regions", len(*actions))

	default:
		s.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	action.ActionStatus = string(constants.ActionApproved)
	s.logger.Infof("[Authorize] account block action authorized successfully: %s", action.RequestAction)

	return action, nil
}
