package bps_action_role_service

import (
	"cbe-super-app-cps-action/internal/constants"
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type bpsActionRoleService struct {
	repo       storage.BPSActionRoleRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewBPSActionRoleService(repo storage.BPSActionRoleRepository, cps service.CPSActionService, logger utils.Logger) service.BPSActionRoleService {
	return &bpsActionRoleService{repo: repo, cpsService: cps, logger: logger}
}

// FindAllWithPagination implements service.bpsActionRoleService.
func (s *bpsActionRoleService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.ActionRole], error) {
	return s.repo.FindAllWithPagination(ctx, filter)
}

// GetByActionCode implements service.bpsActionRoleService.
func (s *bpsActionRoleService) GetByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {
	if actionCode == "" {
		s.logger.Errorf("[GetByActionCode] action code is empty")
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	res, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code == localization.ErrorResourceNotFound.Code {
			s.logger.Errorf("[GetByActionCode] action role not found: %s", actionCode)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		s.logger.Errorf("[GetByActionCode] failed to fetch action role: %v", err)
		return nil, err
	}
	s.logger.Infof("[GetByActionCode] action role retrieved successfully for action code: %s", actionCode)
	return res, nil
}

// Create implements service.bpsActionRoleService.
func (s *bpsActionRoleService) Create(ctx context.Context, req struct {
	ActionCode       string
	ActionName       string
	AssignedMakers   []string
	AssignedCheckers [][]string
}) error {
	s.logger.Infof("[Create] creating action role for action code: %s", req.ActionCode)
	if req.ActionCode == "" || req.ActionName == "" {
		s.logger.Errorf("[Create] invalid input parameters")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("[Create] incomplete user data")
		return errors.New(localization.ErrorUserUnauthorized.Code)
	}

	// build CPS action payload
	payload := model.ActionRoleCPSAction{
		ActionCode:       req.ActionCode,
		ActionName:       req.ActionName,
		AssignedMakers:   req.AssignedMakers,
		AssignedCheckers: req.AssignedCheckers,
		Enabled:          true,
	}
	cpsAction := lib.CpsModelBuilder(
		req.ActionCode,
		maker,
		nil,
		payload,
		string(constants.RequestCreateActionRole),
		constants.CREATE,
	)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[Create] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Create] action role creation request created successfully for action code: %s", req.ActionCode)
	return nil
}

// Update implements service.bpsActionRoleService.
func (s *bpsActionRoleService) Update(ctx context.Context, actionCode string, req actionrole_dto.UpdateActionRoleRequest) error {
	s.logger.Infof("[Update] updating action role for action code: %s", actionCode)
	if actionCode == "" {
		s.logger.Errorf("[Update] action code is empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		s.logger.Errorf("[Update] failed to find action role: %v", err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("[Update] incomplete user data")
		return errors.New(localization.ErrorUserUnauthorized.Code)
	}

	payload := model.ActionRoleCPSAction{
		ActionCode: actionCode,
		ActionName: local_util.NonEmptyString(req.ActionName, old.ActionName),
		Enabled:    old.Enabled,
	}

	if req.AssignedCheckers != nil {
		payload.AssignedMakers = req.AssignedMakers
	}
	if req.AssignedCheckers != nil {
		payload.AssignedCheckers = req.AssignedCheckers
	}

	cpsAction := lib.CpsModelBuilder(
		actionCode,
		maker,
		old,
		payload,
		string(constants.RequestUpdateActionRole),
		constants.UPDATE,
	)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[Update] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Update] action role update request created successfully for action code: %s", actionCode)
	return nil
}

func (s *bpsActionRoleService) Enable(ctx context.Context, actionCode string) error {
	s.logger.Infof("[Enable] enabling action role for action code: %s", actionCode)
	if actionCode == "" {
		s.logger.Errorf("[Enable] action code is empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	// Ensure exists
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		s.logger.Errorf("[Enable] failed to find action role: %v", err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if old.Enabled {
		s.logger.Errorf("[Enable] action role already enabled")
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRoleCPSAction{ActionCode: actionCode, Enabled: true}
	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestEnableActionRole), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cps); err != nil {
		s.logger.Errorf("[Enable] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Enable] action role enable request created successfully for action code: %s", actionCode)
	return nil
}

func (s *bpsActionRoleService) Disable(ctx context.Context, actionCode string) error {
	s.logger.Infof("[Disable] disabling action role for action code: %s", actionCode)
	if actionCode == "" {
		s.logger.Errorf("[Disable] action code is empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		s.logger.Errorf("[Disable] failed to find action role: %v", err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if !old.Enabled {
		s.logger.Errorf("[Disable] action role already disabled")
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRoleCPSAction{ActionCode: actionCode, Enabled: false}
	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestDisableActionRole), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cps); err != nil {
		s.logger.Errorf("[Disable] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Disable] action role disable request created successfully for action code: %s", actionCode)
	return nil
}

// Authorize applies approved CPS actions
func (s *bpsActionRoleService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[Authorize] authorizing action role action: %s", action.RequestAction)
	cur, err := local_util.JsonUnmarshal[model.ActionRoleCPSAction](action.CurrentAction)
	if err != nil {
		s.logger.Errorf("[Authorize] failed to unmarshal action: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
	}
	switch action.ActionType {
	case string(constants.CREATE):
		ar, err := s.bindActionRoleModel(*cur)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to bind action role model: %v", err)
			return nil, err
		}
		ar.ID = bson.NewObjectID()
		ar.CreatedAt = time.Now()
		ar.UpdatedAt = time.Now()
		if err := s.repo.Create(ctx, &ar); err != nil {
			s.logger.Errorf("[Authorize] failed to create action role: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] action role created successfully")
		return action, nil
	case string(constants.UPDATE):
		existing, err := s.repo.FindByActionCode(ctx, action.UniqueId)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to find existing action role: %v", err)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		upd, err := s.bindActionRoleModel(*cur)
		if err != nil {
			s.logger.Errorf("[Authorize] failed to bind action role model: %v", err)
			return nil, err
		}
		upd.UpdatedAt = time.Now()
		upd.CreatedAt = existing.CreatedAt
		if cur.AssignedMakers == nil && cur.AssignedCheckers == nil && cur.ActionName == "" {
			// enable/disable path
			if err := s.repo.EnableOrDisableByActionCode(ctx, existing.ActionCode, cur.Enabled); err != nil {
				s.logger.Errorf("[Authorize] failed to enable/disable action role: %v", err)
				return nil, err
			}
			s.logger.Infof("[Authorize] action role enable/disable completed successfully")
			return action, nil
		}
		if err := s.repo.UpdateByActionCode(ctx, existing.ActionCode, &upd); err != nil {
			s.logger.Errorf("[Authorize] failed to update action role: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] action role updated successfully")
		return action, nil
	default:
		s.logger.Errorf("[Authorize] unsupported action type: %s", action.ActionType)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *bpsActionRoleService) bindActionRoleModel(in model.ActionRoleCPSAction) (model.ActionRole, error) {
	var makers []bson.ObjectID
	for _, id := range in.AssignedMakers {
		obj, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return model.ActionRole{}, errors.New(localization.ErrorInvalidID.Code)
		}
		makers = append(makers, obj)
	}
	var checkers [][]bson.ObjectID
	for _, grp := range in.AssignedCheckers {
		var g []bson.ObjectID
		for _, id := range grp {
			obj, err := bson.ObjectIDFromHex(id)
			if err != nil {
				return model.ActionRole{}, errors.New(localization.ErrorInvalidID.Code)
			}
			g = append(g, obj)
		}
		checkers = append(checkers, g)
	}
	return model.ActionRole{
		ActionCode:      in.ActionCode,
		ActionName:      in.ActionName,
		AssignedMakers:  makers,
		AssignedChecker: checkers,
		Enabled:         in.Enabled,
	}, nil
}
