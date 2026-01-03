package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	cpsactionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	cpsaction "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	mid "cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/service"
	cpsactionsvc "cbe-super-app-cps-action/internal/service/cps_action"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"
	"net/http"

	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type cps_action_resp *model.CPSAction
type cps_action_dto_Resp *cpsactionDto.ActionRequest
type cps_actions_paginated_resp *types.PaginatedResponse[[]*model.CPSAction]

type cpsActionAdapter struct {
	cpsActionApplication service.CPSActionService
	logger               utils.Logger
}

func InitCPSActionAdapter(cpsActionApplication service.CPSActionService, logger utils.Logger) cpsaction.CPSActionAdapter {
	return &cpsActionAdapter{
		logger:               logger,
		cpsActionApplication: cpsActionApplication,
	}
}

// CancelCPSAction implements cps_action.CPSActionAdapter.
func (a *cpsActionAdapter) CancelCPSAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "cancelCpsAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if action.ActionStatus != string(constants.Pending) || action.CurrentCheckerIndex > 0 {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUserForbidden.Message)
		return
	}
	if userData.UserID != action.MakerID {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	action.ActionStatus = string(constants.Canceled)

	if err := a.cpsActionApplication.RejectCPSAction(ctx, actionCode, action); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionCanceled, nil)

}

// ReverseCPSAction reverses an already-approved CPS action
//
//	@Summary		Reverse CPS action (auditor only)
//	@Description	Reverses an APPROVED CPS action by applying the inverse operation. Restricted to auditor roles.
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path	string	true	"Action Code"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"CPS action reversed successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		403	{object}	localization.StandardResponse{data=nil}	"Forbidden"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Action not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/{action_code}/reverse [patch]
func (a *cpsActionAdapter) ReverseCPSAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "reverseCpsAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	// Load action and validate status
	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if action.ActionStatus != string(constants.Approved) {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	// Validate caller is auditor for this module via approver index (role_id + action group)
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	actionName := ""
	if mod, ok := cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(action.RequestAction)); ok {
		actionName = mod
	}
	if repo := mid.GetCPSActionApproveRepo(); repo != nil && actionName != "" {
		rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
		roleID := local_util.FirstHex24(rawRoleID)
		if roleID == "" {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		idxDoc, err := repo.FindByRoleAndAction(ctx, roleID, strings.ToUpper(actionName))
		if err != nil || idxDoc == nil || idxDoc.AuditorIndex == nil {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}
	}

	// Perform reversal via service
	if err := a.cpsActionApplication.ReverseCPSAction(ctx, actionCode); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("cps_action.code", actionCode),
		attribute.String("auditor.id", userData.UserID),
	)
	localization.SendSuccessResponse(w, localization.SuccessCPSActionReversed, nil)
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
func (a *cpsActionAdapter) ApproveCPSAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "approveCpsAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	// Fetch current action first to derive checker_index and attach it to context
	span.SetAttributes(attribute.String("cps_action.code", actionCode))
	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	// Prevent approving an action that is already finalized
	if action.ActionStatus == string(constants.Approved) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionAlreadyApproved)
		return
	}

	if action.ActionStatus == string(constants.Rejected) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionAlreadyRejected)
		return
	}

	if action.ActionStatus == string(constants.Canceled) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionAlreadyRejected)
		return
	}

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUserForbidden.Message)
		return
	}

	TotalCheckerCount := action.CheckerCount
	currentIndex := action.CurrentCheckerIndex

	// Validate approver role's checker_index via cps_action_approver_index (grouped: 0.* -> 1.*, 1.* -> 2.*)
	actionName := ""
	var idxDoc *imodel.CPSActionApproveIndex
	if mod, ok := cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(action.RequestAction)); ok {
		actionName = mod
	}

	if repo := mid.GetCPSActionApproveRepo(); repo != nil && actionName != "" {
		rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
		if rawRoleID == "" {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		roleID := rawRoleID
		UpperCaseAction := strings.ToUpper(actionName)
		idxDoc, err = repo.FindByRoleAndAction(ctx, roleID, UpperCaseAction)
		if err != nil {
			span.RecordError(err)
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		if idxDoc == nil || idxDoc.CheckerIndex == nil {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		Current_role_level := *idxDoc.CheckerIndex
		expected := int32(*idxDoc.CheckerIndex)
		ctx = context.WithValue(ctx, constants.ContextKey("role_checker_index"), *idxDoc.CheckerIndex)
		ctx = context.WithValue(ctx, constants.ContextKey("role_checker_group"), expected)
		r = r.WithContext(ctx)
		if currentIndex == Current_role_level {
			localization.SendBadRequestResponse(w, localization.MsgCPSActionApprovedByThisRole)
			return
		}

		if int64(currentIndex)+1 < int64(Current_role_level) {
			localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
			return
		}

		for _, cu := range action.CheckerUsers {
			if cu.RoleID == roleID || cu.CheckerID == userData.UserID {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
		}
	}

	// Enforce ordering: must approve in sequence
	// if int64(*idxDoc.CheckerIndex) != int64(currentIndex)+1 || int64(*idxDoc.CheckerIndex) > int64(checkerCount) {
	// 	span.RecordError(fmt.Errorf("out of order checker approval"))
	// 	localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
	// 	return
	// }
	// Build approval update inline (only mark Approved on final checker)
	finalStatus := string(constants.Pending)
	if int32(*idxDoc.CheckerIndex) == TotalCheckerCount {
		finalStatus = string(constants.Approved)
	}
	checkerUser := model.Checker{
		CheckerID:          userData.UserID,
		RoleID:             r.Context().Value(constants.ContextKey("role_code")).(string),
		CheckerIndex:       int32(*idxDoc.CheckerIndex),
		CheckerName:        userData.FullName,
		CheckerPhoneNumber: userData.PhoneNumber,
		ApprovedAt:         time.Now(),
	}
	update := &model.CPSAction{
		ActionCode:          action.ActionCode,
		ActionStatus:        finalStatus,
		CurrentCheckerIndex: *idxDoc.CheckerIndex,
		CheckerUsers:        append(action.CheckerUsers, checkerUser),
		Department:          userData.Department,
	}

	// Then, approve the action
	if err := a.cpsActionApplication.ApproveCPSAction(ctx, update); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionAuthorized, nil)
}

// RejectCPSAction rejects a CPS action
//
//	@Summary		Reject CPS action
//	@Description	Rejects a CPS action by action code with rejection reason
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path		string									true	"Action Code"
//	@Param			request		body		cps_action_dto_Resp						true	"Rejection request"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS action rejected successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input or missing rejection reason"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/{action_code}/reject [patch]
func (a *cpsActionAdapter) RejectCPSAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "rejectCpsAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))
	var req cpsactionDto.ActionRequest
	makerData := local_util.ExtractUserFromContext(ctx)

	// Fetch action and attach checker_index context for this approver
	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	// Prevent rejecting an action that is already finalized
	if action.ActionStatus == string(constants.Approved) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionAlreadyApproved)
		return
	}
	if action.ActionStatus == string(constants.Rejected) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionAlreadyRejected)
		return
	}
	if action.ActionStatus == string(constants.Canceled) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionAlreadyCanceled)
		return
	}

	CheckerUser := model.Checker{
		CheckerID:          makerData.UserID,
		RoleID:             r.Context().Value(constants.ContextKey("role_code")).(string),
		CheckerName:        makerData.FullName,
		CheckerPhoneNumber: makerData.PhoneNumber,
		ApprovedAt:         time.Now(),
	}
	if err := a.cpsActionApplication.RejectCPSAction(ctx, actionCode, &model.CPSAction{
		ActionCode:      actionCode,
		ActionStatus:    constants.Rejected,
		RejectionReason: req.RejectionReason,
		CheckerUsers:    append(action.CheckerUsers, CheckerUser),
		Department:      makerData.Department,
	}); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionRejected, nil)
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
//	@Success		200			{object}	localization.StandardResponse{data=cps_actions_paginated_resp}	"CPS actions retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/ [get]
func (a *cpsActionAdapter) GetCPSActionsByDepartment(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionsByDepartment", "handler", "cpsAction")
	defer span.End()

	actions, err := a.cpsActionApplication.GetCPSActionsByDepartment(ctx, userData.Department, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("cps_action.department", userData.Department))
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, actions)
}

func (a *cpsActionAdapter) GetUserCheckedActions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserApprovedCpsActions", "handler", "cpsAction")
	defer span.End()
	userData := local_util.ExtractUserContext(r)
	userID := userData.UserID
	res, err := a.cpsActionApplication.GetUserCheckedActions(ctx, userID, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

func (a *cpsActionAdapter) GetUserCreatedActions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	userData := local_util.ExtractUserContext(r)
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserCreatedActions", "handler", "cpsAction")
	defer span.End()

	userID := userData.UserID
	res, err := a.cpsActionApplication.GetUserCreatedActions(ctx, userID, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

// GetCPSActionByID retrieves a CPS action by ID
//
//	@Summary		Get CPS action by ID
//	@Description	Retrieves a specific CPS action by its ID
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_id	path		string												true	"Action ID"
//	@Success		200			{object}	localization.StandardResponse{data=cps_action_resp}	"CPS action retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}				"Bad request - Invalid action ID"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}				"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/by-id/{action_id} [get]
func (a *cpsActionAdapter) GetCPSActionByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionById", "handler", "cpsAction")
	defer span.End()
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	actionID := chi.URLParam(r, string(constants.ActionID))

	span.SetAttributes(
		attribute.String("cps_action.id", actionID),
		attribute.String("cps_action.department", userData.Department),
	)

	action, err := a.cpsActionApplication.GetCPSActionByID(ctx, actionID, userData.Department)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionFetched, action)
}

// GetCPSActionByActionCode retrieves a CPS action by action code with history
//
//	@Summary		Get CPS action by action code with history
//	@Description	Retrieves a specific CPS action by its action code, including previous and current actions
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path		string												whitespace	true	"Action Code"
//	@Success		200			{object}	localization.StandardResponse{data=cps_action_resp}	"CPS action with history retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}				"Bad request - Invalid action code"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}				"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/by-action-code/{action_code} [get]
func (a *cpsActionAdapter) GetCPSActionByActionCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionByCode", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	span.SetAttributes(
		attribute.String("cps_action.code", actionCode),
		attribute.String("cps_action.department", userData.Department),
	)

	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, userData.Department)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionFetched, action)
}

func (a *cpsActionAdapter) GetUserApproverPendingActions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserApproverPendingActions", "handler", "cpsAction")
	defer span.End()

	// roleCode from context
	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if rawRoleID == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	// fetch checker allocations for this role
	idxRepo := mid.GetCPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}
	_, checkerActions, _, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if checkerActions == nil {
		localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, map[string]interface{}{})
		return
	}

	fmt.Println("Checker Actions", checkerActions)
	// resolve action_names -> request_actions
	var reqs []string
	seen := map[string]struct{}{}
	for _, mod := range checkerActions {
		upper := strings.ToUpper(strings.TrimSpace(mod))
		if lst, ok := cpsactionsvc.RequestActionGroups[upper]; ok {
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

	res, err := a.cpsActionApplication.GetCPSActionsByDepartment(ctx, "", filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

func (a *cpsActionAdapter) GetUserApproverApprovedActions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserApproverApprovedActions", "handler", "cpsAction")
	defer span.End()

	// roleCode from context
	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if rawRoleID == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	idxRepo := mid.GetCPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}
	_, checkerActions, _, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	var reqs []string
	seen := map[string]struct{}{}
	for _, mod := range checkerActions {
		upper := strings.ToUpper(strings.TrimSpace(mod))
		if lst, ok := cpsactionsvc.RequestActionGroups[upper]; ok {
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

	res, err := a.cpsActionApplication.GetCPSActionsByDepartment(ctx, "", filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

func (a *cpsActionAdapter) GetActionCounts(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionCounts", "handler", "cpsAction")
	defer span.End()
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("cps_action.department", userData.Department))

	actions, err := a.cpsActionApplication.GetActionCountsByDepartemnt(ctx, userData.Department)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[CPSAction.GetActionCounts] service failed %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, actions)
}

// ApproverCheckerAllocations returns checker allocations for the caller's role,
// grouped by module with request_actions and derived action_types.
func (a *cpsActionAdapter) ApproverCheckerAllocations(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "approverCheckerAllocations", "handler", "cpsAction")
	defer span.End()
	roleCode, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(roleCode) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}
	repo := mid.GetCPSActionApproveRepo()
	if repo == nil {
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}
	// maker, checker, auditor, portalCards
	_, checkerMods, _, _, err := repo.PopulateUserApproverAllocations(ctx, roleCode)
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
			if group, ok := cpsactionsvc.RequestActionGroups[mod]; ok {
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
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, resp)
}

// ApproverAuditorAllocations returns auditor allocations for the caller's role,
// grouped by module with request_actions and derived action_types.
func (a *cpsActionAdapter) ApproverAuditorAllocations(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "approverAuditorAllocations", "handler", "cpsAction")
	defer span.End()
	roleCode, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(roleCode) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}
	repo := mid.GetCPSActionApproveRepo()
	if repo == nil {
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}
	// maker, checker, auditor, portalCards
	_, _, auditorMods, _, err := repo.PopulateUserApproverAllocations(ctx, roleCode)
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
			if group, ok := cpsactionsvc.RequestActionGroups[mod]; ok {
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
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, resp)
}
