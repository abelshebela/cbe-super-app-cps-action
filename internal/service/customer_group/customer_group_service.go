package customergroup

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	customer_group_dto "cbe-super-app-cps-action/internal/constants/dto/customer_group"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/customer_group/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	shared_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerGroupService struct {
	repo       storage.CustomerGroupRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewCustomerGroupService(repo storage.CustomerGroupRepository, cpsService service.CPSActionService, logger utils.Logger) *customerGroupService {
	return &customerGroupService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *customerGroupService) Create(ctx context.Context, req customer_group_dto.CreateSegmentRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	newSeg := imodel.Segment{
		CustomerGroup:           req.CustomerGroup,
		CustomerGroupLabel:      req.CustomerGroupLabel,
		CustomerSegment:         req.CustomerSegment,
		CustomerSegmentLabel:    req.CustomerSegmentLabel,
		CustomerSubsegment:      req.CustomerSubsegment,
		CustomerSubsegmentLabel: req.CustomerSubsegmentLabel,
		SuperappRole:            req.SuperappRole,
		SuperappRoleLabel:       req.SuperappRoleLabel,
		IsEnabled:               true,
		CreatedAt:               time.Now(),
		LastModifiedAt:          time.Now(),
	}

	cpsActionData := lib.CpsModelBuilder("", makerUser, nil, core.MapSegmentToMap(newSeg), string(constants.RequestCreateCustomerGroup), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustomerGroup][Create] cps action err: %v", err)
		return err
	}
	log.Infof("[CustomerGroup][Create] request created")
	return nil
}

func (s *customerGroupService) Update(ctx context.Context, id string, req customer_group_dto.UpdateSegmentRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustomerGroup][Update] find err: %v", err)
		return err
	}

	updated := *existing
	updated.CustomerGroup = req.CustomerGroup
	updated.CustomerGroupLabel = req.CustomerGroupLabel
	updated.CustomerSegment = req.CustomerSegment
	updated.CustomerSegmentLabel = req.CustomerSegmentLabel
	updated.CustomerSubsegment = req.CustomerSubsegment
	updated.CustomerSubsegmentLabel = req.CustomerSubsegmentLabel
	updated.SuperappRole = req.SuperappRole
	updated.SuperappRoleLabel = req.SuperappRoleLabel
	updated.LastModifiedAt = time.Now()

	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapSegmentToMap(*existing), core.MapSegmentToMap(updated), string(constants.RequestUpdateCustomerGroup), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustomerGroup][Update] cps action err: %v", err)
		return err
	}
	log.Infof("[CustomerGroup][Update] request created")
	return nil
}

func (s *customerGroupService) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustomerGroup][Delete] find err: %v", err)
		return err
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapSegmentToMap(*existing), core.MapSegmentToMap(*existing), string(constants.RequestDeleteCustomerGroup), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustomerGroup][Delete] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *customerGroupService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[CustomerGroup][EnableOrDisable] find err: %v", err)
		return err
	}

	if existing.IsEnabled && enable {
		return errors.New("segment already enabled")
	}
	if !existing.IsEnabled && !enable {
		return errors.New("segment already disabled")
	}

	updated := *existing
	updated.IsEnabled = enable
	updated.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableCustomerGroup
	} else {
		action = constants.RequestDisableCustomerGroup
	}

	cpsActionData := lib.CpsModelBuilder(id, makerUser, core.MapSegmentToMap(*existing), core.MapSegmentToMap(updated), string(action), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[CustomerGroup][EnableOrDisable] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *customerGroupService) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.Segment], error) {
	return s.repo.FindAllWithPagination(ctx, *filterParam)
}

func (s *customerGroupService) FindByID(ctx context.Context, id string) (*imodel.Segment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *customerGroupService) Authorize(ctx context.Context, action *shared_model.CPSAction) (*shared_model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[CustomerGroup][Authorize] action: %s", action.RequestAction)

	seg, err := local_util.JsonUnmarshal[imodel.Segment](action.CurrentAction)
	if err != nil || seg == nil {
		log.Errorf("[CustomerGroup][Authorize] unmarshal err: %v", err)
		return nil, err
	}

	switch action.RequestAction {
	case string(constants.RequestCreateCustomerGroup):
		if err := s.repo.Create(ctx, seg); err != nil {
			log.Errorf("[CustomerGroup][Authorize] create err: %v", err)
			return nil, err
		}
	case string(constants.RequestUpdateCustomerGroup):
		if err := s.repo.Update(ctx, action.UniqueId, seg); err != nil {
			log.Errorf("[CustomerGroup][Authorize] update err: %v", err)
			return nil, err
		}
	case string(constants.RequestDeleteCustomerGroup):
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			log.Errorf("[CustomerGroup][Authorize] delete err: %v", err)
			return nil, err
		}
	case string(constants.RequestEnableCustomerGroup):
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, true); err != nil {
			log.Errorf("[CustomerGroup][Authorize] enable err: %v", err)
			return nil, err
		}
	case string(constants.RequestDisableCustomerGroup):
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, false); err != nil {
			log.Errorf("[CustomerGroup][Authorize] disable err: %v", err)
			return nil, err
		}
	default:
		log.Errorf("[CustomerGroup][Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	return action, nil
}
