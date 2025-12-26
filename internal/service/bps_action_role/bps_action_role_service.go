package bps_action_role_service

import (
	"cbe-super-app-cps-action/internal/constants"
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bpsActionRoleService struct {
	repo       storage.BPSActionRoleRepository
	indexRepo  storage.BPSActionApproveIndexRepository
	roleRepo   storage.RoleRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewBPSActionRoleService(
	repo storage.BPSActionRoleRepository,
	indexRepo storage.BPSActionApproveIndexRepository,
	roleRepo storage.RoleRepository,
	cps service.CPSActionService,
	logger utils.Logger,
) service.BPSActionRoleService {
	return &bpsActionRoleService{
		repo:       repo,
		indexRepo:  indexRepo,
		roleRepo:   roleRepo,
		cpsService: cps,
		logger:     logger,
	}
}

// FindAllWithPagination implements service.bpsActionRoleService.
func (s *bpsActionRoleService) FindAllActionListWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*imodel.BPSActionList], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllActionListWithPagination", "BPS Action LIST", "FindAllActionListWithPagination")
	defer span.End()

	return s.repo.FindAllAccessListWithPagination(ctx, filter)
}

// FindAllWithPagination implements service.bpsActionRoleService.
func (s *bpsActionRoleService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.ActionRole], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "BPS Action Role", "FindAllWithPagination")
	defer span.End()

	return s.repo.FindAllWithPagination(ctx, filter)
}

// GetByActionCode implements service.bpsActionRoleService.
func (s *bpsActionRoleService) GetByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetByActionCode", "BPS Action Role", "GetByActionCode")
	defer span.End()

	if actionCode == "" {
		span.AddEvent("[GetByActionCode] action code is empty")
		s.logger.Errorf("[GetByActionCode] action code is empty")
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	res, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code == localization.ErrorResourceNotFound.Code {
			span.AddEvent("[GetByActionCode] action role not found", trace.WithAttributes(attribute.String("action_code", actionCode)))
			s.logger.Errorf("[GetByActionCode] action role not found: %s", actionCode)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("[GetByActionCode] failed to fetch action role", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", actionCode),
		))
		s.logger.Errorf("[GetByActionCode] failed to fetch action role: %v", err)
		return nil, err
	}
	s.logger.Infof("[GetByActionCode] action role retrieved successfully for action code: %s", actionCode)
	return res, nil
}

// Create implements service.bpsActionRoleService.

func (s *bpsActionRoleService) Create(ctx context.Context, req actionrole_dto.CreateActionRoleRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "BPS Action Role", "Create")
	defer span.End()

	if req.ActionName == "" {
		span.AddEvent("[Create] action name is required")
		return errors.New(localization.ErrorActionNameIsRequired.Code)
	}
	// Format Action Name: Uppercase and replace spaces with underscores
	req.ActionName = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(req.ActionName), " ", "_"))

	// Generate Action Code from Action Name
	req.ActionCode = req.ActionName

	// Check if Action Name already exists
	existing, err := s.repo.FindByActionName(ctx, req.ActionName)
	if err == nil && existing != nil {
		span.AddEvent("[Create] action name already exists", trace.WithAttributes(attribute.String("action_name", req.ActionName)))
		return errors.New(localization.ErrorActionNameAlreadyExists.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("[Create] incomplete user data")
		s.logger.Errorf("[Create] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	if err := s.validateUniqueIDs(req.AssignedMakersRoles); err != nil {
		return err
	}
	if err := s.validateUniqueIDsInGroups(req.AssignedCheckerRoles); err != nil {
		return err
	}
	if err := s.validateUniqueIDs(req.AssignedAuditorRoles); err != nil {
		return err
	}

	// Convert strings to ObjectIDs
	makers := make([]bson.ObjectID, 0, len(req.AssignedMakersRoles))
	for _, id := range req.AssignedMakersRoles {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorInvalidID.Code)
		}
		makers = append(makers, oid)
	}

	var checkers [][]bson.ObjectID
	if !req.IsMakerOnly {
		checkers = make([][]bson.ObjectID, 0, len(req.AssignedCheckerRoles))
		for _, group := range req.AssignedCheckerRoles {
			g := make([]bson.ObjectID, 0, len(group))
			for _, id := range group {
				oid, err := bson.ObjectIDFromHex(id)
				if err != nil {
					return errors.New(localization.ErrorInvalidID.Code)
				}
				g = append(g, oid)
			}
			checkers = append(checkers, g)
		}
	}

	auditors := make([]bson.ObjectID, 0, len(req.AssignedAuditorRoles))
	for _, id := range req.AssignedAuditorRoles {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorInvalidID.Code)
		}
		auditors = append(auditors, oid)
	}

	// build CPS action payload
	payload := model.ActionRole{
		ActionCode:           req.ActionCode,
		ActionName:           req.ActionName,
		AssignedMakersRoles:  makers,
		AssignedCheckerRoles: checkers,
		AssignedAuditorRoles: auditors,
		IsMakerOnly:          req.IsMakerOnly,
		Enabled:              true,
		ApproverCount:        int32(len(checkers)),
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
		span.AddEvent("[Create] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", req.ActionCode),
		))
		s.logger.Errorf("[Create] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Create] action role creation request created successfully for action code: %s", req.ActionCode)
	return nil
}

// Update implements service.bpsActionRoleService.
func (s *bpsActionRoleService) Update(ctx context.Context, actionCode string, req actionrole_dto.UpdateActionRoleRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "BPS Action Role", "Update")
	defer span.End()

	s.logger.Infof("[Update] updating action role for action code: %s", actionCode)
	if actionCode == "" {
		span.AddEvent("[Update] action code is empty")
		s.logger.Errorf("[Update] action code is empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		span.AddEvent("[Update] failed to find action role", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", actionCode),
		))
		s.logger.Errorf("[Update] failed to find action role: %v", err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("[Update] incomplete user data")
		s.logger.Errorf("[Update] incomplete user data")
		return errors.New(localization.ErrorUserUnauthorized.Code)
	}

	if req.AssignedMakersRoles != nil {
		if err := s.validateUniqueIDs(req.AssignedMakersRoles); err != nil {
			return err
		}
	}
	if req.AssignedCheckerRoles != nil {
		if err := s.validateUniqueIDsInGroups(req.AssignedCheckerRoles); err != nil {
			return err
		}
	}
	if req.AssignedAuditorRoles != nil {
		if err := s.validateUniqueIDs(req.AssignedAuditorRoles); err != nil {
			return err
		}
	}

	payload := model.ActionRole{
		ActionCode:  actionCode,
		ActionName:  local_util.NonEmptyString(req.ActionName, old.ActionName),
		Enabled:     old.Enabled,
		IsMakerOnly: req.IsMakerOnly,
	}

	// Handle Makers
	if req.AssignedMakersRoles != nil {
		makers := make([]bson.ObjectID, 0, len(req.AssignedMakersRoles))
		for _, id := range req.AssignedMakersRoles {
			oid, err := bson.ObjectIDFromHex(id)
			if err != nil {
				return errors.New(localization.ErrorInvalidID.Code)
			}
			makers = append(makers, oid)
		}
		payload.AssignedMakersRoles = makers
	} else {
		// Keep existing
		makers := make([]bson.ObjectID, len(old.AssignedMakersRoles))
		for i, r := range old.AssignedMakersRoles {
			makers[i] = r.ID
		}
		payload.AssignedMakersRoles = makers
	}

	// Handle Checkers
	if req.IsMakerOnly {
		payload.AssignedCheckerRoles = [][]bson.ObjectID{}
	} else {
		if req.AssignedCheckerRoles != nil {
			checkers := make([][]bson.ObjectID, 0, len(req.AssignedCheckerRoles))
			for _, group := range req.AssignedCheckerRoles {
				g := make([]bson.ObjectID, 0, len(group))
				for _, id := range group {
					oid, err := bson.ObjectIDFromHex(id)
					if err != nil {
						return errors.New(localization.ErrorInvalidID.Code)
					}
					g = append(g, oid)
				}
				checkers = append(checkers, g)
			}
			payload.AssignedCheckerRoles = checkers
		} else {
			// Keep existing
			checkers := make([][]bson.ObjectID, len(old.AssignedCheckerRoles))
			for i, group := range old.AssignedCheckerRoles {
				g := make([]bson.ObjectID, len(group))
				for j, r := range group {
					g[j] = r.ID
				}
				checkers[i] = g
			}
			payload.AssignedCheckerRoles = checkers
		}
	}

	// Handle Auditors
	if req.AssignedAuditorRoles != nil {
		auditors := make([]bson.ObjectID, 0, len(req.AssignedAuditorRoles))
		for _, id := range req.AssignedAuditorRoles {
			oid, err := bson.ObjectIDFromHex(id)
			if err != nil {
				return errors.New(localization.ErrorInvalidID.Code)
			}
			auditors = append(auditors, oid)
		}
		payload.AssignedAuditorRoles = auditors
	} else {
		// Keep existing
		auditors := make([]bson.ObjectID, len(old.AssignedAuditorRoles))
		for i, r := range old.AssignedAuditorRoles {
			auditors[i] = r.ID
		}
		payload.AssignedAuditorRoles = auditors
	}

	// Recalculate ApproverCount
	payload.ApproverCount = int32(len(payload.AssignedCheckerRoles))

	cpsAction := lib.CpsModelBuilder(
		actionCode,
		maker,
		old,
		payload,
		string(constants.RequestUpdateActionRole),
		constants.UPDATE,
	)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("[Update] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", actionCode),
		))
		s.logger.Errorf("[Update] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Update] action role update request created successfully for action code: %s", actionCode)
	return nil
}

func (s *bpsActionRoleService) Enable(ctx context.Context, actionCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Enable", "BPS Action Role", "Enable")
	defer span.End()

	s.logger.Infof("[Enable] enabling action role for action code: %s", actionCode)
	if actionCode == "" {
		span.AddEvent("[Enable] action code is empty")
		s.logger.Errorf("[Enable] action code is empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	// Ensure exists
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		span.AddEvent("[Enable] failed to find action role", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", actionCode),
		))
		s.logger.Errorf("[Enable] failed to find action role: %v", err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if old.Enabled {
		span.AddEvent("[Enable] action role already enabled", trace.WithAttributes(attribute.String("action_code", actionCode)))
		s.logger.Errorf("[Enable] action role already enabled")
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRole{ActionCode: actionCode, Enabled: true}
	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestEnableActionRole), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cps); err != nil {
		span.AddEvent("[Enable] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", actionCode),
		))
		s.logger.Errorf("[Enable] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Enable] action role enable request created successfully for action code: %s", actionCode)
	return nil
}

func (s *bpsActionRoleService) Disable(ctx context.Context, actionCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Disable", "BPS Action Role", "Disable")
	defer span.End()

	s.logger.Infof("[Disable] disabling action role for action code: %s", actionCode)
	if actionCode == "" {
		span.AddEvent("[Disable] action code is empty")
		s.logger.Errorf("[Disable] action code is empty")
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		span.AddEvent("[Disable] failed to find action role", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", actionCode),
		))
		s.logger.Errorf("[Disable] failed to find action role: %v", err)
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if !old.Enabled {
		span.AddEvent("[Disable] action role already disabled", trace.WithAttributes(attribute.String("action_code", actionCode)))
		s.logger.Errorf("[Disable] action role already disabled")
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRole{ActionCode: actionCode, Enabled: false}
	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestDisableActionRole), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cps); err != nil {
		span.AddEvent("[Disable] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("action_code", actionCode),
		))
		s.logger.Errorf("[Disable] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[Disable] action role disable request created successfully for action code: %s", actionCode)
	return nil
}

// Authorize applies approved CPS actions
func (s *bpsActionRoleService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "BPS Action Role", "Authorize")
	defer span.End()

	s.logger.Infof("[Authorize] authorizing action role action: %s", action.RequestAction)

	cur, err := local_util.JsonUnmarshal[model.ActionRole](action.CurrentAction)

	if err != nil {
		span.AddEvent("[Authorize] failed to unmarshal action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		s.logger.Errorf("[Authorize] failed to unmarshal action: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
	}
	switch action.ActionType {
	case string(constants.CREATE):
		ar, err := s.bindActionRoleModel(*cur)
		if err != nil {
			span.AddEvent("[Authorize] failed to bind action role model", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to bind action role model: %v", err)
			return nil, err
		}
		ar.ID = bson.NewObjectID()
		ar.CreatedAt = time.Now()
		ar.UpdatedAt = time.Now()

		if err := s.UpdateActionList(ctx, cur.ActionName, true); err != nil {
			span.AddEvent("failed to update action list", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("failed to update action list: %v", err)
			return nil, err
		}

		if err := s.repo.Create(ctx, &ar); err != nil {
			span.AddEvent("[Authorize] failed to create action role", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))

			if err := s.UpdateActionList(ctx, cur.ActionName, false); err != nil {
				span.AddEvent("failed to update action list", trace.WithAttributes(attribute.String("error", err.Error())))
				s.logger.Errorf("failed to update action list: %v", err)
				return nil, err
			}
			s.logger.Errorf("[Authorize] failed to create action role: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] action role created successfully") // Sync indices
		s.logger.Infof("Authorize: Syncing indices for Create. Makers: %d, Checkers: %d, Auditors: %d", len(ar.AssignedMakersRoles), len(ar.AssignedCheckerRoles), len(ar.AssignedAuditorRoles))
		if err := s.syncIndices(ctx, "", &ar); err != nil {
			return nil, err
		}
		// Update action.CurrentAction with the new ID and timestamps
		updatedPayload, err := json.Marshal(ar)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		action.CurrentAction = updatedPayload
		return action, nil
	case string(constants.UPDATE):
		existing, err := s.repo.FindByActionCode(ctx, action.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] failed to find existing action role", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to find existing action role: %v", err)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		upd, err := s.bindActionRoleModel(*cur)
		if err != nil {
			span.AddEvent("[Authorize] failed to bind action role model", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to bind action role model: %v", err)
			return nil, err
		}
		upd.UpdatedAt = time.Now()
		upd.CreatedAt = existing.CreatedAt
		if cur.AssignedMakersRoles == nil && cur.AssignedCheckerRoles == nil && cur.ActionName == "" {
			// enable/disable path
			if err := s.repo.EnableOrDisableByActionCode(ctx, existing.ActionCode, cur.Enabled); err != nil {
				span.AddEvent("[Authorize] failed to enable/disable action role", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("unique_id", action.UniqueId),
				))
				s.logger.Errorf("[Authorize] failed to enable/disable action role: %v", err)
				return nil, err
			}
			s.logger.Infof("[Authorize] action role enable/disable completed successfully")
			return action, nil
		}
		if err := s.repo.UpdateByActionCode(ctx, existing.ActionCode, &upd); err != nil {
			span.AddEvent("[Authorize] failed to update action role", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to update action role: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] action role updated successfully") // Sync indices
		if err := s.syncIndices(ctx, existing.ActionName, &upd); err != nil {
			return nil, err
		}
		// Update action.CurrentAction with the updated details
		// Ensure timestamps are set in the payload
		upd.CreatedAt = existing.CreatedAt
		upd.UpdatedAt = time.Now()

		updatedPayload, err := json.Marshal(upd)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		action.CurrentAction = updatedPayload
		return action, nil
	default:
		span.AddEvent("[Authorize] unsupported action type", trace.WithAttributes(attribute.String("action_type", action.ActionType)))
		s.logger.Errorf("[Authorize] unsupported action type: %s", action.ActionType)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *bpsActionRoleService) UpdateActionList(ctx context.Context, actionCode string, status bool) error {
	err := s.repo.UpdateActionList(ctx, actionCode, status)
	return err
}

func (s *bpsActionRoleService) syncIndices(ctx context.Context, oldActionName string, role *model.ActionRole) error {
	indices := s.generateIndices(role)
	s.logger.Infof("syncIndices: Generated %d indices for action %s (oldName: %s)", len(indices), role.ActionName, oldActionName)
	if oldActionName == "" {
		return s.indexRepo.SaveIndices(ctx, indices)
	}
	return s.indexRepo.SyncIndices(ctx, oldActionName, indices)
}

func (s *bpsActionRoleService) generateIndices(role *model.ActionRole) []model.BPSActionApproveIndex {
	var indices []model.BPSActionApproveIndex
	now := time.Now()
	s.logger.Infof("generateIndices: Starting for action %s. Makers: %d, Checkers: %d, Auditors: %d", role.ActionName, len(role.AssignedMakersRoles), len(role.AssignedCheckerRoles), len(role.AssignedAuditorRoles))
	// Makers
	for i, makerID := range role.AssignedMakersRoles {
		idx := int64(i)
		indices = append(indices, model.BPSActionApproveIndex{
			ID:         bson.NewObjectID(),
			RoleId:     makerID.Hex(),
			ActionName: role.ActionName,
			MakerIndex: &idx,
			UpdatedAt:  now,
			CreatedAt:  now,
		})
		s.logger.Infof("generateIndices: Added Maker index for RoleID %s", makerID.Hex())
	}

	// Auditors
	for i, auditorID := range role.AssignedAuditorRoles {
		idx := int64(i)
		found := false
		for j := range indices {
			if indices[j].RoleId == auditorID.Hex() {
				indices[j].AuditorIndex = &idx
				found = true
				break
			}
		}
		if !found {
			indices = append(indices, model.BPSActionApproveIndex{
				ID:           bson.NewObjectID(),
				RoleId:       auditorID.Hex(),
				ActionName:   role.ActionName,
				AuditorIndex: &idx,
				UpdatedAt:    now,
				CreatedAt:    now,
			})
			s.logger.Infof("generateIndices: Added Auditor index for RoleID %s", auditorID.Hex())
		} else {
			s.logger.Infof("generateIndices: Updated Auditor index for RoleID %s", auditorID.Hex())
		}
	}

	// Checkers
	if !role.IsMakerOnly {
		for outer, group := range role.AssignedCheckerRoles {
			for inner, checkerID := range group {
				val := float64(outer) + float64(inner)/10.0
				found := false
				for j := range indices {
					if indices[j].RoleId == checkerID.Hex() {
						indices[j].CheckerIndex = &val
						found = true
						break
					}
				}
				if !found {
					indices = append(indices, model.BPSActionApproveIndex{
						ID:           bson.NewObjectID(),
						RoleId:       checkerID.Hex(),
						ActionName:   role.ActionName,
						CheckerIndex: &val,
						UpdatedAt:    now,
						CreatedAt:    now,
					})
					s.logger.Infof("generateIndices: Added Checker index for RoleID %s (val: %f)", checkerID.Hex(), val)
				} else {
					s.logger.Infof("generateIndices: Updated Checker index for RoleID %s (val: %f)", checkerID.Hex(), val)
				}
			}
		}
	}
	s.logger.Infof("generateIndices: Completed. Total indices: %d", len(indices))

	return indices
}

func (s *bpsActionRoleService) bindActionRoleModel(in model.ActionRole) (model.ActionRole, error) {
	// Since model.ActionRole already has []bson.ObjectID, and json.Unmarshal handles the conversion from hex strings,
	// we just need to return the input.
	return in, nil
}

func (s *bpsActionRoleService) validateUniqueIDs(ids []string) error {
	seen := make(map[string]struct{})
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return errors.New(localization.ErrorInvalidInputParameter.Code)
		}
		seen[id] = struct{}{}
	}
	// Check existence in DB
	exists, err := s.roleRepo.ExistsMany(context.Background(), ids)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	if !exists {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (s *bpsActionRoleService) validateUniqueIDsInGroups(groups [][]string) error {
	seen := make(map[string]struct{})
	var allIDs []string
	for _, group := range groups {
		for _, id := range group {
			if _, ok := seen[id]; ok {
				return errors.New(localization.ErrorInvalidInputParameter.Code)
			}
			seen[id] = struct{}{}
			allIDs = append(allIDs, id)
		}
	}
	// Check existence in DB
	exists, err := s.roleRepo.ExistsMany(context.Background(), allIDs)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	if !exists {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}
