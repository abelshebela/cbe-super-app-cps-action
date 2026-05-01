package bps_action

import (
	"cbe-super-app-cps-action/internal/constants"

	"cbe-super-app-cps-action/internal/constants/localization"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"cbe-super-app-cps-action/internal/constants/types"

	mid "cbe-super-app-cps-action/internal/handlers/middleware"

	"cbe-super-app-cps-action/internal/service"

	"strings"

	bpsActionDto "cbe-super-app-cps-action/internal/constants/dto/bps_action"
	customer_dto "cbe-super-app-cps-action/internal/constants/dto/customer"

	"cbe-super-app-cps-action/internal/storage"

	lobal_util "cbe-super-app-cps-action/pkgs/utils"

	"context"

	"errors"

	// bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	bps_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/trace"
)

type bpsActionPublishPayload struct {
	ActionCode    string `json:"action_code"`
	ActionStatus  string `json:"action_status"`
	AuditorStatus string `json:"auditor_status"`

	RoleCode string `json:"role_code"`

	UserData types.UserContext `json:"user_data"`

	Reason string `json:"reason,omitempty"`
}

type bpsActionService struct {
	repo storage.BPSActionRepository

	roles storage.BPSActionRoleRepository

	customerRepo storage.CustomerRepository

	archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository

	bpsUserRepo storage.BPSUserRepository

	cpsUserRepo storage.CpsUserRepository

	logger utils.Logger

	dispatcher Dispatcher
	cfg        config.VaultConfig
}

func NewBPSActionService(roles storage.BPSActionRoleRepository, repo storage.BPSActionRepository, customerRepo storage.CustomerRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, bpsUserRepo storage.BPSUserRepository, cpsUserRepo storage.CpsUserRepository, logger utils.Logger, dispatcher Dispatcher, cfg config.VaultConfig) service.BPSActionService {

	return &bpsActionService{

		repo: repo,

		logger: logger,

		roles: roles,

		customerRepo: customerRepo,

		archivedLinkedAccountRepo: archivedLinkedAccountRepo,

		bpsUserRepo: bpsUserRepo,

		cpsUserRepo: cpsUserRepo,

		dispatcher: dispatcher,
		cfg:        cfg,
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

	producer := mid.GetClientOrchestrationProducer()

	if producer == nil {

		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	rawRoleID, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	if strings.TrimSpace(rawRoleID) == "" {

		return errors.New(localization.ErrorOperationNotAllowed.Code)

	}

	userData, _ := ctx.Value(constants.ContextKey("user_data")).(types.UserContext)

	action, err := ba.GetBPSActionByActionCode(ctx, actionCode, "")

	if err != nil || action == nil {
		return errors.New(localization.ErrorActionNotFound.Code)
	}

	payload := bpsActionPublishPayload{ActionCode: action.ActionCode, ActionStatus: action.Status, RoleCode: rawRoleID, UserData: userData}

	if err := producer.PublishMessage(ctx, payload, "bps.auditor.claim", constants.BPSApproveTopic, "BPS_AUDITOR_CLAIM"); err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil

}

// AuditorMark records an auditor's mark and advances to the next group or finishes.

func (ba *bpsActionService) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int, customerBar bool) error {

	ctx, span := lobal_util.TraceLogger(ctx, "service", "AuditorMark", "CPSAction", "AuditorMark")
	// checkData := lobal_util.ExtractUserFromContext(ctx)
	defer span.End()

	producer := mid.GetClientOrchestrationProducer()

	if producer == nil {

		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	rawRoleID, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	if strings.TrimSpace(rawRoleID) == "" {

		return errors.New(localization.ErrorOperationNotAllowed.Code)

	}

	// userData, _ := ctx.Value(constants.ContextKey("user_data")).(types.UserContext)

	action, err := ba.GetBPSActionByActionCode(ctx, actionCode, "")

	if err != nil || action == nil {

		ba.logger.Errorf("[BpsActionSvc][AuditorMark] action not found: %s", actionCode)

		return errors.New(localization.ErrorActionNotFound.Code)

	}

	auditorApproval := false
	if auditor.AuditorMark == "MARKEDASRIGHT" {
		auditorApproval = true
	}
	err = MarkActionAsAudited(ctx, ba.repo, actionCode, auditorApproval, auditor.AuditorReason, ba.logger)
	if err != nil {
		ba.logger.Errorf("[BPSAction][AuditorMark] failed to updat eh mark")
		return err
	}

	// payload := bpsActionPublishPayload{ActionCode: action.ActionCode, ActionStatus: action.Status, AuditorStatus: string(auditor.AuditorMark), RoleCode: rawRoleID, UserData: checkData, Reason: auditor.AuditorReason}

	// ba.logger.Infof("[BpsActionSvc][AuditorMark] payload: %+v", payload)
	// if err := producer.PublishMessage(ctx, payload, constants.BPSAuditorMarkTopic, ba.cfg.ACIAATMBlockUnBlockUpdateCode1, "BPS_AUDITOR_MARK"); err != nil {

	// 	span.AddEvent("failed to publish bps auditor mark", trace.WithAttributes(attribute.String("error", err.Error())))
	// 	ba.logger.Errorf("[BpsActionSvc][AuditorMark] failed to publish bps auditor mark: %v", err)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)

	// }

	// ba.logger.Infof("[BpsActionSvc][AuditorMark] auditor mark published for action %s, mark=%s", actionCode, auditor.AuditorMark)

	// if customerBar && auditor.AuditorMark == model.MARKEDASWRONG {
	// 	userCode := action.EntityIdentifyer
	// 	if strings.TrimSpace(userCode) == "" {
	// 		ba.logger.Errorf("[BpsActionSvc][AuditorMark] no user_code (entity_identifier) on action: %s", actionCode)
	// 		span.AddEvent("missing entity_identifier for customer bar")
	// 		return errors.New(localization.ErrorInvalidInputParameter.Code)
	// 	}

	// 	if ba.customerRepo != nil {
	// 		if err := ba.customerRepo.BlockCustomerByUserCode(ctx, userCode); err != nil {
	// 			ba.logger.Errorf("[BpsActionSvc][AuditorMark] failed to block customer %s: %v", userCode, err)
	// 			span.RecordError(err)
	// 			return err
	// 		}
	// 		ba.logger.Infof("[BpsActionSvc][AuditorMark] customer %s blocked successfully for action %s", userCode, actionCode)
	// 		span.AddEvent("customer blocked", trace.WithAttributes(attribute.String("user_code", userCode)))
	// 	}
	// }

	return nil

}

func (ba *bpsActionService) ApproveBPSAction(ctx context.Context, action *bps_model.BPSAction) error {

	ctx, span := lobal_util.TraceLogger(ctx, "service", "ApproveBPSAction", "BPSAction", "ApproveBPSAction")

	defer span.End()

	producer := mid.GetClientOrchestrationProducer()

	if producer == nil {

		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	rawRoleID, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	if strings.TrimSpace(rawRoleID) == "" {

		return errors.New(localization.ErrorOperationNotAllowed.Code)

	}

	userData, _ := ctx.Value(constants.ContextKey("user_data")).(types.UserContext)

	payload := bpsActionPublishPayload{ActionCode: action.ActionCode, ActionStatus: action.Status, RoleCode: rawRoleID, UserData: userData}

	ba.logger.Infof("[BPSAction][ApproveBPSAction] payload: %+v", payload)
	if err := producer.PublishMessage(ctx, payload, "bps.approve", constants.BPSApproveTopic, "BPS_APPROVE"); err != nil {

		span.AddEvent("failed to publish bps approve", trace.WithAttributes(attribute.String("error", err.Error())))

		ba.logger.Errorf("[BPSAction][ApproveBPSAction] failed to publish bps approve: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	return nil

}

func (ba *bpsActionService) RejectBPSAction(ctx context.Context, action_code string, action *bps_model.BPSAction) error {

	ctx, span := lobal_util.TraceLogger(ctx, "service", "RejectCPSAction", "CPSAction", "RejectCPSAction")

	defer span.End()

	producer := mid.GetClientOrchestrationProducer()

	if producer == nil {

		ba.logger.Errorf("[BPSAction][RejectBPSAction] client orchestration producer is nil")
		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	rawRoleID, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	if strings.TrimSpace(rawRoleID) == "" {

		return errors.New(localization.ErrorOperationNotAllowed.Code)

	}

	userData, _ := ctx.Value(constants.ContextKey("user_data")).(types.UserContext)

	reason, _ := ctx.Value(constants.ContextKey("rejection_reason")).(string)

	_ = action_code

	payload := bpsActionPublishPayload{ActionCode: action.ActionCode, ActionStatus: action.Status, RoleCode: rawRoleID, UserData: userData, Reason: reason}
	ba.logger.Infof("[BPSAction][RejectBPSAction] payload: %+v", payload)
	if err := producer.PublishMessage(ctx, payload, "bps.reject", constants.BPSApproveTopic, "BPS_REJECT"); err != nil {

		span.AddEvent("failed to publish bps reject", trace.WithAttributes(attribute.String("error", err.Error())))

		ba.logger.Errorf("[BPSAction][RejectBPSAction] failed to publish bps reject: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)

	}

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

	result, err := ba.repo.SanitizedFindAllWithPaginationBPSActions(ctx, userID, role, *filterParams, RAList)

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

		ba.logger.Errorf("[BpsActionSvc][GetByID] parse id err")

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

		"role_code": role_code,

		"status": string(constants.Pending),

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

		"status": string(constants.Pending),

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

// GetBPSActionDetailByActionCode returns the BPS action together with the related
// customer's member detail, currently linked accounts, and unlinked (archived) accounts.
//
// Enrichment is best-effort: any failure on the customer side is logged but does NOT
// fail the whole call so that the action document is still returned to the caller.
// When the action does not target a customer (no UserInformation.UserID) the enrichment
// fields are returned as nil / empty slices.
func (ba *bpsActionService) GetBPSActionDetailByActionCode(ctx context.Context, uniqueID, department string) (*bpsActionDto.BPSActionDetailResponse, error) {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetBPSActionDetailByActionCode", "BPSAction", "GetBPSActionDetailByActionCode")
	defer span.End()

	action, err := ba.GetBPSActionByActionCode(ctx, uniqueID, department)
	if err != nil {
		return nil, err
	}

	resp := &bpsActionDto.BPSActionDetailResponse{
		Action:           action,
		LinkedAccounts:   []customer_dto.LinkedAccount{},
		UnlinkedAccounts: []model.ArchivedLinkedAccount{},
		Checkers:         ba.resolveActionUsers(ctx, action.CheckerID),
		Auditors:         ba.resolveActionUsers(ctx, action.Auditors.AuditorID),
	}

	// The customer module looks customers up by user_code via the Oracle repo
	// (see customer_oracle.FindCustomerDetailByID). Mirror that here: try
	// UserInformation.UserCode first, then fall back to EntityIdentifyer.
	userCode := strings.TrimSpace(action.UserInformation.UserCode)
	if userCode == "" {
		userCode = strings.TrimSpace(action.EntityIdentifyer)
	}
	if userCode == "" {
		span.AddEvent("no user_code on action; skipping customer enrichment")
		return resp, nil
	}

	if ba.customerRepo != nil {
		if detail, derr := ba.customerRepo.FindCustomerDetailByID(ctx, userCode); derr != nil {
			if derr.Error() == localization.ErrorResourceNotFound.Code {
				ba.logger.Infof("[BpsActionSvc][Detail] no member detail for user_code %s", userCode)
			} else {
				ba.logger.Errorf("[BpsActionSvc][Detail] member detail lookup failed for user_code %s: %v", userCode, derr)
				span.AddEvent("member detail lookup failed", trace.WithAttributes(attribute.String("error", derr.Error())))
			}
		} else if detail != nil {
			resp.MemberDetail = detail
			// Oracle returns currently-linked accounts inline on the detail response.
			if detail.LinkedAccount != nil {
				resp.LinkedAccounts = detail.LinkedAccount
			}
		}
	}

	// Unlinked (archived) accounts live in Mongo and are keyed by customer_number (CIF).
	// Use the CIF returned by the Oracle detail; if we don't have it, skip silently.
	if ba.archivedLinkedAccountRepo != nil && resp.MemberDetail != nil {
		cif := strings.TrimSpace(resp.MemberDetail.PersonalInfo.CustomerNumber)
		if cif != "" {
			if unlinked, uerr := ba.archivedLinkedAccountRepo.FindAllByCustomerNumber(ctx, cif); uerr != nil {
				if uerr.Error() == localization.ErrorResourceNotFound.Code {
					ba.logger.Infof("[BpsActionSvc][Detail] no unlinked accounts for cif %s", cif)
				} else {
					ba.logger.Errorf("[BpsActionSvc][Detail] unlinked accounts lookup failed for cif %s: %v", cif, uerr)
					span.AddEvent("unlinked accounts lookup failed", trace.WithAttributes(attribute.String("error", uerr.Error())))
				}
			} else if unlinked != nil {
				resp.UnlinkedAccounts = unlinked
			}
		}
	}

	return resp, nil
}

// resolveActionUsers resolves a slice of user IDs (Mongo ObjectID hex strings) into
// BPSActionUserInfo records. Each ID is looked up in bps_user first; if not found
// there it falls back to cps_user. IDs that resolve nowhere are returned with just
// the ID populated and Source="" so the FE can still render a placeholder row.
//
// Returns an empty slice (not nil) so the JSON payload always has a list.
func (ba *bpsActionService) resolveActionUsers(ctx context.Context, ids []string) []bpsActionDto.BPSActionUserInfo {
	out := make([]bpsActionDto.BPSActionUserInfo, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		out = append(out, ba.resolveOneActionUser(ctx, id))
	}
	return out
}

// resolveOneActionUser tries bps_user.GetByUserID, then cps_user.FindByID.
// Always returns a populated BPSActionUserInfo with at least the ID.
func (ba *bpsActionService) resolveOneActionUser(ctx context.Context, id string) bpsActionDto.BPSActionUserInfo {
	info := bpsActionDto.BPSActionUserInfo{ID: id}

	if ba.bpsUserRepo != nil {
		if u, err := ba.bpsUserRepo.GetByUserID(ctx, id); err == nil && u != nil {
			info.UserCode = u.UserCode
			info.FullName = u.FullName
			info.UserName = u.Username
			info.PhoneNumber = u.PhoneNumber
			info.JobTitle = u.JobTitle
			info.Role = u.Role
			info.BranchCode = u.BranchCode
			info.BranchName = u.BranchName
			info.Source = "bps_user"
			return info
		} else if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			ba.logger.Infof("[BpsActionSvc][Detail] bps_user lookup for %s failed: %v", id, err)
		}
	}

	if ba.cpsUserRepo != nil {
		if u, err := ba.cpsUserRepo.FindByID(ctx, id); err == nil && u != nil {
			info.UserCode = u.UserCode
			info.FullName = u.FullName
			info.UserName = u.UserName
			info.PhoneNumber = u.PhoneNumber
			info.JobTitle = u.JobTitle
			info.Role = u.Role
			info.Source = "cps_user"
			return info
		} else if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			ba.logger.Infof("[BpsActionSvc][Detail] cps_user lookup for %s failed: %v", id, err)
		}
	}

	return info
}

func (ba *bpsActionService) RollBack(ctx context.Context, action *bps_model.BPSAction) error {
	ctx, span := lobal_util.TraceLogger(ctx, "service", "RollBack", "BPSAction", "RollBack")
	defer span.End()

	// Revert to previous checker state by removing the last appended checker
	previousCheckers := action.CheckerID
	if len(previousCheckers) > 0 {
		previousCheckers = previousCheckers[:len(previousCheckers)-1]
	}

	previousApproved := action.CheckersApproved - 1
	if previousApproved < 0 {
		previousApproved = 0
	}

	err := ba.repo.UpdateCustome(ctx,
		bson.M{"action_code": action.ActionCode},
		bson.M{
			"status":            string(constants.Pending),
			"checker_id":        previousCheckers,
			"checkers_approved": previousApproved,
		},
	)
	if err != nil {
		span.AddEvent("failed to roll back bps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	ba.logger.Infof("[BpsActionSvc][RollBack] reverted: %s", action.ActionCode)
	return nil
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

	filterParams.Filters["checker_users.checker_id"] = userID

	result, err := ba.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")

	if err != nil {

		span.AddEvent("failed to find approver cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))

		return nil, err

	}

	return result, nil

}
