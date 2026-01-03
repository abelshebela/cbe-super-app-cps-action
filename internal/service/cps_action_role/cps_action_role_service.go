package cps_action_role_service

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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type cpsActionRoleService struct {
	repo       storage.CPSActionRoleRepository
	cpsService service.CPSActionService
	indexRepo  storage.CPSActionApproveIndexRepository
	roleRepo   storage.JobRoleRepository
	logger     utils.Logger
}

func NewCPSActionRoleService(
	repo storage.CPSActionRoleRepository,
	indexRepo storage.CPSActionApproveIndexRepository,
	roleRepo storage.JobRoleRepository,
	cps service.CPSActionService,
	logger utils.Logger,
) service.CPSActionRoleService {
	return &cpsActionRoleService{
		repo:       repo,
		indexRepo:  indexRepo,
		roleRepo:   roleRepo,
		cpsService: cps,
		logger:     logger,
	}
}

// FindAllWithPagination implements service.cpsActionRoleService.
func (s *cpsActionRoleService) FindAllActionListWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]imodel.CPSActionList], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "CPSActionRole", "FindAllWithPagination")
	defer span.End()
	result, err := s.repo.FindAllAccessListWithPagination(ctx, filter)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

// FindAllWithPagination implements service.cpsActionRoleService.
func (s *cpsActionRoleService) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.CPSActionRoleResposne], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllWithPagination", "CPSActionRole", "FindAllWithPagination")
	defer span.End()
	result, err := s.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

// GetByActionCode implements service.cpsActionRoleService.
func (s *cpsActionRoleService) GetByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetByActionCode", "CPSActionRole", "GetByActionCode")
	defer span.End()
	if actionCode == "" {
		span.AddEvent("action code is empty", trace.WithAttributes(attribute.String("error", "action code is empty")))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	res, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		span.AddEvent("failed to find by action code", trace.WithAttributes(attribute.String("error", err.Error())))
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "CPSActionRole", "Create")
	defer span.End()
	if req.ActionName == "" {
		span.AddEvent("action name is empty", trace.WithAttributes(attribute.String("error", "action name is empty")))
		return errors.New(localization.ErrorActionNameIsRequired.Code)
	}

	req.ActionName = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(req.ActionName), " ", "_"))
	req.ActionCode = req.ActionName
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("incomplete maker information", trace.WithAttributes(attribute.String("error", "incomplete maker information")))
		return errors.New(localization.ErrorUserUnauthorized.Code)
	}

	existing, err := s.repo.FindByActionName(ctx, req.ActionName)
	if err == nil && existing != nil {
		span.AddEvent("action name already exists", trace.WithAttributes(attribute.String("error", "action name already exists")))
		return errors.New(localization.ErrorActionNameAlreadyExists.Code)
	}

	if err := s.validateUniqueIDs(req.AssignedViewersRoles); err != nil {
		span.AddEvent("failed to validate unique viewer IDs", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if err := s.validateUniqueIDs(req.AssignedMakersRoles); err != nil {
		span.AddEvent("failed to validate unique maker IDs", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if err := s.validateUniqueIDsInGroups(req.AssignedCheckerRoles); err != nil {
		span.AddEvent("failed to validate unique checker IDs in groups", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if err := s.validateUniqueIDs(req.AssignedAuditorRoles); err != nil {
		span.AddEvent("failed to validate unique auditor IDs", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	// Viewers
	var viewers []string
	for _, code := range req.AssignedViewersRoles {
		viewers = append(viewers, code)
	}

	// Makers
	var makers []string
	if !req.IsViewOnly {
		for _, code := range req.AssignedMakersRoles {
			makers = append(makers, code)
		}
	}

	// Checkers
	var checkers [][]string
	if !req.IsMakerOnly && !req.IsViewOnly {
		checkers = make([][]string, 0, len(req.AssignedCheckerRoles))
		for _, group := range req.AssignedCheckerRoles {
			g := make([]string, 0, len(group))
			for _, code := range group {
				g = append(g, code)
			}
			checkers = append(checkers, g)
		}
	}

	// Auditor
	auditors := make([]string, 0, len(req.AssignedAuditorRoles))
	if !req.IsViewOnly {
		for _, code := range req.AssignedAuditorRoles {

			auditors = append(auditors, code)
		}
	}

	// build CPS action payload
	payload := imodel.CPSActionRole{
		ActionCode:           req.ActionCode,
		ActionName:           req.ActionName,
		PortalCardName:       req.PortalCardName,
		AssignedMakersRoles:  makers,
		AssignedViewersRoles: viewers,
		AssignedCheckerRoles: checkers,
		AssignedAuditorRoles: auditors,
		IsMakerOnly:          req.IsMakerOnly || int32(len(checkers)) == 0,
		IsViweOnly:           req.IsViewOnly,
		Enabled:              true,
		ApproverCount:        int32(len(checkers)),
	}

	cpsAction := lib.CpsModelBuilder(
		constants.Empty,
		maker,
		nil,
		payload,
		string(constants.RequestCreateCpsActionRole),
		constants.CREATE,
	)

	err = s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

// Update implements service.CPSActionRoleService.
func (s *cpsActionRoleService) Update(ctx context.Context, actionCode string, req actionrole_dto.UpdateActionRoleRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "CPSActionRole", "Update")
	defer span.End()
	if actionCode == "" {
		span.AddEvent("action code is empty", trace.WithAttributes(attribute.String("error", "action code is empty")))
		return errors.New(localization.ErrorActionNameIsRequired.Code)
	}

	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		span.AddEvent("failed to find by action code", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("incomplete maker information", trace.WithAttributes(attribute.String("error", "incomplete maker information")))
		return errors.New(localization.ErrorUserUnauthorized.Code)
	}

	if req.AssignedViewersRoles != nil {
		if err := s.validateUniqueIDs(req.AssignedViewersRoles); err != nil {
			span.AddEvent("failed to validate unique maker IDs", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
	}

	if req.AssignedMakersRoles != nil {
		if err := s.validateUniqueIDs(req.AssignedMakersRoles); err != nil {
			span.AddEvent("failed to validate unique maker IDs", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
	}
	if req.AssignedCheckerRoles != nil {
		if err := s.validateUniqueIDsInGroups(req.AssignedCheckerRoles); err != nil {
			span.AddEvent("failed to validate unique checker IDs in groups", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
	}
	if req.AssignedAuditorRoles != nil {
		if err := s.validateUniqueIDs(req.AssignedAuditorRoles); err != nil {
			span.AddEvent("failed to validate unique auditor IDs", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
	}

	payload := imodel.CPSActionRole{
		ActionCode:     actionCode,
		PortalCardName: req.PortalCardName,
		ActionName:     local_util.NonEmptyString(req.ActionName, old.ActionName),
		IsMakerOnly:    req.IsMakerOnly || int32(len(req.AssignedCheckerRoles)) == 0,
		IsViweOnly:     req.IsViewOnly || (int32(len(req.AssignedViewersRoles)) > 0 && len(req.AssignedViewersRoles) == 0 && len(req.AssignedCheckerRoles) == 0),
		Enabled:        true,
		ApproverCount:  int32(len(req.AssignedCheckerRoles)),
	}

	if req.AssignedViewersRoles != nil {
		viewers := make([]string, 0, len(req.AssignedViewersRoles))
		for _, code := range req.AssignedViewersRoles {
			viewers = append(viewers, code)
		}
		payload.AssignedViewersRoles = viewers
	}

	if req.AssignedMakersRoles != nil {
		makers := make([]string, 0, len(req.AssignedMakersRoles))
		for _, code := range req.AssignedMakersRoles {
			makers = append(makers, code)
		}

		payload.AssignedMakersRoles = makers
	}

	if req.IsMakerOnly {
		payload.AssignedCheckerRoles = [][]string{}
	} else if req.IsViewOnly {
		payload.AssignedMakersRoles = []string{}
		payload.AssignedAuditorRoles = []string{}
		payload.AssignedCheckerRoles = [][]string{}
	} else {
		if req.AssignedCheckerRoles != nil {
			checkers := make([][]string, 0, len(req.AssignedCheckerRoles))
			for _, group := range req.AssignedCheckerRoles {
				g := make([]string, 0, len(group))
				for _, code := range group {
					g = append(g, code)
				}
				checkers = append(checkers, g)
			}
			payload.AssignedCheckerRoles = checkers
		}
	}
	if req.AssignedAuditorRoles != nil {
		auditors := make([]string, 0, len(req.AssignedAuditorRoles))
		for _, code := range req.AssignedAuditorRoles {
			auditors = append(auditors, code)
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
		string(constants.RequestUpdateCpsActionRole),
		constants.UPDATE,
	)

	err = s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (s *cpsActionRoleService) Enable(ctx context.Context, actionCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Enable", "CPSActionRole", "Enable")
	defer span.End()

	if actionCode == "" {
		span.AddEvent("action code is empty", trace.WithAttributes(attribute.String("error", "action code is empty")))
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		span.AddEvent("failed to find by action code", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if old.Enabled {
		span.AddEvent("action role already enabled", trace.WithAttributes(attribute.String("error", "action role already enabled")))
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRole{ActionCode: actionCode, Enabled: true}

	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestEnableActionRole), constants.UPDATE)
	err = s.cpsService.CreateCPSAction(ctx, &cps)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (s *cpsActionRoleService) Disable(ctx context.Context, actionCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "Disable", "CPSActionRole", "Disable")
	defer span.End()
	if actionCode == "" {
		span.AddEvent("action code is empty", trace.WithAttributes(attribute.String("error", "action code is empty")))
		return errors.New(localization.ErrorInvalidInputParameter.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		span.AddEvent("failed to find by action code", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	if !old.Enabled {
		span.AddEvent("action role already disabled", trace.WithAttributes(attribute.String("error", "action role already disabled")))
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.ActionRole{ActionCode: actionCode, Enabled: false}
	cps := lib.CpsModelBuilder(actionCode, maker, old, payload, string(constants.RequestDisableActionRole), constants.UPDATE)
	err = s.cpsService.CreateCPSAction(ctx, &cps)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

// Authorize applies approved CPS actions
func (s *cpsActionRoleService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "CPSActionRole", "Authorize")
	defer span.End()

	switch action.ActionType {
	case string(constants.CREATE):

		cur, err := local_util.JsonUnmarshal[imodel.CPSActionRole](action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to unmarshal action", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("failed to unmarshal action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}

		if err := s.UpdateActionList(ctx, cur.ActionName, true); err != nil {
			span.AddEvent("failed to update action list", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("failed to update action list: %v", err)
			return nil, err
		}

		ar := *cur
		ar.ID = bson.NewObjectID()
		ar.CreatedAt = time.Now()
		ar.UpdatedAt = time.Now()

		if err := s.repo.Create(ctx, &ar); err != nil {
			span.AddEvent("failed to create action role", trace.WithAttributes(attribute.String("error", err.Error())))
			if err := s.UpdateActionList(ctx, cur.ActionName, false); err != nil {
				span.AddEvent("failed to update action list", trace.WithAttributes(attribute.String("error", err.Error())))
				s.logger.Errorf("failed to update action list: %v", err)
			}
			return nil, err
		}

		s.logger.Infof("Authorize: Syncing indices for Create. Makers: %d, Checkers: %d, Auditors: %d", len(ar.AssignedMakersRoles), len(ar.AssignedCheckerRoles), len(ar.AssignedAuditorRoles))
		if err := s.syncIndices(ctx, "", &ar); err != nil {
			span.AddEvent("failed to sync indices for create", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		// Update action.CurrentAction with the new ID and timestamps
		updatedPayload, err := json.Marshal(ar)
		if err != nil {
			span.AddEvent("failed to marshal updated payload", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		action.CurrentAction = updatedPayload
		return action, nil

	case string(constants.UPDATE):
		prev, err := local_util.JsonUnmarshal[model.CPSActionRoleResposne](action.PreviousAction)
		if err != nil {
			span.AddEvent("failed to unmarshal action", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("failed to unmarshal action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}

		new, err := local_util.JsonUnmarshal[imodel.CPSActionRole](action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to unmarshal new action", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("failed to unmarshal new action: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionFormat.Code)
		}

		if err := s.repo.UpdateByActionCode(ctx, prev.ActionCode, new); err != nil {
			span.AddEvent("failed to update action role", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}

		s.logger.Infof("Authorize: Syncing indices for Create. Makers: %d, Checkers: %d, Auditors: %d", len(new.AssignedMakersRoles), len(new.AssignedCheckerRoles), len(new.AssignedAuditorRoles))
		if err := s.syncIndices(ctx, new.ActionCode, new); err != nil {
			span.AddEvent("failed to sync indices for create", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil
	default:
		span.AddEvent("unsupported action type", trace.WithAttributes(attribute.String("action_type", action.ActionType)))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *cpsActionRoleService) UpdateActionList(ctx context.Context, actionCode string, status bool) error {
	err := s.repo.UpdateActionList(ctx, actionCode, status)
	return err
}

func (s *cpsActionRoleService) validateUniqueIDs(ids []string) error {
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

func (s *cpsActionRoleService) validateUniqueIDsInGroups(groups [][]string) error {
	seen := make(map[string]struct{})
	var allIDs []string
	for _, group := range groups {
		for _, id := range group {
			// if _, ok := seen[id]; ok {
			// 	return errors.New(localization.ErrorInvalidInputParameter.Code)
			// }
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
func (s *cpsActionRoleService) syncIndices(ctx context.Context, oldActionName string, role *imodel.CPSActionRole) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "syncIndices", "CPSActionRole", "syncIndices")
	defer span.End()
	span.SetAttributes(attribute.String("old_action_name", oldActionName))
	indices := s.generateIndices(role)
	s.logger.Infof("syncIndices: Generated %d indices for action %s (oldName: %s)", len(indices), role.ActionName, oldActionName)

	if oldActionName == "" {
		return s.indexRepo.SaveIndices(ctx, indices)
	}

	err := s.indexRepo.SyncIndices(ctx, oldActionName, indices)
	if err != nil {
		span.AddEvent("failed to sync indices", trace.WithAttributes(attribute.String("error", err.Error())))
	}
	return err
}

func (s *cpsActionRoleService) generateIndices(role *imodel.CPSActionRole) []imodel.CPSActionApproveIndex {
	_, span := local_util.TraceLogger(context.Background(), "service", "generateIndices", "CPSActionRole", "generateIndices")
	defer span.End()

	var indices []imodel.CPSActionApproveIndex
	now := time.Now()
	s.logger.Infof("generateIndices: Starting for action %s. Viewers: %d,%s. Makers: %d, Checkers: %d, Auditors: %d", role.ActionName, len(role.AssignedViewersRoles), len(role.AssignedMakersRoles), len(role.AssignedCheckerRoles), len(role.AssignedAuditorRoles))

	// Makers
	if len(role.AssignedMakersRoles) > 0 {
		for j, makerID := range role.AssignedMakersRoles {
			idx := int64(j + 1)
			indices = append(indices, imodel.CPSActionApproveIndex{
				ID:             bson.NewObjectID(),
				RoleId:         makerID,
				PortalCardName: role.PortalCardName,
				ActionName:     role.ActionName,
				MakerIndex:     &idx,
				UpdatedAt:      now,
				CreatedAt:      now,
			})
			span.AddEvent("maker index generated", trace.WithAttributes(attribute.String("role_id", makerID)))
			s.logger.Infof("generateIndices: Added Maker index for RoleID %s", makerID)
		}
	}

	// Viewers
	if len(role.AssignedViewersRoles) > 0 {
		for k, viewerID := range role.AssignedViewersRoles {
			idx := int64(k + 1)
			found := false
			for j := range indices {
				if indices[j].RoleId == viewerID {
					indices[j].ViewerIndex = &idx
					found = true
					break
				}
			}
			if !found {
				indices = append(indices, imodel.CPSActionApproveIndex{
					ID:             bson.NewObjectID(),
					RoleId:         viewerID,
					PortalCardName: role.PortalCardName,
					ActionName:     role.ActionName,
					ViewerIndex:    &idx,
					UpdatedAt:      now,
					CreatedAt:      now,
				})

				span.AddEvent("viewer index generated", trace.WithAttributes(attribute.String("role_id", viewerID)))
				s.logger.Infof("generateIndices: Added Viewer index for RoleID %s", viewerID)
			} else {
				span.AddEvent("viewer index updated", trace.WithAttributes(attribute.String("role_id", viewerID)))
				s.logger.Infof("generateIndices: Updated Viewer index for RoleID %s", viewerID)
			}
		}
	}

	// Auditors
	if len(role.AssignedAuditorRoles) > 0 {
		for k, auditorID := range role.AssignedAuditorRoles {
			idx := int64(k + 1)
			found := false
			for j := range indices {
				if indices[j].RoleId == auditorID {
					indices[j].AuditorIndex = &idx
					found = true
					break
				}
			}
			if !found {
				indices = append(indices, imodel.CPSActionApproveIndex{
					ID:             bson.NewObjectID(),
					RoleId:         auditorID,
					PortalCardName: role.PortalCardName,
					ActionName:     role.ActionName,
					AuditorIndex:   &idx,
					UpdatedAt:      now,
					CreatedAt:      now,
				})

				span.AddEvent("auditor index generated", trace.WithAttributes(attribute.String("role_id", auditorID)))
				s.logger.Infof("generateIndices: Added Auditor index for RoleID %s", auditorID)
			} else {
				span.AddEvent("auditor index updated", trace.WithAttributes(attribute.String("role_id", auditorID)))
				s.logger.Infof("generateIndices: Updated Auditor index for RoleID %s", auditorID)
			}
		}
	}

	// Checkers
	if !role.IsMakerOnly {
		for outer, group := range role.AssignedCheckerRoles {
			for inner, checkerID := range group {
				val := float64(outer+1) + float64(inner+1)/10.0
				found := false
				for j := range indices {
					if indices[j].RoleId == checkerID {
						indices[j].CheckerIndex = &val
						found = true
						break
					}
				}
				if !found {
					indices = append(indices, imodel.CPSActionApproveIndex{
						ID:             bson.NewObjectID(),
						RoleId:         checkerID,
						PortalCardName: role.PortalCardName,
						ActionName:     role.ActionName,
						CheckerIndex:   &val,
						UpdatedAt:      now,
						CreatedAt:      now,
					})

					span.AddEvent("checker index generated", trace.WithAttributes(attribute.String("role_id", checkerID), attribute.Float64("value", val)))
					s.logger.Infof("generateIndices: Added Checker index for RoleID %s (val: %f)", checkerID, val)
				} else {
					span.AddEvent("checker index updated", trace.WithAttributes(attribute.String("role_id", checkerID), attribute.Float64("value", val)))
					s.logger.Infof("generateIndices: Updated Checker index for RoleID %s (val: %f)", checkerID, val)
				}
			}
		}
	}
	span.AddEvent("generate indices completed", trace.WithAttributes(attribute.Int("total_indices", len(indices))))
	s.logger.Infof("generateIndices: Completed. Total indices: %d", len(indices))

	return indices
}
