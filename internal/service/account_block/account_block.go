package accountblock

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service/account_block/core"
	cps_constants "github.com/abelshebela/cbe-super-app-cps-action/internal/service/cps_action"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	account_block_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/account_block"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func (s *accountBlockService) GetBranchById(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetBranchById", "BlockAccount", "GetBranchById")
	defer span.End()

	block, err := s.repo.GetBranchesByIds(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	if len(block) == 0 {
		span.AddEvent("Branch not found", trace.WithAttributes(attribute.String("id", id)))
		return nil, errors.New(localization.ErrorBranchNotFound.Code)
	}

	return block[0], nil
}

func (s *accountBlockService) GetRegionById(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetRegionById", "BlockAccount", "GetRegionById")
	defer span.End()

	block, err := s.repo.GetRegionsByIds(ctx, []string{id})
	if err != nil {
		return nil, err
	}
	if len(block) == 0 {
		span.AddEvent("Region not found", trace.WithAttributes(attribute.String("id", id)))
		return nil, errors.New(localization.ErrorRegionNotFound.Code)
	}

	return block[0], nil
}

func (s *accountBlockService) GetDistrictById(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetDistrictById", "BlockAccount", "GetDistrictById")
	defer span.End()

	ids := strings.Split(id, ",")
	block, err := s.repo.GetDistrictsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	if len(block) == 0 {
		span.AddEvent("District not found", trace.WithAttributes(attribute.String("id", id)))
		return nil, errors.New(localization.ErrorDistrictNotFound.Code)
	}

	return block[0], nil
}

// func (s *accountBlockService) GetCityById(ctx context.Context, id string) (*imodel.AccountBlock, error) {
// 	ctx, span := local_util.TraceLogger(ctx, "service", "GetCityById", "BlockAccount", "GetCityById")
// 	defer span.End()

// 	block, err := s.repo.GetCitiesByIds(ctx, []string{id})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(block) == 0 {
// 		span.AddEvent("City not found", trace.WithAttributes(attribute.String("id", id)))
// 		return nil, errors.New(localization.ErrorCityNotFound.Code)
// 	}

// 	return block[0], nil
// }

func (s *accountBlockService) GetAllBranches(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	var isEnabled *bool

	if v, ok := filterParams.Filters["status_check"]; ok {
		if boolVal, ok := v.(bool); ok {
			isEnabled = &boolVal
		}
	}

	data, err := s.repo.FindAllBranchesWithPagination(ctx, *filterParams)
	if err != nil {
		return nil, err
	}

	filterParams.Search = strings.TrimSpace(filterParams.Search)
	if data.Data != nil && isEnabled != nil {
		if *isEnabled == data.Data[0].IsEnabled && (filterParams.Search == data.Data[0].Code || strings.EqualFold(filterParams.Search, data.Data[0].Name)) {
			if *isEnabled {
				return &types.PaginatedResponse[[]*imodel.AccountBlock]{}, errors.New(localization.ErrorBranchAlreadyEnabled.Code)
			} else {
				return &types.PaginatedResponse[[]*imodel.AccountBlock]{}, errors.New(localization.ErrorBranchAlreadyDisabled.Code)
			}
		}
	}

	return data, nil
}

func (s *accountBlockService) GetAllRegions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllRegions", "BlockAccount", "GetAllRegions")
	defer span.End()

	regions, err := s.repo.FindAllRegionsWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch regions", trace.WithAttributes(attribute.Int("page", filterParams.Page), attribute.Int("per_page", filterParams.PerPage), attribute.String("error", err.Error())))
		return nil, err
	}
	return regions, nil
}

func (s *accountBlockService) GetAllDistricts(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	return s.repo.FindAllDistrictsWithPagination(ctx, *filterParams)
}

// func (s *accountBlockService) GetAllCities(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
// 	return s.repo.FindAllCitiesWithPagination(ctx, *filterParams)
// }

func (s *accountBlockService) EnableOrDisableBranches(ctx context.Context, branchIds []string, reason string, enabled bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableBranches", "BlockAccount", "EnableOrDisableBranches")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)

	for _, id := range branchIds {
		ok := local_util.IsOracleHexID(id)
		if id == "" || !ok {
			return errors.New(localization.ErrorInvalidID.Code)
		}
	}

	log.Infof("[AccBlockSvc][EnableDisableBranches] count: %d enabled: %v", len(branchIds), enabled)
	var alreadyEnabled []string
	var alreadyDisabled []string
	var previousAction []types.EnableDisableAction
	var currentAction []types.EnableDisableAction

	branches, err := s.repo.GetBranchesByIds(ctx, branchIds)
	if err != nil {
		span.AddEvent("Failed to get branches", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccBlockSvc][EnableDisableBranches] get err: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(branches) != len(branchIds) {
		log.Errorf("[AccBlockSvc][EnableDisableBranches] not found")
		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
	}

	var districtIDs []string
	for _, branch := range branches {
		if branch.DistrictID != nil {
			districtIDs = append(districtIDs, *branch.DistrictID)
		}
	}

	// Validate districts if enabling
	if enabled {
		districts, err := s.repo.GetDistrictsByIds(ctx, districtIDs)
		if err != nil {
			return err
		}
		districtMap := make(map[string]bool)
		for _, district := range districts {
			districtMap[district.ID] = district.IsEnabled
		}

		// for _, branch := range branches {
		// 	if branch.DistrictID == nil {
		// 		return errors.New(localization.ErrorCannotEnableBranch.Code)
		// 	}
		// 	if isDistrictEnabled, exists := districtMap[*branch.DistrictID]; !exists || !isDistrictEnabled {
		// 		return errors.New(localization.ErrorCannotEnableBranch.Code)
		// 	}
		// }
	}

	fullname := ctx.Value(constants.ContextKey("full_name")).(string)
	for _, branch := range branches {
		if enabled {
			if branch.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, branch.Name)
			}
		} else {
			if !branch.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, branch.Name)
			}
		}

		previousAction = append(previousAction, types.EnableDisableAction{
			ID:      branch.ID,
			Name:    branch.Name,
			Enabled: branch.IsEnabled,
		})
		currentAction = append(currentAction, types.EnableDisableAction{
			ID:      branch.ID,
			Name:    branch.Name,
			Enabled: enabled,
			Reason: types.Reason{
				Reason:    reason,
				CreatedAt: time.Now(),
				CreatedBy: fullname,
			},
		})
	}

	if len(alreadyEnabled) > 0 {
		return errors.New(localization.ErrorBranchAlreadyEnabled.Code)
	} else if len(alreadyDisabled) > 0 {
		return errors.New(localization.ErrorBranchAlreadyDisabled.Code)
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
		log.Errorf("[AccBlockSvc][EnableDisableBranches] cps action err: %v", err)
		return err
	}
	log.Infof("[AccBlockSvc][EnableDisableBranches] done count: %d", len(branchIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableRegions(ctx context.Context, regionIds []string, reason string, enabled bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableRegions", "BlockAccount", "EnableOrDisableRegions")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)

	log.Infof("[AccBlockSvc][EnableDisableRegions] count: %d enabled: %v", len(regionIds), enabled)

	var alreadyEnabled []string
	var alreadyDisabled []string
	var previousAction []types.EnableDisableAction
	var currentAction []types.EnableDisableAction
	regions, err := s.repo.GetRegionsByIds(ctx, regionIds)
	if err != nil {
		span.AddEvent("Failed to get regions", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccBlockSvc][EnableDisableRegions] get err: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(regions) != len(regionIds) {
		log.Errorf("[AccBlockSvc][EnableDisableRegions] not found")
		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
	}

	fullname := ctx.Value(constants.ContextKey("full_name")).(string)
	for _, region := range regions {
		if enabled {
			if region.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, region.Name)
			}
		} else {
			if !region.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, region.Name)
			}
		}

		previousAction = append(previousAction, types.EnableDisableAction{
			ID:      region.ID,
			Name:    region.Name,
			Enabled: region.IsEnabled,
		})
		currentAction = append(currentAction, types.EnableDisableAction{
			ID:      region.ID,
			Name:    region.Name,
			Enabled: enabled,
			Reason: types.Reason{
				Reason:    reason,
				CreatedAt: time.Now(),
				CreatedBy: fullname,
			},
		})
	}

	if len(alreadyEnabled) > 0 {
		return errors.New(localization.ErrorRegionAlreadyEnabled.Code)
	} else if len(alreadyDisabled) > 0 {
		return errors.New(localization.ErrorRegionAlreadyDisabled.Code)
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
		log.Errorf("[AccBlockSvc][EnableDisableRegions] cps action err: %v", err)
		return err
	}
	log.Infof("[AccBlockSvc][EnableDisableRegions] done count: %d", len(regionIds))
	return nil
}

func (s *accountBlockService) EnableOrDisableDistricts(ctx context.Context, districtIds []string, reason string, enabled bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableDistricts", "BlockAccount", "EnableOrDisableDistricts")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)

	log.Infof("[AccBlockSvc][EnableDisableDistricts] count: %d enabled: %v", len(districtIds), enabled)

	var alreadyEnabled []string
	var alreadyDisabled []string
	var previousAction []types.EnableDisableAction
	var currentAction []types.EnableDisableAction

	districts, err := s.repo.GetDistrictsByIds(ctx, districtIds)
	if err != nil {
		span.AddEvent("Failed to get districts", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccBlockSvc][EnableDisableDistricts] get err: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(districts) != len(districtIds) {
		log.Errorf("[AccBlockSvc][EnableDisableDistricts] not found")
		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
	}

	var regionIDs []string
	for _, district := range districts {
		if district.RegionID != nil {
			regionIDs = append(regionIDs, *district.RegionID)
		}
	}

	// if enabled {
	// 	regions, err := s.repo.GetRegionsByIds(ctx, regionIDs)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	regionMap := make(map[string]bool)
	// 	for _, region := range regions {
	// 		regionMap[region.ID] = region.IsEnabled
	// 	}

	// for _, district := range districts {
	// 	if district.RegionID == nil {
	// 		return errors.New(localization.ErrorCannotEnableDistrict.Code)
	// 	}
	// 	if isRegionEnabled, exists := regionMap[*district.RegionID]; !exists || !isRegionEnabled {
	// 		return errors.New(localization.ErrorCannotEnableDistrict.Code)
	// 	}
	// }
	// }

	fullname := ctx.Value(constants.ContextKey("full_name")).(string)
	for _, district := range districts {
		if enabled {
			if district.IsEnabled {
				alreadyEnabled = append(alreadyEnabled, district.Name)
			}
		} else {
			if !district.IsEnabled {
				alreadyDisabled = append(alreadyDisabled, district.Name)
			}
		}

		previousAction = append(previousAction, types.EnableDisableAction{
			ID:      district.ID,
			Name:    district.Name,
			Enabled: district.IsEnabled,
		})
		currentAction = append(currentAction, types.EnableDisableAction{
			ID:      district.ID,
			Name:    district.Name,
			Enabled: enabled,
			Reason: types.Reason{
				Reason:    reason,
				CreatedAt: time.Now(),
				CreatedBy: fullname,
			},
		})
	}

	if len(alreadyEnabled) > 0 {
		return errors.New(localization.ErrorDistrictAlreadyEnabled.Code)
	} else if len(alreadyDisabled) > 0 {
		return errors.New(localization.ErrorDistrictAlreadyDisabled.Code)
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
		log.Errorf("[AccBlockSvc][EnableDisableDistricts] cps action err: %v", err)
		return err
	}
	log.Infof("[AccBlockSvc][EnableDisableDistricts] done count: %d", len(districtIds))
	return nil
}

// func (s *accountBlockService) EnableOrDisableCities(ctx context.Context, ids []string, reason string, enabled bool) error {

// 	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableCities", "BlockAccount", "EnableOrDisableCities")
// 	defer span.End()

// 	log.Infof("[AccBlockSvc][EnableDisableCities] count: %d enabled: %v", len(ids), enabled)

// 	var alreadyEnabled []string
// 	var alreadyDisabled []string
// 	var previousAction []types.EnableDisableAction
// 	var currentAction []types.EnableDisableAction
// 	cities, err := s.repo.GetCitiesByIds(ctx, ids)
// 	if err != nil {
// 		span.AddEvent("Failed to get cities", trace.WithAttributes(attribute.String("error", err.Error())))
// 		log.Errorf("[AccBlockSvc][EnableDisableCities] get err: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	if len(cities) != len(ids) {
// 		log.Errorf("[AccBlockSvc][EnableDisableCities] not found")
// 		return errors.New(localization.ErrorOneOrMoreInvalidCodes.Code)
// 	}

// 	var districtIDs []string
// 	for _, city := range cities {
// 		if city.DistrictID != nil {
// 			districtIDs = append(districtIDs, *city.DistrictID)
// 		}
// 	}

// 	if enabled {
// 		districts, err := s.repo.GetDistrictsByIds(ctx, districtIDs)
// 		if err != nil {
// 			return err
// 		}
// 		districtMap := make(map[string]bool)
// 		for _, district := range districts {
// 			districtMap[district.ID] = district.IsEnabled
// 		}

// 		for _, city := range cities {
// 			if city.DistrictID == nil {
// 				return errors.New(localization.ErrorCannotEnablCity.Code)
// 			}
// 			if isDistrictEnabled, exists := districtMap[*city.DistrictID]; !exists || !isDistrictEnabled {
// 				return errors.New(localization.ErrorCannotEnablCity.Code)
// 			}
// 		}
// 	}

// 	fullname := ctx.Value(constants.ContextKey("full_name")).(string)
// 	for _, city := range cities {
// 		if enabled {
// 			if city.IsEnabled {
// 				alreadyEnabled = append(alreadyEnabled, city.Name)
// 			}
// 		} else {
// 			if !city.IsEnabled {
// 				alreadyDisabled = append(alreadyDisabled, city.Name)
// 			}
// 		}

// 		previousAction = append(previousAction, types.EnableDisableAction{
// 			ID:      city.ID,
// 			Name:    city.Name,
// 			Enabled: city.IsEnabled,
// 		})
// 		currentAction = append(currentAction, types.EnableDisableAction{
// 			ID:      city.ID,
// 			Name:    city.Name,
// 			Enabled: enabled,
// 			Reason: types.Reason{
// 				Reason:    reason,
// 				CreatedAt: time.Now(),
// 				CreatedBy: fullname,
// 			},
// 		})
// 	}

// 	if len(alreadyEnabled) > 0 {
// 		return errors.New(localization.ErrorCityAlreadyEnabled.Code)
// 	} else if len(alreadyDisabled) > 0 {
// 		return errors.New(localization.ErrorCityAlreadyDisabled.Code)
// 	}

// 	var actionType constants.ActionType
// 	var requestActionType constants.RequestAction
// 	if enabled {
// 		actionType = constants.ActionType(cps_constants.ActionEnable)
// 		requestActionType = constants.RequestEnableCities
// 	} else {
// 		actionType = constants.ActionType(cps_constants.ActionDisable)
// 		requestActionType = constants.RequestDisableCities
// 	}

// 	cpsAction := core.GenerateCPSAction(ctx, "CITY", enabled, previousAction, currentAction, reason, actionType, requestActionType)

// 	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
// 		span.AddEvent("Failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
// 		log.Errorf("[AccBlockSvc][EnableDisableCities] cps action err: %v", err)
// 		return err
// 	}
// 	log.Infof("[AccBlockSvc][EnableDisableCities] done count: %d", len(ids))
// 	return nil
// }

func (s *accountBlockService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[AccBlockSvc][Authorize] action: %s", action.RequestAction)

	actions, err := local_util.JsonUnmarshal[[]types.EnableDisableAction](action.CurrentAction)
	if err != nil {
		log.Errorf("[AccBlockSvc][Authorize] unmarshal err: %v", err)
		return nil, err
	}

	switch constants.RequestAction(action.RequestAction) {
	case constants.RequestEnableBranches:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableBranches(ctx, ids, &(*actions)[0].Reason, true); err != nil {
			log.Errorf("[AccBlockSvc][Authorize] enable branches err: %v", err)
			return nil, err
		}
		log.Infof("[AccBlockSvc][Authorize] enabled %d branches", len(*actions))

	case constants.RequestDisableBranches:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableBranches(ctx, ids, &(*actions)[0].Reason, false); err != nil {
			log.Errorf("[AccBlockSvc][Authorize] disable branches err: %v", err)
			return nil, err
		}
		log.Infof("[AccBlockSvc][Authorize] disabled %d branches", len(*actions))

	// case constants.RequestEnableCities:
	// 	var ids []string
	// 	for _, act := range *actions {
	// 		ids = append(ids, act.ID)
	// 	}
	// 	if err := s.repo.EnableOrDisableCities(ctx, ids, &(*actions)[0].Reason, true); err != nil {
	// 		log.Errorf("[AccBlockSvc][Authorize] enable cities err: %v", err)
	// 		return nil, err
	// 	}
	// 	log.Infof("[AccBlockSvc][Authorize] enabled %d cities", len(*actions))

	// case constants.RequestDisableCities:
	// 	var ids []string
	// 	for _, act := range *actions {
	// 		ids = append(ids, act.ID)
	// 	}
	// 	if err := s.repo.EnableOrDisableCities(ctx, ids, &(*actions)[0].Reason, false); err != nil {
	// 		log.Errorf("[AccBlockSvc][Authorize] disable cities err: %v", err)
	// 		return nil, err
	// 	}
	// 	log.Infof("[AccBlockSvc][Authorize] disabled %d cities", len(*actions))

	case constants.RequestEnableDistricts:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableDistricts(ctx, ids, &(*actions)[0].Reason, true); err != nil {
			log.Errorf("[AccBlockSvc][Authorize] enable districts err: %v", err)
			return nil, err
		}
		log.Infof("[AccBlockSvc][Authorize] enabled %d districts", len(*actions))

	case constants.RequestDisableDistricts:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableDistricts(ctx, ids, &(*actions)[0].Reason, false); err != nil {
			log.Errorf("[AccBlockSvc][Authorize] disable districts err: %v", err)
			return nil, err
		}
		log.Infof("[AccBlockSvc][Authorize] disabled %d districts", len(*actions))

	case constants.RequestEnableRegions:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableRegions(ctx, ids, &(*actions)[0].Reason, true); err != nil {
			log.Errorf("[AccBlockSvc][Authorize] enable regions err: %v", err)
			return nil, err
		}
		log.Infof("[AccBlockSvc][Authorize] enabled %d regions", len(*actions))

	case constants.RequestDisableRegions:
		var ids []string
		for _, act := range *actions {
			ids = append(ids, act.ID)
		}
		if err := s.repo.EnableOrDisableRegions(ctx, ids, &(*actions)[0].Reason, false); err != nil {
			log.Errorf("[AccBlockSvc][Authorize] disable regions err: %v", err)
			return nil, err
		}
		log.Infof("[AccBlockSvc][Authorize] disabled %d regions", len(*actions))

	default:
		log.Errorf("[AccBlockSvc][Authorize] unsupported: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	action.ActionStatus = string(constants.ActionApproved)
	log.Infof("[AccBlockSvc][Authorize] done: %s", action.RequestAction)

	return action, nil
}

func (s *accountBlockService) GetAccountBlockDetails(ctx context.Context, id string, filter *types.Filter) (*types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAccountBlockDetails", "BlockAccount", "GetAccountBlockDetails")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)

	result, err := s.repo.GetAccountBlockDetails(ctx, id, *filter)
	if err != nil {
		if err.Error() == localization.ErrorBranchNotFound.Code || err.Error() == localization.ErrorCityNotFound.Code || err.Error() == localization.ErrorDistrictNotFound.Code || err.Error() == localization.ErrorRegionNotFound.Code || err.Error() == localization.ErrorActionNotFound.Code {
			return &types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse]{}, nil
		}
		log.Errorf("[AccBlockSvc][GetDetails] id=%s err: %v", id, err)
		return nil, err
	}

	return result, nil
}

func (s *accountBlockService) GetPreviousReasons(ctx context.Context, entityType string, identifier string) (*account_block_dto.PreviousDisableReasonsResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPreviousReasons", "BlockAccount", "GetPreviousReasons")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, s.logger)

	reasons, err := s.repo.GetPreviousReasons(ctx, entityType, identifier)
	if err != nil {
		log.Errorf("[AccBlockSvc][GetPreviousReasons] type=%s identifier=%s err: %v", entityType, identifier, err)
		return nil, err
	}

	return &account_block_dto.PreviousDisableReasonsResponse{
		EntityType:    entityType,
		Identifier:    identifier,
		DisableReason: reasons,
	}, nil
}
