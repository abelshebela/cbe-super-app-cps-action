package bps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"strings"

	bpsActionDto "cbe-super-app-cps-action/internal/constants/dto/bps_action"
	"cbe-super-app-cps-action/internal/storage"
	lobal_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bpsActionService struct {
	repo       storage.BPSActionRepository
	roles      storage.BPSActionRoleRepository
	logger     utils.Logger
	dispatcher Dispatcher
}

func NewBPSActionService(roles storage.BPSActionRoleRepository, repo storage.BPSActionRepository, logger utils.Logger, dispatcher Dispatcher) service.BPSActionService {
	return &bpsActionService{
		repo:       repo,
		logger:     logger,
		roles:      roles,
		dispatcher: dispatcher,
	}
}

// IsMakerOnlyForRequest returns true if the module mapped from requestAction is configured as maker-only in CPSActionRole.
func (ba *bpsActionService) IsMakerOnlyForRequest(ctx context.Context, requestAction string) (bool, error) {
	if mod, ok := ResolveModuleForRA(RequestAction(requestAction)); ok && ba.roles != nil {
		role, err := ba.roles.FindByActionName(ctx, mod)
		if err != nil || role == nil {
			return false, errors.New(localization.ErrorOperationNotAllowed.Code)
		}
		return role.IsMakerOnly, nil
	}
	return false, errors.New(localization.ErrorOperationNotAllowed.Code)
}

// AuditorClaim sets auditor status to INPROGRESS when baller belongs to the active group.
func (ba *bpsActionService) AuditorClaim(ctx context.Context, actionCode string, activeGroup int) error {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "AuditorClaim", "CPSAction", "AuditorClaim")
	defer span.End()
	// act, err := ba.repo.SanitizedFindOne(ctx, bson.M{"action_code": actionCode})
	// if err != nil || act == nil {
	// 	return errors.New(localization.ErrorResourceNotFound.Code)
	// }
	// // compute active group from record
	// current := int64(0)
	// if act.Au > 0 {
	// 	current = int64(act.CurrentAuditorIndex)
	// }

	// expected := int64(activeGroup)
	// if current != 0 && current != expected {
	// 	return errors.New(localization.ErrorOperationNotAllowed.Code)
	// }
	// // idempotent move to INPROGRESS
	// upd := bps_model.BPSAction{ActionCode: actionCode}
	// upd.AuditorStatus = "INPROGRESS"
	// _, err = ba.repo.Update(ctx, actionCode, upd)
	return nil
}

// AuditorMark records an auditor's mark and advances to the next group or finishes.
func (ba *bpsActionService) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int) error {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "AuditorMark", "CPSAction", "AuditorMark")
	defer span.End()

	// act, err := ba.repo.SanitizedFindOne(ctx, bson.M{"action_code": actionCode})
	// if err != nil || act == nil {
	// 	ba.logger.Errorf("failed to find action", trace.WithAttributes(attribute.String("error", err.Error())))
	// 	return errors.New(localization.ErrorActionNotFound.Code)
	// }
	// // prevent multiple marks within the same group (any-one quorum)
	// grp := int(activeGroup)
	// for _, au := range act.AuditorUsers {
	// 	if int(au.AuditorIndex) == grp {
	// 		return errors.New(localization.ErrorOperationNotAllowed.Code)
	// 	}
	// }

	// // append this auditor
	// nextUsers := append(act.AuditorUsers, auditor)

	// upd := model.CPSAction{ActionCode: actionCode}
	// upd.AuditorUsers = nextUsers

	// // advance group or finish
	// if act.AuditorCount > 0 && int32(activeGroup) >= act.AuditorCount {
	// 	upd.AuditorStatus = model.AuditorStatus(constants.AUDITORCHECKED)
	// 	upd.CurrentAuditorIndex = float64(activeGroup)
	// } else {
	// 	upd.AuditorStatus = model.AuditorStatus(constants.AUDITORINPROGRESS)
	// 	upd.CurrentAuditorIndex = float64(activeGroup + 1)
	// }
	// _, err = ba.repo.UpdateByActionCode(ctx, actionCode, upd)
	return nil
}

func (ba *bpsActionService) ApproveBPSAction(ctx context.Context, action *bps_model.BPSAction) error {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "ApproveCPSAction", "CPSAction", "ApproveCPSAction")
	defer span.End()

	// data, err := ba.repo.Update(ctx, action.ActionCode, *action)
	// if err != nil {
	// 	span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	// 	return err
	// }

	// if action.ActionStatus != string(constants.Approved) {
	// 	return nil
	// }
	// approve, err := ba.dispatcher.Authorize(ctx, data)
	// if err != nil && approve == nil {
	// 	span.AddEvent("failed to authorize cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	// 	RollErr := ba.RollBack(ctx, action)
	// 	if err.Error() == localization.ErrorTimeoutError.Code {
	// 		return err
	// 	}
	// 	if RollErr != nil {
	// 		span.AddEvent("failed to roll back cps action", trace.WithAttributes(attribute.String("error", RollErr.Error())))
	// 		return RollErr
	// 	}
	// 	return err
	// }
	return nil

}
func (ba *bpsActionService) RejectBPSAction(ctx context.Context, action_code string, action *bps_model.BPSAction) error {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "RejectCPSAction", "CPSAction", "RejectCPSAction")
	defer span.End()
	// _, err := ba.repo.Update(ctx, action_code, *action)
	// if err != nil {
	// 	span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	// 	return err
	// }
	return nil
}
func (ba *bpsActionService) GetBPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionsByDepartment", "CPSAction", "GetCPSActionsByDepartment")
	defer span.End()
	result, err := ba.repo.SanitizedFindAllWithPagination(ctx, *filterParams, department)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ba *bpsActionService) GetBPSActionsForApprover(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionsForApprover", "CPSAction", "GetCPSActionsForApprover")
	defer span.End()
	result, err := ba.repo.SanitizedFindAllWithPaginationForApprover(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ba *bpsActionService) GetBPSActionsForAuditor(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionsForAuditor", "CPSAction", "GetCPSActionsForAuditor")
	defer span.End()
	result, err := ba.repo.SanitizedFindAllWithPaginationForAuditor(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ba *bpsActionService) GetBPSActions(ctx context.Context, userID, role string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActions", "CPSAction", "GetCPSActions")
	defer span.End()

	result, err := ba.repo.SanitizedFindAllWithPaginationCPSActions(ctx, userID, role, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	return result, nil
}

func (ba *bpsActionService) GetBPSActionByID(ctx context.Context, id, department string) (*bps_model.BPSAction, error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionByID", "CPSAction", "GetCPSActionByID")
	defer span.End()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		span.AddEvent("failed to parse the string to bson object", trace.WithAttributes(attribute.String("error", err.Error())))
		ba.logger.Errorf("their is error when try to parse the string to bson object in service")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	action, err := ba.repo.SanitizedFindOne(ctx, bson.M{"_id": objID, "department": department})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}
func (ba *bpsActionService) GetBPSActionByUniqueID(ctx context.Context, requestAction, role_code string) (*bps_model.BPSAction, error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionByUniqueID", "CPSAction", "GetCPSActionByUniqueID")
	defer span.End()

	filter := bson.M{
		"role_code":      role_code,
		"action_status":  string(constants.Pending),
		"request_action": requestAction,
	}

	action, err := ba.repo.SanitizedFindOne(context.Background(), filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ba *bpsActionService) GetBPSActionByForUpdate(ctx context.Context) (*bps_model.BPSAction, error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionByForUpdate", "CPSAction", "GetCPSActionByUniqueID")
	defer span.End()

	var reqs []string
	seen := map[string]struct{}{}

	actionName, _ := ctx.Value(constants.ContextKey("action_name")).(string)

	if lst, ok := RequestActionGroups[actionName]; ok {
		for _, ra := range lst {
			key := string(ra)
			if _, ok := seen[key]; ok {
				continue
			}

			if strings.Contains(key, constants.CREATE) {
				continue
			}
			seen[key] = struct{}{}
			reqs = append(reqs, key)
		}
	}

	filter := bson.M{
		"action_status":  string(constants.Pending),
		"request_action": bson.M{"$in": reqs},
	}

	action, err := ba.repo.SanitizedFindOne(context.Background(), filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ba *bpsActionService) GetBPSActionByActionCode(ctx context.Context, uniqueID, department string) (*bps_model.BPSAction, error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionByActionCode", "CPSAction", "GetCPSActionByActionCode")
	defer span.End()

	action, err := ba.repo.SanitizedFindOne(ctx, bson.M{"action_code": uniqueID})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ba *bpsActionService) RollBack(ctx context.Context, action *bps_model.BPSAction) error {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "RollBack", "CPSAction", "RollBack")
	defer span.End()
	err := ba.repo.UpdateCustome(ctx, bson.M{"action_code": action.ActionCode}, bson.M{"action_status": string(constants.Pending), "checker_id": "", "checker_name": "", "checker_phone_number": ""})
	if err != nil {
		span.AddEvent("failed to update custom", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
	// return ba.repo.UpdateCustome(ctx, bson.M{"action_code": action.ActionCode}, bson.M{"action_status": string(constants.Pending), "checker_users": []types.Checker{}, "current_checker_index": float32(0)})
}

func (ba *bpsActionService) GetActionCountsByDepartemnt(ctx context.Context, department string) (*bpsActionDto.BPSActionCountResponse, error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetActionCountsByDepartemnt", "CPSAction", "GetActionCountsByDepartemnt")
	defer span.End()
	count, err := ba.repo.GetCountByDepartment(ctx, department)
	if err != nil {
		span.AddEvent("failed to get count by department", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return count, nil
}

func (ba *bpsActionService) GetUserAuthorizerIndex(ctx context.Context, requestAction constants.RequestAction) (imodel.BPSActionApproveIndex, error) {
	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	var approverData imodel.BPSActionApproveIndex

	if mod, ok := ResolveModuleForRA(RequestAction(requestAction)); ok && ba.roles != nil {
		if approver, err := ba.roles.FindApproverByActionName(ctx, strings.ToUpper(mod), roleCode); err == nil {
			approverData = approver
		}
	} else {
		return imodel.BPSActionApproveIndex{}, errors.New(localization.ErrorOperationNotAllowed.Message)
	}

	return approverData, nil
}

func (ba *bpsActionService) GetUserCreatedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetUserCreatedActions", "CPSAction", "GetUserCreatedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID

	result, err := ba.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find pending cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ba *bpsActionService) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetUserCheckedActions", "CPSAction", "GetUserCheckedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID

	result, err := ba.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find approver cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}
