package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"strings"

	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type cpsActionService struct {
	repo       storage.CPSActionRepository
	roles      storage.CPSActionRoleRepository
	logger     utils.Logger
	dispatcher Dispatcher
}

// IsMakerOnlyForRequest returns true if the module mapped from requestAction is configured as maker-only in CPSActionRole.
func (ca *cpsActionService) IsMakerOnlyForRequest(ctx context.Context, requestAction string) (bool, error) {
	if mod, ok := ResolveModuleForRA(RequestAction(requestAction)); ok && ca.roles != nil {
		role, err := ca.roles.FindByActionName(ctx, mod)
		if err != nil || role == nil {
			return false, errors.New(localization.ErrorOperationNotAllowed.Code)
		}
		return role.IsMakerOnly, nil
	}
	return false, errors.New(localization.ErrorOperationNotAllowed.Code)
}

// AuditorClaim sets auditor status to INPROGRESS when caller belongs to the active group.
func (ca *cpsActionService) AuditorClaim(ctx context.Context, actionCode string, activeGroup int) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuditorClaim", "CPSAction", "AuditorClaim")
	defer span.End()
	act, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": actionCode})
	if err != nil || act == nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	// compute active group from record
	current := int64(0)
	if act.CurrentAuditorIndex > 0 {
		current = int64(act.CurrentAuditorIndex)
	}

	expected := int64(activeGroup)
	if current != 0 && current != expected {
		return errors.New(localization.ErrorOperationNotAllowed.Code)
	}
	// idempotent move to INPROGRESS
	upd := model.CPSAction{ActionCode: actionCode}
	upd.AuditorStatus = "INPROGRESS"
	_, err = ca.repo.Update(ctx, actionCode, upd)
	return err
}

// AuditorMark records an auditor's mark and advances to the next group or finishes.
func (ca *cpsActionService) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuditorMark", "CPSAction", "AuditorMark")
	defer span.End()
	act, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": actionCode})
	if err != nil || act == nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	// prevent multiple marks within the same group (any-one quorum)
	grp := int(activeGroup)
	for _, au := range act.AuditorUsers {
		if int(au.AuditorIndex) == grp {
			return errors.New(localization.ErrorOperationNotAllowed.Code)
		}
	}

	// append this auditor
	nextUsers := append(act.AuditorUsers, auditor)

	upd := model.CPSAction{ActionCode: actionCode}
	upd.AuditorUsers = nextUsers

	// advance group or finish
	if act.AuditorCount > 0 && int32(activeGroup) >= act.AuditorCount {
		upd.AuditorStatus = "CHECKED"
		upd.CurrentAuditorIndex = float64(activeGroup)
	} else {
		upd.AuditorStatus = "NOTCHECKED"
		upd.CurrentAuditorIndex = float64(activeGroup + 1)
	}
	_, err = ca.repo.Update(ctx, actionCode, upd)
	return err
}

func NewCPSActionService(roles storage.CPSActionRoleRepository, repo storage.CPSActionRepository, logger utils.Logger, dispatcher Dispatcher) service.CPSActionService {
	return &cpsActionService{
		repo:       repo,
		logger:     logger,
		roles:      roles,
		dispatcher: dispatcher,
	}
}

func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateCPSAction", "CPSAction", "CreateCPSAction")
	defer span.End()

	roleCode := ctx.Value(constants.ContextKey("role_code")).(string)
	existing, err := ca.GetCPSActionByUniqueID(ctx, cpsAction.RequestAction, roleCode)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		span.AddEvent("failed to get cps action by unique id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if existing != nil {
		span.AddEvent("pending cps action exists", trace.WithAttributes(attribute.String("error", "pending cps action exists")))
		return errors.New(localization.ErrorPendingCpsActionExists.Code)
	}

	cpsAction.RoleCode = roleCode
	err = ca.repo.Save(ctx, cpsAction)
	if err != nil {
		span.AddEvent("failed to save cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (ca *cpsActionService) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveCPSAction", "CPSAction", "ApproveCPSAction")
	defer span.End()

	data, err := ca.repo.Update(ctx, action.ActionCode, *action)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if action.ActionStatus != string(constants.Approved) {
		return nil
	}
	approve, err := ca.dispatcher.Authorize(ctx, data)
	if err != nil && approve == nil {
		span.AddEvent("failed to authorize cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		RollErr := ca.RollBack(ctx, action)
		if err.Error() == localization.ErrorTimeoutError.Code {
			return err
		}
		if RollErr != nil {
			span.AddEvent("failed to roll back cps action", trace.WithAttributes(attribute.String("error", RollErr.Error())))
			return RollErr
		}
		return err
	}
	return nil

}
func (ca *cpsActionService) RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "RejectCPSAction", "CPSAction", "RejectCPSAction")
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}
func (ca *cpsActionService) CancelCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CancelCPSAction", "CPSAction", "CancelCPSAction")
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}
func (ca *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsByDepartment", "CPSAction", "GetCPSActionsByDepartment")
	defer span.End()
	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, department)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetCPSActionsForApprover(ctx context.Context, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForApprover", "CPSAction", "GetCPSActionsForApprover")
	defer span.End()
	result, err := ca.repo.SanitizedFindAllWithPaginationForApprover(ctx, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetCPSActionsForAuditor(ctx context.Context, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForAuditor", "CPSAction", "GetCPSActionsForAuditor")
	defer span.End()
	result, err := ca.repo.SanitizedFindAllWithPaginationForAuditor(ctx, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByID", "CPSAction", "GetCPSActionByID")
	defer span.End()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		span.AddEvent("failed to parse the string to bson object", trace.WithAttributes(attribute.String("error", err.Error())))
		ca.logger.Errorf("their is error when try to parse the string to bson object in service")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	action, err := ca.repo.SanitizedFindOne(ctx, bson.M{"_id": objID, "department": department})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}
func (ca *cpsActionService) GetCPSActionByUniqueID(ctx context.Context, requestAction, role_code string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByUniqueID", "CPSAction", "GetCPSActionByUniqueID")
	defer span.End()

	filter := bson.M{
		"role_code":      role_code,
		"action_status":  string(constants.Pending),
		"request_action": requestAction,
	}

	action, err := ca.repo.SanitizedFindOne(context.Background(), filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ca *cpsActionService) GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByActionCode", "CPSAction", "GetCPSActionByActionCode")
	defer span.End()

	action, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": uniqueID})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ca *cpsActionService) RollBack(ctx context.Context, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "RollBack", "CPSAction", "RollBack")
	defer span.End()
	err := ca.repo.UpdateCustome(ctx, bson.M{"action_code": action.ActionCode}, bson.M{"action_status": string(constants.Pending), "checker_id": "", "checker_name": "", "checker_phone_number": ""})
	if err != nil {
		span.AddEvent("failed to update custom", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
	// return ca.repo.UpdateCustome(ctx, bson.M{"action_code": action.ActionCode}, bson.M{"action_status": string(constants.Pending), "checker_users": []types.Checker{}, "current_checker_index": float32(0)})
}

func (ca *cpsActionService) GetActionCountsByDepartemnt(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetActionCountsByDepartemnt", "CPSAction", "GetActionCountsByDepartemnt")
	defer span.End()
	count, err := ca.repo.GetCountByDepartment(ctx, department)
	if err != nil {
		span.AddEvent("failed to get count by department", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return count, nil
}

func (ca *cpsActionService) GetUserAuthorizerIndex(ctx context.Context, requestAction constants.RequestAction) (imodel.CPSActionApproveIndex, error) {
	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	var approverData imodel.CPSActionApproveIndex

	if mod, ok := ResolveModuleForRA(RequestAction(requestAction)); ok && ca.roles != nil {
		if approver, err := ca.roles.FindApproverByActionName(ctx, strings.ToUpper(mod), roleCode); err == nil {
			approverData = approver
		}
	} else {
		return imodel.CPSActionApproveIndex{}, errors.New(localization.ErrorOperationNotAllowed.Message)
	}

	return approverData, nil
}

func (ca *cpsActionService) GetUserCreatedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCreatedActions", "CPSAction", "GetUserCreatedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID

	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find pending cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCheckedActions", "CPSAction", "GetUserCheckedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID

	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find approver cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}
