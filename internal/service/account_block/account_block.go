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

	account_block_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "GetBranchById", "BlockAccount", "GetBranchById")
	defer span.End()

	block, err := s.repo.GetBranchesByIds(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	if len(block) == 0 {
		span.AddEvent("Branch not found", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
		return nil, errors.New(localization.ErrorBranchNotFound.Code)
	}

	return block[0], nil
}

func (s *accountBlockService) GetRegionById(ctx context.Context, id string) (*model.AccountBlock, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetRegionById", "BlockAccount", "GetRegionById")
	defer span.End()

	block, err := s.repo.GetRegionsByIds(ctx, []string{id})
	if len(block) == 0 {
		span.AddEvent("Region not found", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
		return nil, errors.New(localization.ErrorRegionNotFound.Code)
	}

	return block[0], nil
}

func (s *accountBlockService) GetDistrictById(ctx context.Context, id string) (*model.AccountBlock, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetDistrictById", "BlockAccount", "GetDistrictById")
	defer span.End()

	block, err := s.repo.GetDistrictsByIds(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	if len(block) == 0 {
		span.AddEvent("District not found", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
		return nil, errors.New(localization.ErrorDistrictNotFound.Code)
	}

	return block[0], nil
}

func (s *accountBlockService) GetCityById(ctx context.Context, id string) (*model.AccountBlock, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCityById", "BlockAccount", "GetCityById")
	defer span.End()

	block, err := s.repo.GetCitiesByIds(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	if len(block) == 0 {
		span.AddEvent("City not found", trace.WithAttributes(attribute.String("id", id), attribute.String("error", err.Error())))
		return nil, errors.New(localization.ErrorCityNotFound.Code)
	}

	return block[0], nil
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

	branches, err := s.repo.GetBranchesByIds(ctx, branchIds)
	if err != nil {
		span.AddEvent("Failed to get branches", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableBranches] failed to get branches: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(branches) != len(branchIds) {
		s.logger.Errorf("[EnableOrDisableBranches] one or more branches not found")
		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
	}

	var cityIDs []string
	for _, branch := range branches {
		cityIDs = append(cityIDs, branch.CityID.Hex())
	}

	// Validate cities if enabling
	if enabled {
		cities, err := s.repo.GetCitiesByIds(ctx, cityIDs)
		if err != nil {
			return err
		}
		cityMap := make(map[string]bool)
		for _, city := range cities {
			cityMap[city.ID.Hex()] = city.IsEnabled
		}

		for _, branch := range branches {
			if isCityEnabled, exists := cityMap[branch.CityID.Hex()]; !exists || !isCityEnabled {
				return errors.New(localization.ErrorCannotEnableBranch.Code)
			}
		}
	}

	for _, branch := range branches {
		if enabled {
			if branch.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, branch.ID.Hex())
			}
		} else {
			if !branch.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, branch.ID.Hex())
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
			Reason:  reason,
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
	regions, err := s.repo.GetRegionsByIds(ctx, regionIds)
	if err != nil {
		span.AddEvent("Failed to get regions", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableRegions] failed to get regions: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(regions) != len(regionIds) {
		s.logger.Errorf("[EnableOrDisableRegions] one or more regions not found")
		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code) // Corrected error code usage
	}

	for _, region := range regions {
		if enabled {
			if region.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, region.ID.Hex())
			}
		} else {
			if !region.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, region.ID.Hex())
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
			Reason:  reason,
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

	districts, err := s.repo.GetDistrictsByIds(ctx, districtIds)
	if err != nil {
		span.AddEvent("Failed to get districts", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableDistricts] failed to get districts: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(districts) != len(districtIds) {
		s.logger.Errorf("[EnableOrDisableDistricts] one or more districts not found")
		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
	}

	var regionIDs []string
	for _, district := range districts {
		regionIDs = append(regionIDs, district.RegionID.Hex())
	}

	if enabled {
		regions, err := s.repo.GetRegionsByIds(ctx, regionIDs)
		if err != nil {
			return err
		}
		regionMap := make(map[string]bool)
		for _, region := range regions {
			regionMap[region.ID.Hex()] = region.IsEnabled
		}

		for _, district := range districts {
			if isRegionEnabled, exists := regionMap[district.RegionID.Hex()]; !exists || !isRegionEnabled {
				return errors.New(localization.ErrorCannotEnableDistrict.Code)
			}
		}
	}

	for _, district := range districts {
		if enabled {
			if district.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, district.ID.Hex())
			}
		} else {
			if !district.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, district.ID.Hex())
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
			Reason:  reason,
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
	cities, err := s.repo.GetCitiesByIds(ctx, ids)
	if err != nil {
		span.AddEvent("Failed to get cities", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[EnableOrDisableCities] failed to get cities: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(cities) != len(ids) {
		s.logger.Errorf("[EnableOrDisableCities] one or more cities not found")
		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
	}

	var districtIDs []string
	for _, city := range cities {
		districtIDs = append(districtIDs, city.DistrictID.Hex())
	}

	if enabled {
		districts, err := s.repo.GetDistrictsByIds(ctx, districtIDs)
		if err != nil {
			return err
		}
		districtMap := make(map[string]bool)
		for _, district := range districts {
			districtMap[district.ID.Hex()] = district.IsEnabled
		}

		for _, city := range cities {
			if isDistrictEnabled, exists := districtMap[city.DistrictID.Hex()]; !exists || !isDistrictEnabled {
				return errors.New(localization.ErrorCannotEnableDistrict.Code)
			}
		}
	}

	for _, city := range cities {
		if enabled {
			if city.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, city.ID.Hex())
			}
		} else {
			if !city.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, city.ID.Hex())
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
			Reason:  reason,
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
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableBranches(ctx, ids, (*actions)[0].Reason, true); err != nil {
			s.logger.Errorf("[Authorize] failed to enable branches: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully enabled %d branches", len(*actions))

	case constants.RequestDisableBranches:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableBranches(ctx, ids, (*actions)[0].Reason, false); err != nil {
			s.logger.Errorf("[Authorize] failed to disable branches: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully disabled %d branches", len(*actions))

	case constants.RequestEnableCities:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableCities(ctx, ids, (*actions)[0].Reason, true); err != nil {
			s.logger.Errorf("[Authorize] failed to enable cities: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully enabled %d cities", len(*actions))

	case constants.RequestDisableCities:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableCities(ctx, ids, (*actions)[0].Reason, false); err != nil {
			s.logger.Errorf("[Authorize] failed to disable cities: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully disabled %d cities", len(*actions))

	case constants.RequestEnableDistricts:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableDistricts(ctx, ids, (*actions)[0].Reason, true); err != nil {
			s.logger.Errorf("[Authorize] failed to enable districts: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully enabled %d districts", len(*actions))

	case constants.RequestDisableDistricts:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableDistricts(ctx, ids, (*actions)[0].Reason, false); err != nil {
			s.logger.Errorf("[Authorize] failed to disable districts: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully disabled %d districts", len(*actions))

	case constants.RequestEnableRegions:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableRegions(ctx, ids, (*actions)[0].Reason, true); err != nil {
			s.logger.Errorf("[Authorize] failed to enable regions: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully enabled %d regions", len(*actions))

	case constants.RequestDisableRegions:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableRegions(ctx, ids, (*actions)[0].Reason, false); err != nil {
			s.logger.Errorf("[Authorize] failed to disable regions: %v", err)
			return nil, err
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

func (s *accountBlockService) GetAccountBlockDetails(ctx context.Context, id string, filter *types.Filter) (*types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAccountBlockDetails", "BlockAccount", "GetAccountBlockDetails")
	defer span.End()

	result, err := s.repo.GetAccountBlockDetails(ctx, id, *filter)
	if err != nil {
		if err.Error() == localization.ErrorBranchNotFound.Code || err.Error() == localization.ErrorCityNotFound.Code || err.Error() == localization.ErrorDistrictNotFound.Code || err.Error() == localization.ErrorRegionNotFound.Code || err.Error() == localization.ErrorActionNotFound.Code {
			return &types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse]{}, nil
		}
		s.logger.Errorf("[GetAccountBlockDetails] failed to get account block details for id %s: %v", id, err)
		return nil, err
	}

	return result, nil
}
