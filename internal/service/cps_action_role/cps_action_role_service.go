package bps_action_role_service

import (
	"cbe-super-app-cps-action/internal/constants"
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
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
)

type cpsActionRoleService struct {
	repo       storage.CPSActionRoleRepository
	cpsService service.CPSActionService
	indexRepo  storage.CPSActionApproveIndexRepository
	roleRepo   storage.RoleRepository
	logger     utils.Logger
}

func NewCPSActionRoleService(
	repo storage.CPSActionRoleRepository,
	indexRepo storage.CPSActionApproveIndexRepository,
	roleRepo storage.RoleRepository,
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
	existing, err := s.repo.FindByActionName(ctx, req.ActionName)
	if err == nil && existing != nil {
		return errors.New(localization.ErrorActionNameAlreadyExists.Code)
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
	var makers []bson.ObjectID
	for _, id := range req.AssignedMakersRoles {
		obj, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return errors.New(localization.ErrorInvalidID.Code)
		}
		makers = append(makers, obj)
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
		string(constants.RequestCreateCpsActionRole),
		constants.CREATE,
	)
	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

// Update implements service.CPSActionRoleService.
func (s *cpsActionRoleService) Update(ctx context.Context, actionCode string, req actionrole_dto.UpdateActionRoleRequest) error {
	if actionCode == "" {
		return errors.New(localization.ErrorActionNameIsRequired.Code)
	}
	old, err := s.repo.FindByActionCode(ctx, actionCode)
	if err != nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
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
		string(constants.RequestUpdateCpsActionRole),
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
	cur, err := local_util.JsonUnmarshal[model.CPSActionRole](action.CurrentAction)
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
		s.logger.Infof("Authorize: Syncing indices for Create. Makers: %d, Checkers: %d, Auditors: %d", len(ar.AssignedMakersRoles), len(ar.AssignedCheckersRoles), len(ar.AssignedAuditorRoles))
		if err := s.syncIndices(ctx, "", &ar); err != nil {
			return nil, err
		}
		// Update action.CurrentAction with the new ID and timestamps
		updatedPayload, err := json.Marshal(ar)
		if err != nil {
			return nil, err
		}
		action.CurrentAction = updatedPayload
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
		if cur.AssignedMakersRoles == nil && cur.AssignedCheckersRoles == nil && cur.ActionName == "" {
			// enable/disable path
			if err := s.repo.EnableOrDisableByActionCode(ctx, existing.ActionCode, cur.Enabled); err != nil {
				return nil, err
			}
			return action, nil
		}
		if err := s.repo.UpdateByActionCode(ctx, existing.ActionCode, &upd); err != nil {
			return nil, err
		}
		// Sync indices
		if err := s.syncIndices(ctx, existing.ActionName, &upd); err != nil {
			return nil, err
		}
		upd.CreatedAt = existing.CreatedAt
		upd.UpdatedAt = time.Now()

		updatedPayload, err := json.Marshal(upd)
		if err != nil {
			return nil, err
		}
		action.CurrentAction = updatedPayload
		return action, nil
	default:
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}

func (s *cpsActionRoleService) bindActionRoleModel(in model.CPSActionRole) (model.CPSActionRole, error) {
	// Since model.ActionRole already has []bson.ObjectID, and json.Unmarshal handles the conversion from hex strings,
	// we just need to return the input.
	return in, nil
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
func (s *cpsActionRoleService) syncIndices(ctx context.Context, oldActionName string, role *model.CPSActionRole) error {
	indices := s.generateIndices(role)
	s.logger.Infof("syncIndices: Generated %d indices for action %s (oldName: %s)", len(indices), role.ActionName, oldActionName)
	if oldActionName == "" {
		return s.indexRepo.SaveIndices(ctx, indices)
	}
	return s.indexRepo.SyncIndices(ctx, oldActionName, indices)
}

func (s *cpsActionRoleService) generateIndices(role *model.CPSActionRole) []model.CPSActionApproveIndex {
	var indices []model.CPSActionApproveIndex
	now := time.Now()
	s.logger.Infof("generateIndices: Starting for action %s. Makers: %d, Checkers: %d, Auditors: %d", role.ActionName, len(role.AssignedMakersRoles), len(role.AssignedCheckersRoles), len(role.AssignedAuditorRoles))
	// Makers
	for i, makerID := range role.AssignedMakersRoles {
		idx := int64(i)
		indices = append(indices, model.CPSActionApproveIndex{
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
			indices = append(indices, model.CPSActionApproveIndex{
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
		for outer, group := range role.AssignedCheckersRoles {
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
					indices = append(indices, model.CPSActionApproveIndex{
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
