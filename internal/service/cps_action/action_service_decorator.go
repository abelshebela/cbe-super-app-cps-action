package cpsaction

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/constants"
	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
)

type cpsActionServiceWithRoles struct {
	base  service.CPSActionService
	roles storage.CPSActionRoleRepository
}

// WithActionRolePolicy wraps a base CPSActionService and injects CPSActionRole
// policy on CPSAction creation (checker_count and action_status).
func WithActionRolePolicy(base service.CPSActionService, roles storage.CPSActionRoleRepository) service.CPSActionService {
	return &cpsActionServiceWithRoles{base: base, roles: roles}
}

func (s *cpsActionServiceWithRoles) IsMakerOnlyForRequest(ctx context.Context, requestAction string) (bool, error) {
	return s.base.IsMakerOnlyForRequest(ctx, requestAction)
}

func (s *cpsActionServiceWithRoles) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	actType := strings.ToUpper(strings.TrimSpace(cpsAction.ActionType))
	req := strings.ToUpper(strings.TrimSpace(cpsAction.RequestAction))
	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	if actType == string(constants.ActionCreate) || actType == string(constants.ActionUpdate) || actType == string(constants.ActionDelete) ||
		strings.Contains(req, "ENABLE") || strings.Contains(req, "DISABLE") {

		if cpsAction.CheckerUsers == nil {
			cpsAction.CheckerUsers = []model.Checker{}
		}

		cpsAction.CurrentCheckerIndex = 0.0

		if mod, ok := ResolveModuleForRA(RequestAction(cpsAction.RequestAction)); ok && s.roles != nil {

			var role *imodel.CPSActionRole
			var approverData imodel.CPSActionApproveIndex

			if r, err := s.roles.FindByActionName(ctx, mod); err == nil && r != nil {
				role = r
			}
			if role != nil {
				ctx = context.WithValue(ctx, constants.ContextKey("is_maker_only"), role.IsMakerOnly)
				ctx = context.WithValue(ctx, constants.ContextKey("action_name"), mod)
				types.SetIsMakerOnly(ctx, role.IsMakerOnly)
			}
			if approver, err := s.roles.FindApproverByActionName(ctx, strings.ToUpper(mod), roleCode); err == nil {
				approverData = approver
			}

			if role == nil || approverData.ID.IsZero() {
				return localization.ErrorOperationNotAllowed
			}

			if approverData.MakerIndex == nil {
				return localization.ErrorOperationNotAllowed
			}

			if role.IsMakerOnly {
				cpsAction.CheckerCount = 0
			} else if role.ApproverCount > 0 {
				cpsAction.CheckerCount = int32(role.ApproverCount)
			} else {
				cpsAction.CheckerCount = 0
			}

			if err := s.base.CreateCPSAction(ctx, cpsAction); err != nil {
				return err
			}

			if role != nil {
				fmt.Printf("CPS Action created with role policy: %+v\n", role)
				if role.IsMakerOnly {
					cpsAction.ActionStatus = string(constants.Approved)
					if err := s.base.ApproveCPSAction(ctx, cpsAction); err != nil {
						return err
					}
				}

			}

		}
	}

	return nil
}

func (s *cpsActionServiceWithRoles) AuditorClaim(ctx context.Context, actionCode string, activeGroup int) error {
	return nil
}
func (s *cpsActionServiceWithRoles) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int) error {
	if err := s.base.AuditorMark(ctx, actionCode, auditor, activeGroup); err != nil {
		return err
	}
	return nil
}

func (s *cpsActionServiceWithRoles) GetUserAuthorizerIndex(ctx context.Context, requestAction constants.RequestAction) (imodel.CPSActionApproveIndex, error) {
	return s.base.GetUserAuthorizerIndex(ctx, requestAction)
}
func (s *cpsActionServiceWithRoles) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	return s.base.GetUserCheckedActions(ctx, userID, filterParams)
}
func (s *cpsActionServiceWithRoles) GetUserCreatedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	return s.base.GetUserCreatedActions(ctx, userID, filterParams)
}

func (s *cpsActionServiceWithRoles) GetCPSActionsForApprover(ctx context.Context, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	return s.base.GetCPSActionsForApprover(ctx, RAList, filterParams)
}

func (s *cpsActionServiceWithRoles) GetCPSActionsForAuditor(ctx context.Context, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	return s.base.GetCPSActionsForAuditor(ctx, RAList, filterParams)
}

func (s *cpsActionServiceWithRoles) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {
	return s.base.ApproveCPSAction(ctx, action)
}

func (s *cpsActionServiceWithRoles) RejectCPSAction(ctx context.Context, actionCode string, action *model.CPSAction) error {
	return s.base.RejectCPSAction(ctx, actionCode, action)
}

func (s *cpsActionServiceWithRoles) CancelCPSAction(ctx context.Context, actionCode string, action *model.CPSAction) error {
	return s.base.CancelCPSAction(ctx, actionCode, action)
}

func (s *cpsActionServiceWithRoles) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	return s.base.GetCPSActionsByDepartment(ctx, department, filterParams)
}

func (s *cpsActionServiceWithRoles) GetActionCountsByDepartemnt(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error) {
	// Delegate; type alias not available here, forward to base
	return s.base.GetActionCountsByDepartemnt(ctx, department)
}

func (s *cpsActionServiceWithRoles) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	return s.base.GetCPSActionByID(ctx, id, department)
}

func (s *cpsActionServiceWithRoles) GetCPSActionByUniqueID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	return s.base.GetCPSActionByUniqueID(ctx, id, department)
}

func (s *cpsActionServiceWithRoles) GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	return s.base.GetCPSActionByActionCode(ctx, uniqueID, department)
}

// ReverseCPSAction delegates to the base implementation to satisfy service.CPSActionService
func (s *cpsActionServiceWithRoles) ReverseCPSAction(ctx context.Context, actionCode string) error {
	return s.base.ReverseCPSAction(ctx, actionCode)
}
