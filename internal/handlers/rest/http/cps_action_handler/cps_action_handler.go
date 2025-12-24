package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	cpsactionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	cpsaction "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	mid "cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/service"
	cpsactionsvc "cbe-super-app-cps-action/internal/service/cps_action"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
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

	if action.ActionStatus != string(constants.Pending) {
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
		rawRoleID, _ := r.Context().Value(constants.ContextKey("role_id")).(string)
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
	localization.SendSuccessResponse(w, localization.SuccessCPSActionAuthorized, nil)
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
	var idxDoc *model.CPSActionApproveIndex
	if mod, ok := cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(action.RequestAction)); ok {
		actionName = mod
	}
	if repo := mid.GetCPSActionApproveRepo(); repo != nil && actionName != "" {
		rawRoleID, _ := r.Context().Value(constants.ContextKey("role_id")).(string)
		if rawRoleID == "" {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}
		roleID := local_util.FirstHex24(rawRoleID)
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
		// expected := int32(*idxDoc.CheckerIndex)
		// ctx = context.WithValue(ctx, constants.ContextKey("role_checker_index"), *idxDoc.CheckerIndex)
		// ctx = context.WithValue(ctx, constants.ContextKey("role_checker_group"), expected)
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
	// if int32(*idxDoc.CheckerIndex) != currentIndex+1 || int32(*idxDoc.CheckerIndex) > checkerCount {
	// 	span.RecordError(fmt.Errorf("out of order checker approval"))
	// 	localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
	// 	return
	// }
	// Build approval update inline (only mark Approved on final checker)
	finalStatus := string(constants.Pending)
	if int32(currentIndex+1) == TotalCheckerCount {
		finalStatus = string(constants.Approved)
	}
	checkerUser := model.Checker{
		CheckerID:          userData.UserID,
		RoleID:             r.Context().Value(constants.ContextKey("role_id")).(string),
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
		fmt.Printf("Errors : %v\n", err)
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCanceled, nil)
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
	idx32 := int32(action.CurrentCheckerIndex) + 1
	r = r.WithContext(context.WithValue(r.Context(), constants.ContextKey("checker_index"), idx32))
	ctx = context.WithValue(ctx, constants.ContextKey("checker_index"), idx32)

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorCPSActionRejectionPayloadDecodeFailed.Message)
		return
	}

	if req.Validate() != nil {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionRejectionReason)
		return
	}

	span.SetAttributes(
		attribute.String("cps_action.code", actionCode),
		attribute.String("cps_action.rejection_reason", req.RejectionReason),
	)

	// Enforce ordering and role-based approver index on rejection as well
	checkerCount := action.CheckerCount
	currentIndex := int32(action.CurrentCheckerIndex)
	if idx32 != currentIndex+1 || idx32 > checkerCount {
		span.RecordError(fmt.Errorf("out of order checker rejection"))
		localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
		return
	}
	{
		actionName := ""
		if mod, ok := cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(action.RequestAction)); ok {
			actionName = mod
		}
		if repo := mid.GetCPSActionApproveRepo(); repo != nil && actionName != "" {
			roleID, _ := r.Context().Value(constants.ContextKey("role_id")).(string)
			if roleID == "" {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			idxDoc, err := repo.FindByRoleAndAction(ctx, roleID, actionName)
			if err != nil {
				span.RecordError(err)
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			if idxDoc == nil || idxDoc.CheckerIndex == nil {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
			expected := int32(*idxDoc.CheckerIndex)
			ctx = context.WithValue(ctx, constants.ContextKey("role_checker_index"), *idxDoc.CheckerIndex)
			ctx = context.WithValue(ctx, constants.ContextKey("role_checker_group"), expected)
			r = r.WithContext(ctx)
			if expected != idx32 || expected != currentIndex+1 || expected > checkerCount {
				localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
				return
			}
			for _, cu := range action.CheckerUsers {
				if cu.RoleID == roleID || cu.CheckerID == userData.UserID || cu.CheckerIndex == expected {
					localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
					return
				}
			}
		}
	}

	CheckerUser := model.Checker{
		CheckerID:          userData.UserID,
		RoleID:             r.Context().Value(constants.ContextKey("role_id")).(string),
		CheckerIndex:       idx32,
		CheckerName:        userData.FullName,
		CheckerPhoneNumber: userData.PhoneNumber,
		ApprovedAt:         time.Now(),
	}
	if err := a.cpsActionApplication.RejectCPSAction(ctx, actionCode, &model.CPSAction{
		ActionCode:          actionCode,
		ActionStatus:        constants.Rejected,
		RejectionReason:     req.RejectionReason,
		CheckerUsers:        append(action.CheckerUsers, CheckerUser),
		CurrentCheckerIndex: float64(idx32),
		Department:          userData.Department,
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
