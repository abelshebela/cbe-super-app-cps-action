package cpsaction

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/constants"
	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/constants/localization"

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

func (s *cpsActionServiceWithRoles) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	actType := strings.ToUpper(strings.TrimSpace(cpsAction.ActionType))
	req := strings.ToUpper(strings.TrimSpace(cpsAction.RequestAction))

	if actType == string(constants.ActionCreate) || actType == string(constants.ActionUpdate) || actType == string(constants.ActionDelete) ||
		strings.Contains(req, "ENABLE") || strings.Contains(req, "DISABLE") {

		if cpsAction.CheckerUsers == nil {
			cpsAction.CheckerUsers = []model.Checker{}
		}

		cpsAction.CurrentCheckerIndex = 0.0

		if mod, ok := ResolveModuleForRA(RequestAction(cpsAction.RequestAction)); ok && s.roles != nil {

			var role *model.CPSActionRole

			if r, err := s.roles.FindByActionName(ctx, strings.ToUpper(mod)); err == nil && r != nil {
				role = r
			}

			if role == nil {
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
					s.base.ApproveCPSAction(ctx, cpsAction)
				}

			}

		}
	}

	return nil
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
