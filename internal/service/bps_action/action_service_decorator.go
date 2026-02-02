package bps_action

import (
	"context"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/constants"
	bpsActionDto "cbe-super-app-cps-action/internal/constants/dto/bps_action"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
)

type bpsActionServiceWithRoles struct {
	base  service.BPSActionService
	roles storage.BPSActionRoleRepository
}

// WithActionRolePolicy wraps a base CPSActionService and injects CPSActionRole
// policy on CPSAction creation (checker_count and action_status).
func WithActionRolePolicy(base service.BPSActionService, roles storage.BPSActionRoleRepository) service.BPSActionService {
	return &bpsActionServiceWithRoles{base: base, roles: roles}
}

func (s *bpsActionServiceWithRoles) IsMakerOnlyForRequest(ctx context.Context, requestAction string) (bool, error) {
	return s.base.IsMakerOnlyForRequest(ctx, requestAction)
}

func (s *bpsActionServiceWithRoles) AuditorClaim(ctx context.Context, actionCode string, activeGroup int) error {
	return nil
}
func (s *bpsActionServiceWithRoles) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int) error {
	if err := s.base.AuditorMark(ctx, actionCode, auditor, activeGroup); err != nil {
		return err
	}
	return nil
}

func (s *bpsActionServiceWithRoles) GetUserAuthorizerIndex(ctx context.Context, requestAction constants.RequestAction) (imodel.BPSActionApproveIndex, error) {
	return s.base.GetUserAuthorizerIndex(ctx, requestAction)
}
func (s *bpsActionServiceWithRoles) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	return s.base.GetUserCheckedActions(ctx, userID, filterParams)
}
func (s *bpsActionServiceWithRoles) GetUserCreatedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	return s.base.GetUserCreatedActions(ctx, userID, filterParams)
}

func (s *bpsActionServiceWithRoles) GetBPSActionsForApprover(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	return s.base.GetBPSActionsForApprover(ctx, userID, RAList, filterParams)
}

func (s *bpsActionServiceWithRoles) GetBPSActionsForAuditor(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	return s.base.GetBPSActionsForAuditor(ctx, userID, RAList, filterParams)
}

func (s *bpsActionServiceWithRoles) GetBPSActions(ctx context.Context, userID, role string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	return s.base.GetBPSActions(ctx, userID, role, RAList, filterParams)
}

func (s *bpsActionServiceWithRoles) ApproveBPSAction(ctx context.Context, action *bps_model.BPSAction) error {
	return s.base.ApproveBPSAction(ctx, action)
}

func (s *bpsActionServiceWithRoles) RejectBPSAction(ctx context.Context, actionCode string, action *bps_model.BPSAction) error {
	return s.base.RejectBPSAction(ctx, actionCode, action)
}

func (s *bpsActionServiceWithRoles) GetBPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	return s.base.GetBPSActionsByDepartment(ctx, department, filterParams)
}

func (s *bpsActionServiceWithRoles) GetActionCountsByDepartemnt(ctx context.Context, department string) (*bpsActionDto.BPSActionCountResponse, error) {
	// Delegate; type alias not available here, forward to base
	return s.base.GetActionCountsByDepartemnt(ctx, department)
}

func (s *bpsActionServiceWithRoles) GetBPSActionByID(ctx context.Context, id, department string) (*bps_model.BPSAction, error) {
	return s.base.GetBPSActionByID(ctx, id, department)
}

func (s *bpsActionServiceWithRoles) GetBPSActionByUniqueID(ctx context.Context, id, department string) (*bps_model.BPSAction, error) {
	return s.base.GetBPSActionByUniqueID(ctx, id, department)
}

func (s *bpsActionServiceWithRoles) GetBPSActionByActionCode(ctx context.Context, uniqueID, department string) (*bps_model.BPSAction, error) {
	return s.base.GetBPSActionByActionCode(ctx, uniqueID, department)
}
