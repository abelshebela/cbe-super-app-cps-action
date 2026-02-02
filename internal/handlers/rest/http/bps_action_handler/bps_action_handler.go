package bps_action_handler

import (
	"cbe-super-app-cps-action/internal/constants"
	bpsactionDto "cbe-super-app-cps-action/internal/constants/dto/bps_action"
	bps_actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/bps_action_role"
	bpsaction "cbe-super-app-cps-action/internal/constants/interfaces/bps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	mid "cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/service"
	bpsactionsvc "cbe-super-app-cps-action/internal/service/bps_action"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"

	"strings"
	"time"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type bps_action_resp *bps_model.BPSAction
type bps_action_dto_Resp *bpsactionDto.ActionRequest
type bps_actions_paginated_resp *types.PaginatedResponse[[]*bps_model.BPSAction]

type bpsActionAdapter struct {
	bpsActionApplication service.BPSActionService
	logger               utils.Logger
}

func InitBPSActionAdapter(bpsActionApplication service.BPSActionService, logger utils.Logger) bpsaction.BPSActionAdapter {
	return &bpsActionAdapter{
		logger:               logger,
		bpsActionApplication: bpsActionApplication,
	}
}

// AuditorAction handles both claim (start) and mark (complete) for auditor in a single endpoint
// POST /actions/{action_code}/auditor
func (a *bpsActionAdapter) AuditorAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "auditorAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	// Load action (to resolve module/request_action)
	action, err := a.bpsActionApplication.GetBPSActionByActionCode(ctx, actionCode, "")
	if err != nil || action == nil {
		span.RecordError(err)
		if err.Error() == localization.ErrorActionNotFound.Code {
			localization.SendErrorByCodeResponse(w, localization.ErrorActionDataNotFound.Code)
			return
		}
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	// Determine caller
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}
	if strings.Contains(strings.ToUpper(userData.UserRole), constants.Maker) {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(rawRoleID) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	// Parse request to decide claim vs mark
	var reqBody bps_actionrole_dto.AuditorMarkRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil && err != io.EOF {
			localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
			return
		}
	}

	if action.MakerID == userData.UserID {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	// Publish via AuditorMark (claim is not used)
	auditor := model.Auditor{
		AuditorID:          userData.UserID,
		RoleID:             rawRoleID,
		AuditorIndex:       0,
		AuditorName:        userData.FullName,
		AuditorPhoneNumber: userData.PhoneNumber,
		AuditorReason:      reqBody.Reason,
		AuditorMark:        model.AuditorMark(strings.ToUpper(strings.TrimSpace(reqBody.Mark))),
		ApprovedAt:         time.Now(),
	}
	ctx = context.WithValue(ctx, constants.ContextKey("user_data"), userData)
	if err := a.bpsActionApplication.AuditorMark(ctx, actionCode, auditor, 0); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSActionChecked, nil)
}

// ApproveCPSAction approves a CPS action
//
//	@Summary		Approve CPS action
//	@Description	Approves a CPS action by action code
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path		string									true	"Action Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS action approved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid action code"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/{action_code}/approve [patch]
func (a *bpsActionAdapter) ApproveBPSAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "approveCpsAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(rawRoleID) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUserForbidden.Message)
		return
	}
	if strings.Contains(strings.ToUpper(userData.UserRole), constants.Maker) {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	span.SetAttributes(attribute.String("cps_action.code", actionCode))
	action, err := a.bpsActionApplication.GetBPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	ctx = context.WithValue(ctx, constants.ContextKey("user_data"), userData)
	if err := a.bpsActionApplication.ApproveBPSAction(ctx, action); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBPSActionAuthorized, nil)
}

// RejectCPSAction rejects a CPS action
//
//	@Summary		Reject CPS action
//	@Description	Rejects a CPS action by action code with rejection reason
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path		string									true	"Action Code"
//	@Param			request		body		bps_action_dto_Resp						true	"Rejection request"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS action rejected successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input or missing rejection reason"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/{action_code}/reject [patch]
func (a *bpsActionAdapter) RejectBPSAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "rejectCpsAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))
	var req bpsactionDto.ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(rawRoleID) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUserForbidden.Message)
		return
	}

	// Fetch action
	action, err := a.bpsActionApplication.GetBPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	ctx = context.WithValue(ctx, constants.ContextKey("user_data"), userData)
	ctx = context.WithValue(ctx, constants.ContextKey("rejection_reason"), req.RejectionReason)
	if err := a.bpsActionApplication.RejectBPSAction(ctx, actionCode, action); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBPSActionRejected, nil)
}

// GetCPSActionsByDepartment retrieves CPS actions by department
//
//	@Summary		Get CPS actions by department
//	@Description	Retrieves a paginated list of CPS actions for the user's department
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"		default(1)
//	@Param			per_page	query		int																false	"Items per page"	default(10)
//	@Param			search		query		string															false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=bps_actions_paginated_resp}	"CPS actions retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/ [get]
func (a *bpsActionAdapter) GetBPSActionsByDepartment(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionsByDepartment", "handler", "cpsAction")
	defer span.End()

	actions, err := a.bpsActionApplication.GetBPSActionsByDepartment(ctx, userData.Department, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("cps_action.department", userData.Department))
	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, actions)
}

func (a *bpsActionAdapter) GetUserCheckedActions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserApprovedCpsActions", "handler", "cpsAction")
	defer span.End()
	userData := local_util.ExtractUserContext(r)
	userID := userData.UserID
	res, err := a.bpsActionApplication.GetUserCheckedActions(ctx, userID, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, res)
}

// GetCPSActionByID retrieves a CPS action by ID
//
//	@Summary		Get CPS action by ID
//	@Description	Retrieves a specific CPS action by its ID
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_id	path		string												true	"Action ID"
//	@Success		200			{object}	localization.StandardResponse{data=bps_action_resp}	"CPS action retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}				"Bad request - Invalid action ID"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}				"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/by-id/{action_id} [get]
func (a *bpsActionAdapter) GetBPSActionByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionById", "handler", "cpsAction")
	defer span.End()
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}
	if strings.Contains(strings.ToUpper(userData.UserRole), constants.Maker) {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	actionID := chi.URLParam(r, string(constants.ActionID))

	span.SetAttributes(
		attribute.String("cps_action.id", actionID),
		attribute.String("cps_action.department", userData.Department),
	)

	action, err := a.bpsActionApplication.GetBPSActionByID(ctx, actionID, userData.Department)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBPSActionFetched, action)
}

// GetCPSActionByActionCode retrieves a BPS action by action code with history
//
//	@Summary		Get CPS action by action code with history
//	@Description	Retrieves a specific CPS action by its action code, including previous and current actions
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path		string												whitespace	true	"Action Code"
//	@Success		200			{object}	localization.StandardResponse{data=bps_action_resp}	"CPS action with history retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}				"Bad request - Invalid action code"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}				"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/by-action-code/{action_code} [get]
func (a *bpsActionAdapter) GetBPSActionByActionCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionByCode", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}
	if strings.Contains(strings.ToUpper(userData.UserRole), constants.Maker) {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	span.SetAttributes(
		attribute.String("cps_action.code", actionCode),
		attribute.String("cps_action.department", userData.Department),
	)

	action, err := a.bpsActionApplication.GetBPSActionByActionCode(ctx, actionCode, userData.Department)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBPSActionFetched, action)
}

func (a *bpsActionAdapter) GetUserApproverActions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	userID := local_util.ExtractUserContext(r).UserID

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserApproverPendingActions", "handler", "cpsAction")
	defer span.End()

	// roleCode from context
	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if rawRoleID == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	// fetch checker allocations for this role
	idxRepo := mid.GetBPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}
	_, checkerActions, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if checkerActions == nil {
		localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, map[string]interface{}{})
		return
	}

	a.logger.Infof("Checker Actions: %v", checkerActions)

	// resolve action_names -> request_actions
	var reqs []string
	seen := map[string]struct{}{}
	for _, mod := range checkerActions {
		upper := strings.ToUpper(strings.TrimSpace(mod))
		if lst, ok := bpsactionsvc.RequestActionGroups[upper]; ok {
			for _, ra := range lst {
				key := string(ra)
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				reqs = append(reqs, key)
			}
		}
	}
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	if len(reqs) > 0 {
		filterParams.Filters["request_action"] = map[string]interface{}{"$in": reqs}
	}

	res, err := a.bpsActionApplication.GetBPSActionsForApprover(ctx, userID, reqs, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, res)
}

func (a *bpsActionAdapter) GetUserAuditorActions(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserApproverPendingActions", "handler", "cpsAction")
	defer span.End()

	filterParams := local_util.ExtractFilterParams(r)
	var allocation []string
	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	// roleCode from context
	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if rawRoleID == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	// fetch checker allocations for this role
	idxRepo := mid.GetBPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	_, _, auditorAllocations, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if auditorAllocations == nil {
		localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, map[string]interface{}{})
		return
	}

	for _, v := range auditorAllocations {
		if slices.Contains(allocation, v) {
			continue
		}
		allocation = append(allocation, v)
	}
	// resolve action_names -> request_actions
	var reqs []string
	seen := map[string]struct{}{}
	for _, mod := range allocation {
		upper := strings.ToUpper(strings.TrimSpace(mod))
		if lst, ok := bpsactionsvc.RequestActionGroups[upper]; ok {
			for _, ra := range lst {
				key := string(ra)
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				reqs = append(reqs, key)
			}
		}
	}

	// do not force action_status; let API-provided filters decide
	userID := local_util.ExtractUserContext(r).UserID
	res, err := a.bpsActionApplication.GetBPSActionsForAuditor(ctx, userID, reqs, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, res)
}

func (a *bpsActionAdapter) GetUserApproverApprovedActions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserApproverApprovedActions", "handler", "cpsAction")
	defer span.End()

	// roleCode from context
	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if rawRoleID == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	idxRepo := mid.GetBPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}
	_, checkerActions, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	var reqs []string
	seen := map[string]struct{}{}
	for _, mod := range checkerActions {
		upper := strings.ToUpper(strings.TrimSpace(mod))
		if lst, ok := bpsactionsvc.RequestActionGroups[upper]; ok {
			for _, ra := range lst {
				key := string(ra)
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				reqs = append(reqs, key)
			}
		}
	}
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	if len(reqs) > 0 {
		filterParams.Filters["request_action"] = map[string]interface{}{"$in": reqs}
	}
	// do not force action_status; let API-provided filters decide

	userID := local_util.ExtractUserContext(r).UserID
	res, err := a.bpsActionApplication.GetBPSActionsForApprover(ctx, userID, reqs, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, res)
}

func (a *bpsActionAdapter) GetActionCounts(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionCounts", "handler", "cpsAction")
	defer span.End()
	if _, err := local_util.ParseUserContext(r); err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	userID := userContext.UserID

	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(rawRoleID) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	requestedRole := r.URL.Query().Get("role")
	if err := local_util.NoSpecialChars(requestedRole); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	idxRepo := mid.GetBPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	makerActions, checkerActions, auditorActions, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if makerActions == nil && checkerActions == nil && auditorActions == nil {
		resp := &bpsactionDto.BPSActionCountResponse{
			Approved:   0,
			Rejected:   0,
			Inprogress: 0,
			Completed:  0,
		}
		localization.SendSuccessResponse(w, localization.SuccessBPSActionCount, resp)
		return
	}

	// resolve action_names -> request_actions (same as GetUserApproverActions)
	var reqs []string
	seen := map[string]struct{}{}

	if makerActions != nil && requestedRole == "maker" {
		for _, mod := range makerActions {
			upper := strings.ToUpper(strings.TrimSpace(mod))
			if lst, ok := bpsactionsvc.RequestActionGroups[upper]; ok {
				for _, ra := range lst {
					key := string(ra)
					if _, ok := seen[key]; ok {
						continue
					}
					seen[key] = struct{}{}
					reqs = append(reqs, key)
				}
			}
		}

	}

	if checkerActions != nil && requestedRole == "checker" {
		for _, mod := range checkerActions {
			upper := strings.ToUpper(strings.TrimSpace(mod))
			if lst, ok := bpsactionsvc.RequestActionGroups[upper]; ok {
				for _, ra := range lst {
					key := string(ra)
					if _, ok := seen[key]; ok {
						continue
					}
					seen[key] = struct{}{}
					reqs = append(reqs, key)
				}
			}
		}
	}

	if auditorActions != nil && requestedRole == "auditor" {
		for _, mod := range auditorActions {
			upper := strings.ToUpper(strings.TrimSpace(mod))
			if lst, ok := bpsactionsvc.RequestActionGroups[upper]; ok {
				for _, ra := range lst {
					key := string(ra)
					if _, ok := seen[key]; ok {
						continue
					}
					seen[key] = struct{}{}
					reqs = append(reqs, key)
				}
			}
		}
	}

	// If no mapped request actions, return zero counts
	if len(reqs) == 0 {
		resp := &bpsactionDto.BPSActionCountResponse{
			Approved:   0,
			Rejected:   0,
			Inprogress: 0,
			Completed:  0,
		}
		localization.SendSuccessResponse(w, localization.SuccessBPSActionCount, resp)
		return
	}

	// Helper to compute counts per status using pagination meta totals
	buildFilter := func(status, auditorStatus string) *types.Filter {
		f := &types.Filter{
			Page:    constants.DefaultPage,
			PerPage: 1,
			Search:  "",
			Filters: map[string]interface{}{},
		}
		f.Filters[string(constants.ActionStatus)] = status

		if auditorStatus != "" {
			f.Filters[string(constants.AuditorStatusDBFieldName)] = auditorStatus
		}
		return f
	}

	var approvedCount, rejectedCount, inprogressAuditCount, completedAuditCount int

	// Approved
	if auditorActions != nil && requestedRole == "auditor" {
		if res, err := a.bpsActionApplication.GetBPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Approved), string(model.AUDITORNOTCHECKED))); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			approvedCount = int(res.Meta.TotalDocs)
		}
	} else {
		if res, err := a.bpsActionApplication.GetBPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Approved), "")); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			approvedCount = int(res.Meta.TotalDocs)
		}

	}

	// Rejected
	if auditorActions != nil && requestedRole == "auditor" {
		if res, err := a.bpsActionApplication.GetBPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Rejected), string(model.AUDITORNOTCHECKED))); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			rejectedCount = int(res.Meta.TotalDocs)
		}
	} else {
		if res, err := a.bpsActionApplication.GetBPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Rejected), "")); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			rejectedCount = int(res.Meta.TotalDocs)
		}
	}

	// Audior's Inprogress Count
	if res, err := a.bpsActionApplication.GetBPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(model.AUDITORINPROGRESS), "")); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	} else if res != nil && res.Meta.TotalDocs > 0 {
		inprogressAuditCount = int(res.Meta.TotalDocs)
	}

	// Audior's Completed Count
	if res, err := a.bpsActionApplication.GetBPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(model.AUDITORNOTCHECKED), "")); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	} else if res != nil && res.Meta.TotalDocs > 0 {
		completedAuditCount = int(res.Meta.TotalDocs)
	}

	resp := &bpsactionDto.BPSActionCountResponse{
		Approved:   approvedCount,
		Rejected:   rejectedCount,
		Inprogress: inprogressAuditCount,
		Completed:  completedAuditCount,
	}

	localization.SendSuccessResponse(w, localization.SuccessBPSActionCount, resp)
}

func (a *bpsActionAdapter) GetAuthorizerIndex(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserAuthorizerIndex", "handler", "cpsAction")
	defer span.End()
	requestAction := chi.URLParam(r, "request_action")

	if requestAction == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	authorizerIndex, err := a.bpsActionApplication.GetUserAuthorizerIndex(ctx, constants.RequestAction(requestAction))
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[CPSAction.GetActionCounts] service failed %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, authorizerIndex)

}

// ApproverCheckerAllocations returns checker allocations for the caller's role,
// grouped by module with request_actions and derived action_types.
func (a *bpsActionAdapter) ApproverCheckerAllocations(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "approverCheckerAllocations", "handler", "cpsAction")
	defer span.End()
	roleCode, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(roleCode) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}
	repo := mid.GetBPSActionApproveRepo()
	if repo == nil {
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}
	// maker, checker, auditor, portalCards
	_, checkerMods, _, err := repo.PopulateUserApproverAllocations(ctx, roleCode)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	type modInfo struct {
		RequestActions []string `json:"request_actions"`
		ActionTypes    []string `json:"action_types"`
	}
	deriveTypes := func(actions []string) []string {
		seen := map[string]struct{}{}
		for _, ra := range actions {
			u := strings.ToUpper(ra)
			switch {
			case strings.Contains(u, "CREATE"):
				seen["CREATE"] = struct{}{}
			case strings.Contains(u, "UPDATE"):
				seen["UPDATE"] = struct{}{}
			case strings.Contains(u, "DELETE"):
				seen["DELETE"] = struct{}{}
			case strings.Contains(u, "ENABLE"):
				seen["ENABLE"] = struct{}{}
			case strings.Contains(u, "DISABLE"):
				seen["DISABLE"] = struct{}{}
			}
		}
		out := make([]string, 0, len(seen))
		for k := range seen {
			out = append(out, k)
		}
		return out
	}
	build := func(mods []string) (map[string]modInfo, []string) {
		m := map[string]modInfo{}
		uniq := map[string]struct{}{}
		for _, raw := range mods {
			mod := strings.ToUpper(strings.TrimSpace(raw))
			if mod == "" {
				continue
			}
			uniq[mod] = struct{}{}
			var reqs []string
			if group, ok := bpsactionsvc.RequestActionGroups[mod]; ok {
				reqs = make([]string, 0, len(group))
				for _, ga := range group {
					reqs = append(reqs, string(ga))
				}
			} else {
				reqs = []string{}
			}
			m[mod] = modInfo{RequestActions: reqs, ActionTypes: deriveTypes(reqs)}
		}
		list := make([]string, 0, len(uniq))
		for k := range uniq {
			list = append(list, k)
		}
		return m, list
	}
	checkerMap, checkerList := build(checkerMods)
	resp := map[string]interface{}{
		"modules_checker":   checkerList,
		"by_module_checker": checkerMap,
		"statuses":          []string{"PENDING", "APPROVED", "REJECTED"},
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, resp)
}

// ApproverAuditorAllocations returns auditor allocations for the caller's role,
// grouped by module with request_actions and derived action_types.
func (a *bpsActionAdapter) ApproverAuditorAllocations(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "approverAuditorAllocations", "handler", "cpsAction")
	defer span.End()
	roleCode, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(roleCode) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}
	repo := mid.GetBPSActionApproveRepo()
	if repo == nil {
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}
	// maker, checker, auditor, portalCards
	_, _, auditorMods, err := repo.PopulateUserApproverAllocations(ctx, roleCode)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	type modInfo struct {
		RequestActions []string `json:"request_actions"`
		ActionTypes    []string `json:"action_types"`
	}
	deriveTypes := func(actions []string) []string {
		seen := map[string]struct{}{}
		for _, ra := range actions {
			u := strings.ToUpper(ra)
			switch {
			case strings.Contains(u, "CREATE"):
				seen["CREATE"] = struct{}{}
			case strings.Contains(u, "UPDATE"):
				seen["UPDATE"] = struct{}{}
			case strings.Contains(u, "DELETE"):
				seen["DELETE"] = struct{}{}
			case strings.Contains(u, "ENABLE"):
				seen["ENABLE"] = struct{}{}
			case strings.Contains(u, "DISABLE"):
				seen["DISABLE"] = struct{}{}
			}
		}
		out := make([]string, 0, len(seen))
		for k := range seen {
			out = append(out, k)
		}
		return out
	}
	build := func(mods []string) (map[string]modInfo, []string) {
		m := map[string]modInfo{}
		uniq := map[string]struct{}{}
		for _, raw := range mods {
			mod := strings.ToUpper(strings.TrimSpace(raw))
			if mod == "" {
				continue
			}
			uniq[mod] = struct{}{}
			var reqs []string
			if group, ok := bpsactionsvc.RequestActionGroups[mod]; ok {
				reqs = make([]string, 0, len(group))
				for _, ga := range group {
					reqs = append(reqs, string(ga))
				}
			} else {
				reqs = []string{}
			}
			m[mod] = modInfo{RequestActions: reqs, ActionTypes: deriveTypes(reqs)}
		}
		list := make([]string, 0, len(uniq))
		for k := range uniq {
			list = append(list, k)
		}
		return m, list
	}
	auditorMap, auditorList := build(auditorMods)
	resp := map[string]interface{}{
		"modules_auditor":   auditorList,
		"by_module_auditor": auditorMap,
		"statuses":          []string{"PENDING", "APPROVED", "REJECTED"},
	}
	localization.SendSuccessResponse(w, localization.SuccessBPSActionsRetrieved, resp)
}
