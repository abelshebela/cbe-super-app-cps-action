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
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type cpsActionRoleService struct {
	repo       storage.CPSActionRoleRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewCPSActionRoleService(repo storage.CPSActionRoleRepository, cps service.CPSActionService, logger utils.Logger) service.CPSActionRoleService {
	return &cpsActionRoleService{repo: repo, cpsService: cps, logger: logger}
}

// FindAllWithPagination implements service.bpsActionRoleService.
func (s *cpsActionRoleService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.CPSActionRole], error) {
	return s.repo.FindAllWithPagination(ctx, filter)
}

// GetByActionCode implements service.bpsActionRoleService.
func (s *cpsActionRoleService) GetByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {
	if actionCode == "" {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	res, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code == localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		return nil, err
	}
	return res, nil
}

// Create implements service.CPSActionRoleService.
func (s *cpsActionRoleService) Create(ctx context.Context, req actionrole_dto.CreateActionRoleRequest) error {
	if req.ActionName == "" {
		return errors.New(localization.ErrorActionNameIsRequired.Code)
	}
	req.ActionName = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(req.ActionName), " ", "_"))
	req.ActionCode = req.ActionName
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		return errors.New(localization.ErrorUserUnauthorized.Code)
	}

	var makers []bson.ObjectID
	for _, id := range req.AssignedMakersRoles {
		obj, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorInvalidID.Code)
		}
		makers = append(makers, obj)
	}
	var checkers [][]bson.ObjectID
	for _, grp := range req.AssignedCheckerRoles {
		var g []bson.ObjectID
		for _, id := range grp {
			obj, err := bson.ObjectIDFromHex(id)
			if err != nil {
				return errors.New(localization.ErrorInvalidID.Code)
			}
			g = append(g, obj)
		}
		checkers = append(checkers, g)
	}

	// build CPS action payload
	payload := model.ActionRole{
		ActionCode:           req.ActionCode,
		ActionName:           req.ActionName,
		AssignedMakersRoles:  makers,
		AssignedCheckerRoles: checkers,
		Enabled:              true,
	}
	cpsAction := lib.CpsModelBuilder(
		req.ActionCode,
		maker,
		nil,
		payload,
		string(constants.RequestCreateActionRole),
		constants.CREATE,
	)
	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

// Update implements service.CPSActionRoleService.
func (s *cpsActionRoleService) Update(ctx context.Context, actionCode string, req actionrole_dto.UpdateActionRoleRequest) error {
	if actionCode == "" {
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		return errors.New(localization.ErrorUserUnauthorized.Code)
	}

	payload := model.ActionRole{
		ActionCode: actionCode,
		ActionName: local_util.NonEmptyString(req.ActionName, old.ActionName),
		Enabled:    old.Enabled,
	}

	if req.AssignedCheckerRoles != nil {
		payload.AssignedMakersRoles = req.AssignedMakersRoles
	}
	if req.AssignedCheckerRoles != nil {
		payload.AssignedCheckerRoles = req.AssignedCheckerRoles
	}

	cpsAction := lib.CpsModelBuilder(
		actionCode,
		maker,
		old,
		payload,
		string(constants.RequestUpdateActionRole),
		constants.UPDATE,
	)
	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *cpsActionRoleService) Enable(ctx context.Context, actionCode string) error {
	if actionCode == "" {
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	// Ensure exists
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if old.Enabled {
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRole{ActionCode: actionCode, Enabled: true}
	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestEnableActionRole), constants.UPDATE)
	return s.cpsService.CreateCPSAction(ctx, &cps)
}

func (s *cpsActionRoleService) Disable(ctx context.Context, actionCode string) error {
	if actionCode == "" {
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if !old.Enabled {
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRole{ActionCode: actionCode, Enabled: false}
	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestDisableActionRole), constants.UPDATE)
	return s.cpsService.CreateCPSAction(ctx, &cps)
}

// Authorize applies approved CPS actions
func (s *cpsActionRoleService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	cur, err := local_util.JsonUnmarshal[model.ActionRole](action.CurrentAction)
	if err != nil {
		s.logger.Errorf("failed to unmarshal action: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
	}
	switch action.ActionType {
	case string(constants.CREATE):
		ar, err := s.bindActionRoleModel(*cur)
		if err != nil {
			return nil, err
		}
		ar.ID = bson.NewObjectID()
		ar.CreatedAt = time.Now()
		ar.UpdatedAt = time.Now()
		if err := s.repo.Create(ctx, &ar); err != nil {
			return nil, err
		}
		return action, nil
	case string(constants.UPDATE):
		existing, err := s.repo.FindByActionCode(ctx, action.UniqueId)
		if err != nil {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		upd, err := s.bindActionRoleModel(*cur)
		if err != nil {
			return nil, err
		}
		upd.UpdatedAt = time.Now()
		upd.CreatedAt = existing.CreatedAt
		if cur.AssignedCheckerRoles == nil && cur.AssignedCheckerRoles == nil && cur.ActionName == "" {
			// enable/disable path
			if err := s.repo.EnableOrDisableByActionCode(ctx, existing.ActionCode, cur.Enabled); err != nil {
				return nil, err
			}
			return action, nil
		}
		if err := s.repo.UpdateByActionCode(ctx, existing.ActionCode, &upd); err != nil {
			return nil, err
		}
		return action, nil
	default:
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *cpsActionRoleService) bindActionRoleModel(in model.ActionRole) (model.CPSActionRole, error) {
	var makers []bson.ObjectID
	for _, id := range in.AssignedMakersRoles {
		// obj, err := bson.ObjectIDFromHex(id)
		// if err != nil {
		// 	return model.ActionRole{}, errors.New(localization.ErrorInvalidID.Code)
		// }
		makers = append(makers, id)
	}
	var checkers [][]bson.ObjectID
	for _, grp := range in.AssignedCheckerRoles {
		var g []bson.ObjectID
		for _, id := range grp {
			// obj, err := bson.ObjectIDFromHex(id)
			// if err != nil {
			// 	return model.ActionRole{}, errors.New(localization.ErrorInvalidID.Code)
			// }
			g = append(g, id)
		}
		checkers = append(checkers, g)
	}
	return model.CPSActionRole{
		ActionCode:           in.ActionCode,
		ActionName:           in.ActionName,
		AssignedMakersRoles:  makers,
		AssignedCheckerRoles: checkers,
		Enabled:              in.Enabled,
	}, nil
}
