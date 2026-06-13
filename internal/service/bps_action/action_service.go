package bps_action

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

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

	local_util "cbe-super-app-cps-action/pkgs/utils"

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
	repo                      storage.BPSActionRepository
	roles                     storage.BPSActionRoleRepository
	customerRepo              storage.CustomerRepository
	actionLogRepo             storage.UserActionLogRepository
	archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository
	bpsUserRepo               storage.BPSUserRepository
	cpsUserRepo               storage.CpsUserRepository
	logger                    utils.Logger

	dispatcher Dispatcher
	cfg        config.VaultConfig
}

func NewBPSActionService(roles storage.BPSActionRoleRepository, repo storage.BPSActionRepository, customerRepo storage.CustomerRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, bpsUserRepo storage.BPSUserRepository, cpsUserRepo storage.CpsUserRepository, actionLogRepo storage.UserActionLogRepository, logger utils.Logger, dispatcher Dispatcher, cfg config.VaultConfig) service.BPSActionService {

	return &bpsActionService{

		repo: repo,

		logger: logger,

		roles: roles,

		customerRepo: customerRepo,

		archivedLinkedAccountRepo: archivedLinkedAccountRepo,

		bpsUserRepo: bpsUserRepo,

		cpsUserRepo: cpsUserRepo,

		actionLogRepo: actionLogRepo,

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

	if ba.actionLogRepo != nil {
		if logErr := ba.actionLogRepo.UpdateAuditorActionStatusByActionCode(ctx, actionCode, string(constants.AUDITORINPROGRESS)); logErr != nil {
			span.AddEvent("failed to update auditor action status on claim", trace.WithAttributes(attribute.String("error", logErr.Error())))
		}
	}

	return nil

}

// AuditorMark records an auditor's mark and advances to the next group or finishes.

func (ba *bpsActionService) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int, customerBar bool) error {
	log := local_util.LoggerFromCtx(ctx, ba.logger)

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

		log.Errorf("[BpsActionSvc][AuditorMark] action not found: %s", actionCode)

		return errors.New(localization.ErrorActionNotFound.Code)

	}

	// Determine auditor approval based on mark
	auditorApproval := false
	isCustomerBarred := false
	markStr := string(auditor.AuditorMark)
	if markStr == string(imodel.MARKEDASRIGHT) {
		auditorApproval = true
	} else if markStr == string(imodel.MARKEDASWRONG) && customerBar {
		// Block customer when marked as wrong with customer_bar=true
		userCode := action.UserInformation.UserCode
		if userCode != "" && ba.customerRepo != nil {
			if blockErr := ba.customerRepo.BlockCustomerByUserCode(ctx, userCode); blockErr != nil {
				log.Errorf("[BpsActionSvc][AuditorMark] failed to block customer %s: %v", userCode, blockErr)
				// Continue with audit even if blocking fails - don't block the audit process
			} else {
				log.Infof("[BpsActionSvc][AuditorMark] customer %s blocked successfully", userCode)
				isCustomerBarred = true
			}
		}
	}
	err = MarkActionAsAudited(ctx, ba.repo, actionCode, auditorApproval, auditor.AuditorReason, ba.logger, isCustomerBarred)
	if err != nil {
		log.Errorf("[BPSAction][AuditorMark] failed to update mark")
		return err
	}

	// BPS has a single auditor — once marked, the audit is complete (CHECKED).
	// Log user action with customer barred flag if applicable
	ba.logUserActionWithCustomerBar(ctx, action, imodel.AUDITOR, action.Status, imodel.AuditorMark(auditor.AuditorMark), "", "", isCustomerBarred)
	if ba.actionLogRepo != nil {
		if logErr := ba.actionLogRepo.AuditorMarkLogsByActionCode(ctx, actionCode, string(auditor.AuditorMark), string(constants.AUDITORCHECKED), isCustomerBarred); logErr != nil {
			span.AddEvent("failed to propagate auditor mark to logs", trace.WithAttributes(attribute.String("error", logErr.Error())))
		}
	}

	return nil

}

func (ba *bpsActionService) ApproveBPSAction(ctx context.Context, action *bps_model.BPSAction) error {
	log := local_util.LoggerFromCtx(ctx, ba.logger)

	ctx, span := lobal_util.TraceLogger(ctx, "service", "ApproveBPSAction", "BPSAction", "ApproveBPSAction")

	defer span.End()

	producer := mid.GetClientOrchestrationProducer()

	if producer == nil {

		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	rawRoleID, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	checkerLevel, _ := ctx.Value(constants.ContextKey("role_checker_group")).(string)

	if strings.TrimSpace(rawRoleID) == "" {

		return errors.New(localization.ErrorOperationNotAllowed.Code)

	}

	userData, _ := ctx.Value(constants.ContextKey("user_data")).(types.UserContext)

	payload := bpsActionPublishPayload{ActionCode: action.ActionCode, ActionStatus: action.Status, RoleCode: rawRoleID, UserData: userData}

	log.Infof("[BPSAction][ApproveBPSAction] payload: %+v", payload)
	if err := producer.PublishMessage(ctx, payload, "bps.approve", constants.BPSApproveTopic, "BPS_APPROVE"); err != nil {

		span.AddEvent("failed to publish bps approve", trace.WithAttributes(attribute.String("error", err.Error())))

		log.Errorf("[BPSAction][ApproveBPSAction] failed to publish bps approve: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	ba.logUserAction(ctx, action, imodel.CHECKER, action.Status, "", checkerLevel, "")
	if ba.actionLogRepo != nil {
		if logErr := ba.actionLogRepo.ApproveUserActionsByActionCode(ctx, action.ActionCode); logErr != nil {
			span.AddEvent("failed to bulk-update action log on approve", trace.WithAttributes(attribute.String("error", logErr.Error())))
		}
	}

	return nil

}

func (ba *bpsActionService) RejectBPSAction(ctx context.Context, action_code string, action *bps_model.BPSAction) error {
	log := local_util.LoggerFromCtx(ctx, ba.logger)

	ctx, span := lobal_util.TraceLogger(ctx, "service", "RejectBPSAction", "BPSAction", "RejectBPSAction")

	defer span.End()

	producer := mid.GetClientOrchestrationProducer()

	if producer == nil {

		log.Errorf("[BPSAction][RejectBPSAction] client orchestration producer is nil")
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
	log.Infof("[BPSAction][RejectBPSAction] payload: %+v", payload)
	if err := producer.PublishMessage(ctx, payload, "bps.reject", constants.BPSApproveTopic, "BPS_REJECT"); err != nil {

		span.AddEvent("failed to publish bps reject", trace.WithAttributes(attribute.String("error", err.Error())))

		log.Errorf("[BPSAction][RejectBPSAction] failed to publish bps reject: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)

	}

	ba.logUserAction(ctx, action, imodel.CHECKER, action.Status, "", "", "")
	if ba.actionLogRepo != nil {
		if logErr := ba.actionLogRepo.RejectUserActionsByActionCode(ctx, action.ActionCode); logErr != nil {
			span.AddEvent("failed to bulk-update action log on reject", trace.WithAttributes(attribute.String("error", logErr.Error())))
		}
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
	log := local_util.LoggerFromCtx(ctx, ba.logger)

	log.Infof("[BpsActionSvc][GetBPSActionsForApprover] Retrieving BPS actions for user %s with filter %+v", userID, filterParams.Filters)
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}

	// Detect status filter to route approved/rejected through the user action log.
	statusFilter := ""
	if v, ok := filterParams.Filters["status"]; ok {
		if s, ok := v.(string); ok {
			statusFilter = strings.ToUpper(strings.TrimSpace(s))
		}
	}

	if (statusFilter == "APPROVED" || statusFilter == "REJECTED") && ba.actionLogRepo != nil {
		logFilter := map[string]interface{}{
			"action_type":         string(bps_model.BPSActions),
			"given_action_status": statusFilter,
			// "user_action_responsibilities": string(bps_model.CHECKER),
			"username": userID,
		}

		actionCodes, err := ba.actionLogRepo.GetActionCodesByFilter(ctx, logFilter)
		if err != nil {
			span.AddEvent("failed to get action codes from log", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}

		delete(filterParams.Filters, "status")

		if len(actionCodes) == 0 {
			return &types.PaginatedResponse[[]*bps_model.BPSAction]{
				Data: []*bps_model.BPSAction{},
				Meta: lobal_util.BuildPaginationMeta(0, filterParams.Page, filterParams.PerPage),
			}, nil
		}
		filterParams.Filters["action_code"] = bson.M{"$in": actionCodes}
	}

	log.Infof("[BpsActionSvc][GetBPSActionsForApprover] final filters for repo query: %+v", filterParams.Filters)

	result, err := ba.repo.SanitizedFindAllWithPaginationForApprover(ctx, userID, *filterParams, RAList)

	if err != nil {

		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))

		return nil, err

	}

	log.Infof("[BpsActionSvc][GetBPSActionsForApprover] ----------------found %d actions for user %s with filter %+v", len(result.Data), userID, filterParams.Filters)
	return result, nil

}

func (ba *bpsActionService) GetBPSActionsForAuditor(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*bps_model.BPSAction], error) {

	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetBPSActionsForAuditor", "BPSAction", "GetBPSActionsForAuditor")

	defer span.End()

	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}

	// Extract filters
	levels := extractStringSlice(filterParams.Filters, "levels")
	services := extractStringSlice(filterParams.Filters, "services")

	// Check for customer_bared filter
	var auditorCustomerBared *bool
	if v, ok := filterParams.Filters["customer_bared"]; ok {
		if b, ok := v.(bool); ok {
			auditorCustomerBared = &b
		}
	}

	// NOTCHECKED: BPS tracks un-audited state via auditors.audited on the document itself.
	// No user_action_log entry exists yet for these actions, so bypass the log lookup entirely
	// and query the BPS repo directly.
	auditorStatus, _ := filterParams.Filters["auditor_status"].(string)
	customerBarred, _ := filterParams.Filters["customer_bared"].(bool)

	if strings.ToUpper(auditorStatus) == string(constants.AUDITORNOTCHECKED) && !customerBarred {
		result, err := ba.repo.SanitizedFindAllWithPaginationForAuditor(ctx, userID, *filterParams, RAList)
		if err != nil {
			span.AddEvent("failed to find all with pagination for auditor (NOTCHECKED)", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return result, nil
	}

	// Build log filter
	logFilter := bps_model.UserActionLogActionCodeFilter{
		RequestActions:       RAList,
		Levels:               levels,
		Services:             services,
		Responsibilities:     []string{string(bps_model.AUDITOR)},
		AuditorCustomerBared: auditorCustomerBared,
	}

	// level+claim filtering: all pairs must be satisfied within the same action_code.
	if levelClaimPairs := bpsExtractLevelClaimPairs(filterParams.Filters); len(levelClaimPairs) > 0 {
		logFilter.LevelClaimPairs = levelClaimPairs
	}

	// Get action codes from user_action_log
	if ba.actionLogRepo != nil {
		actionCodes, err := ba.actionLogRepo.GetActionCodesByActionLogFilter(ctx, logFilter)
		if err != nil {
			span.AddEvent("failed to get action codes by action log filter", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		if len(actionCodes) == 0 {
			return &types.PaginatedResponse[[]*bps_model.BPSAction]{
				Data: []*bps_model.BPSAction{},
				Meta: lobal_util.BuildPaginationMeta(0, filterParams.Page, filterParams.PerPage),
			}, nil
		}
		filterParams.Filters["action_code"] = bson.M{"$in": actionCodes}
	}

	result, err := ba.repo.SanitizedFindAllWithPaginationForAuditor(ctx, userID, *filterParams, RAList)

	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	return result, nil
}

func bpsExtractLevelClaimPairs(filters map[string]interface{}) []bps_model.LevelClaimPair {
	v, ok := filters["level_claim_pairs"]
	if !ok {
		return nil
	}
	if pairs, ok := v.([]bps_model.LevelClaimPair); ok {
		return pairs
	}
	return nil
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
	log := local_util.LoggerFromCtx(ctx, ba.logger)

	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetCPSActionByID", "CPSAction", "GetCPSActionByID")

	defer span.End()

	objID, err := bson.ObjectIDFromHex(id)

	if err != nil {

		span.AddEvent("failed to parse the string to bson object", trace.WithAttributes(attribute.String("error", err.Error())))

		log.Errorf("[BpsActionSvc][GetByID] parse id err")

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
	log := local_util.LoggerFromCtx(ctx, ba.logger)

	ctx, span := lobal_util.TraceLogger(ctx, "service", "GetBPSActionDetailByActionCode", "BPSAction", "GetBPSActionDetailByActionCode")
	defer span.End()

	action, err := ba.GetBPSActionByActionCode(ctx, uniqueID, department)
	if err != nil {
		return nil, err
	}

	maker := ba.resolveOneActionUser(ctx, action.MakerID)

	rejectionReason := ""
	if ActionStatus(action.Status) == ActionRejected {
		rejectionReason = action.ActionReason.ActionNote
	}

	resp := &bpsActionDto.BPSActionDetailResponse{
		Action:                 action,
		LinkedAccounts:         []customer_dto.LinkedAccount{},
		UnlinkedAccounts:       []model.ArchivedLinkedAccount{},
		Maker:                  &maker,
		Checkers:               ba.resolveActionUsers(ctx, action.CheckerID),
		Auditors:               ba.resolveActionUsers(ctx, action.Auditors.AuditorID),
		CheckerLevels:          buildCheckerLevels(action),
		RejectionReason:        rejectionReason,
		AuditorRejectionReason: action.Auditors.Reason,
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
				log.Infof("[BpsActionSvc][Detail] no member detail for user_code %s", userCode)
			} else {
				log.Errorf("[BpsActionSvc][Detail] member detail lookup failed for user_code %s: %v", userCode, derr)
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
					log.Infof("[BpsActionSvc][Detail] no unlinked accounts for cif %s", cif)
				} else {
					log.Errorf("[BpsActionSvc][Detail] unlinked accounts lookup failed for cif %s: %v", cif, uerr)
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
	log := local_util.LoggerFromCtx(ctx, ba.logger)

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
			log.Infof("[BpsActionSvc][Detail] bps_user lookup for %s failed: %v", id, err)
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
			info.Department = u.Department.Hex()
			info.Source = "cps_user"
			return info
		} else if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			log.Infof("[BpsActionSvc][Detail] cps_user lookup for %s failed: %v", id, err)
		}
	}

	return info
}

// buildCheckerLevels derives per-slot approval status from the flat CheckerID
// and CheckerTime slices stored on the action.
//
//   - Slots below CheckersApproved are APPROVED (with their timestamp).
//   - When the overall status is REJECTED the last populated slot is REJECTED.
//   - All remaining slots beyond the populated ones are PENDING.
func buildCheckerLevels(action *bps_model.BPSAction) []bpsActionDto.CheckerLevelInfo {
	needed := action.CheckersNeeded
	if needed <= 0 {
		needed = len(action.CheckerID)
	}
	if needed == 0 {
		return []bpsActionDto.CheckerLevelInfo{}
	}

	levels := make([]bpsActionDto.CheckerLevelInfo, needed)
	for i := 0; i < needed; i++ {
		level := bpsActionDto.CheckerLevelInfo{
			Level:  i + 1,
			Status: "PENDING",
		}
		if i < len(action.CheckerID) {
			level.CheckerID = action.CheckerID[i]
		}
		if i < action.CheckersApproved {
			level.Status = "APPROVED"
			if i < len(action.CheckerTime) {
				t := action.CheckerTime[i]
				level.ApprovedAt = &t
			}
		}
		levels[i] = level
	}

	// If overall status is REJECTED, the last populated checker slot is the one
	// that triggered rejection — mark it REJECTED instead of PENDING.
	if string(action.Status) == string(ActionRejected) && len(action.CheckerID) > 0 {
		lastIdx := len(action.CheckerID) - 1
		if lastIdx >= action.CheckersApproved && lastIdx < needed {
			var ts *time.Time
			if lastIdx < len(action.CheckerTime) {
				t := action.CheckerTime[lastIdx]
				ts = &t
			}
			levels[lastIdx].Status = "REJECTED"
			levels[lastIdx].ApprovedAt = ts
		}
	}

	return levels
}

func (ba *bpsActionService) RollBack(ctx context.Context, action *bps_model.BPSAction) error {
	log := local_util.LoggerFromCtx(ctx, ba.logger)

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
	log.Infof("[BpsActionSvc][RollBack] reverted: %s", action.ActionCode)
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

func (ba *bpsActionService) logUserAction(ctx context.Context, action *bps_model.BPSAction, responsibility imodel.UserActionResponsibility, givenStatus string, auditorMark imodel.AuditorMark, checkerLevel, auditorLevel string) {
	reqLog := local_util.LoggerFromCtx(ctx, ba.logger)
	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	userData := local_util.ExtractUserFromContext(ctx)
	if ba.actionLogRepo == nil {
		return
	}

	// userData, _ := ctx.Value(constants.ContextKey("user_data")).(types.UserContext)

	userOID, _ := bson.ObjectIDFromHex(userData.UserID)
	actionOID := action.ID

	serviceName := ""
	if mod, ok := ResolveModuleForRA(RequestAction(action.RequestAction)); ok {
		serviceName = mod
	}

	actionLog := &imodel.UserActionLog{
		ID:                         bson.NewObjectID(),
		ActionID:                   actionOID,
		ActionCode:                 action.ActionCode,
		GivenActionStatus:          givenStatus,
		GivenAuditorStatus:         auditorMark,
		RequestAction:              constants.RequestAction(action.RequestAction),
		ActionTakenServiceName:     serviceName,
		CheckerLevel:               checkerLevel,
		AuditorLevel:               auditorLevel,
		UserRoleCode:               roleCode,
		UserID:                     userOID,
		Username:                   userData.UserName,
		UserPhone:                  userData.PhoneNumber,
		ActionType:                 imodel.BPSActions,
		UserActionResponsibilities: responsibility,
		CreatedAt:                  time.Now(),
	}

	if err := ba.actionLogRepo.Upsert(ctx, actionLog); err != nil {
		reqLog.Errorf("[BpsActionSvc][logUserAction] failed to log action: %v", err)
	}
}

func (ba *bpsActionService) logUserActionWithCustomerBar(ctx context.Context, action *bps_model.BPSAction, responsibility imodel.UserActionResponsibility, givenStatus string, auditorMark imodel.AuditorMark, checkerLevel, auditorLevel string, auditorCustomerBared bool) {
	reqLog := local_util.LoggerFromCtx(ctx, ba.logger)
	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	userData := local_util.ExtractUserFromContext(ctx)
	if ba.actionLogRepo == nil {
		return
	}

	userOID, _ := bson.ObjectIDFromHex(userData.UserID)
	actionOID := action.ID

	serviceName := ""
	if mod, ok := ResolveModuleForRA(RequestAction(action.RequestAction)); ok {
		serviceName = mod
	}

	actionLog := &imodel.UserActionLog{
		ID:                         bson.NewObjectID(),
		ActionID:                   actionOID,
		ActionCode:                 action.ActionCode,
		GivenActionStatus:          givenStatus,
		GivenAuditorStatus:         auditorMark,
		RequestAction:              constants.RequestAction(action.RequestAction),
		ActionTakenServiceName:     serviceName,
		CheckerLevel:               checkerLevel,
		AuditorLevel:               auditorLevel,
		UserRoleCode:               roleCode,
		UserID:                     userOID,
		Username:                   userData.UserName,
		UserPhone:                  userData.PhoneNumber,
		ActionType:                 imodel.BPSActions,
		UserActionResponsibilities: responsibility,
		AuditorCustomerBared:       auditorCustomerBared,
		CreatedAt:                  time.Now(),
	}

	if err := ba.actionLogRepo.Upsert(ctx, actionLog); err != nil {
		reqLog.Errorf("[BpsActionSvc][logUserActionWithCustomerBar] failed to log action: %v", err)
	}
}

// extractStringSlice safely extracts a string slice from filters map.
// Handles nil map, single string, and []interface{} (common after JSON unmarshalling).
func extractStringSlice(filters map[string]interface{}, key string) []string {
	if filters == nil {
		return nil
	}
	v, ok := filters[key]
	if !ok {
		return nil
	}
	// Direct string slice
	if ss, ok := v.([]string); ok {
		return ss
	}
	// Single string
	if s, ok := v.(string); ok && s != "" {
		return []string{s}
	}
	// []interface{} after JSON unmarshalling
	if arr, ok := v.([]interface{}); ok {
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

// ReinstateCustomer reinstates/unblocks a customer that was barred by auditor
func (ba *bpsActionService) ReinstateCustomer(ctx context.Context, actionCode, userCode, reason string) error {
	log := local_util.LoggerFromCtx(ctx, ba.logger)

	ctx, span := lobal_util.TraceLogger(ctx, "service", "ReinstateCustomer", "BPSAction", "ReinstateCustomer")
	defer span.End()

	// 1. Unblock the customer
	if userCode != "" && ba.customerRepo != nil {
		if err := ba.customerRepo.UNBlockCustomerByUserCode(ctx, userCode); err != nil {
			log.Errorf("[BpsActionSvc][ReinstateCustomer] failed to unblock customer %s: %v", userCode, err)
			return err
		}
		log.Infof("[BpsActionSvc][ReinstateCustomer] customer %s unblocked successfully", userCode)
	}

	// 2. Update user_action_log entries for this action_code
	// Set IsCustomerReinstated = true and ReinstateReason
	if ba.actionLogRepo != nil && actionCode != "" {
		if err := ba.actionLogRepo.UpdateReinstateStatus(ctx, actionCode, reason); err != nil {
			log.Errorf("[BpsActionSvc][ReinstateCustomer] failed to update user_action_log for action %s: %v", actionCode, err)
			// Don't return error - continue even if log update fails
		} else {
			log.Infof("[BpsActionSvc][ReinstateCustomer] updated user_action_log for action %s", actionCode)
		}
	}

	return nil
}
