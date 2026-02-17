package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	cpsactionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	cps_actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/cps_action_role"
	cpsaction "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	mid "cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/service"
	cpsactionsvc "cbe-super-app-cps-action/internal/service/cps_action"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type cps_action_resp *model.CPSAction
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

// AuditorAction handles both claim (start) and mark (complete) for auditor in a single endpoint
//
//	@Summary		Auditor action (claim or mark)
//	@Description	Allows an auditor to claim (start) or mark (complete) an audit for a CPS action. Send empty mark to claim, or include mark (APPROVED/REJECTED) and reason to complete.
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path		string									true	"Action Code"
//	@Param			request		body		cps_actionrole_dto.AuditorMarkRequest	true	"Auditor mark request (empty mark to claim, or mark + reason to complete)"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"Auditor action processed successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input or operation not allowed"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/{action_code}/auditor [patch]
func (a *cpsActionAdapter) AuditorAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "auditorAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	// Load action (to resolve module/request_action)
	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, "")
	if err != nil || action == nil {
		span.RecordError(err)
		if err.Error() == localization.ErrorActionNotFound.Code {
			localization.SendErrorByCodeResponse(w, localization.ErrorActionDataNotFound.Code)
			return
		}
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	if action.AuditorStatus == model.AuditorStatus(constants.AUDITORCHECKED) {
		localization.SendBadRequestResponse(w, localization.ErrorAuditorActionOnThisActionCompleted.Message)
		return
	}

	if action.ActionStatus == string(constants.Pending) || action.ActionStatus == string(constants.Canceled) || action.AuditorCount == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}
	// Determine caller allocations and active auditor group
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	actionName := ""
	if mod, ok := cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(action.RequestAction)); ok {
		actionName = mod
	}

	repo := mid.GetCPSActionApproveRepo()
	if repo == nil || actionName == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	rawRoleID, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
	if strings.TrimSpace(rawRoleID) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	UpperCaseAction := strings.ToUpper(actionName)
	idxDoc, err := repo.FindByRoleAndAction(ctx, rawRoleID, UpperCaseAction, action.Version)
	if err != nil || idxDoc == nil || idxDoc.AuditorIndex == nil {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	// Active group is the integer part of the auditor index (e.g., 1.* -> 1)
	activeGroup := int(*idxDoc.AuditorIndex)
	if action.CurrentAuditorIndex >= float64(activeGroup) {
		localization.SendBadRequestResponse(w, localization.ErrorAuditorActionOnThisRoleCompleted.Message)
		return
	}

	if (action.CurrentAuditorIndex + 1) > float64(activeGroup) {
		localization.SendBadRequestResponse(w, localization.ErrorAuditorActionWaitForPreviousAuditor.Message)
		return
	}

	currentIndex := action.CurrentAuditorIndex

	Current_role_level := *idxDoc.AuditorIndex
	expected := int32(*idxDoc.AuditorIndex)
	ctx = context.WithValue(ctx, constants.ContextKey("role_checker_index"), *idxDoc.AuditorIndex)
	ctx = context.WithValue(ctx, constants.ContextKey("role_checker_group"), expected)
	r = r.WithContext(ctx)
	if currentIndex == float64(Current_role_level) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionApprovedByThisRole)
		return
	}

	if int64(currentIndex)+1 < int64(Current_role_level) {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
		return
	}

	if action.MakerID == userData.UserID {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	for _, au := range action.AuditorUsers {
		if au.RoleID == rawRoleID || au.AuditorID == userData.UserID {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}
	}

	// Parse request to decide claim vs mark
	var reqBody cps_actionrole_dto.AuditorMarkRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	if action.MakerID == userData.UserID {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	if strings.TrimSpace(reqBody.Mark) == "" {
		// Claim path -> move status to INPROGRESS when allowed
		if err := a.cpsActionApplication.AuditorClaim(ctx, actionCode, activeGroup); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		}
		localization.SendSuccessResponse(w, localization.SuccessCPSActionChecked, nil)
		return
	}

	// Mark path -> record auditor mark and advance group/finish
	auditor := model.Auditor{
		AuditorID:          userData.UserID,
		RoleID:             rawRoleID,
		AuditorIndex:       int32(activeGroup),
		AuditorName:        userData.FullName,
		AuditorPhoneNumber: userData.PhoneNumber,
		AuditorReason:      reqBody.Reason,
		AuditorMark:        model.AuditorMark(strings.ToUpper(strings.TrimSpace(reqBody.Mark))),
		ApprovedAt:         time.Now(),
	}

	if err := a.cpsActionApplication.AuditorMark(ctx, actionCode, auditor, activeGroup); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionChecked, nil)
}

// CancelCPSAction cancels a CPS action
//
//	@Summary		Cancel CPS action
//	@Description	Cancels a PENDING CPS action. Only the maker who created the action can cancel it, and only if no checker has approved it yet.
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			action_code	path		string									true	"Action Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS action canceled successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - Action cannot be canceled"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"Action not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/{action_code}/cancel [patch]
func (a *cpsActionAdapter) CancelCPSAction(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "cancelCpsAction", "handler", "cpsAction")
	defer span.End()
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	var req cps_actionrole_dto.CancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

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
	action.CanceledReason = string(req.CancelReason)

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

		idxDoc, err := repo.FindByRoleAndAction(ctx, roleID, strings.ToUpper(actionName), action.Version)
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
		idxDoc, err = repo.FindByRoleAndAction(ctx, roleID, UpperCaseAction, action.Version)
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
		if currentIndex == float64(Current_role_level) {
			localization.SendBadRequestResponse(w, localization.MsgCPSActionApprovedByThisRole)
			return
		}

		if int64(currentIndex)+1 < int64(Current_role_level) {
			localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
			return
		}

		if action.MakerID == userData.UserID {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
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
		CurrentCheckerIndex: float64(*idxDoc.CheckerIndex),
		CheckerUsers:        append(action.CheckerUsers, checkerUser),
		RoleCode:            r.Context().Value(constants.ContextKey("role_code")).(string),
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
//	@Param			request		body		cpsaction.ActionRequest	true	"Rejection request"
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	// Fetch action and attach checker_index context for this approver
	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, "")
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	currentIndex := action.CurrentCheckerIndex

	//---------------------------------------------------
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
		idxDoc, err = repo.FindByRoleAndAction(ctx, roleID, UpperCaseAction, action.Version)
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
		if currentIndex == float64(Current_role_level) {
			localization.SendBadRequestResponse(w, localization.MsgCPSActionApprovedByThisRole)
			return
		}

		if int64(currentIndex)+1 < int64(Current_role_level) {
			localization.SendBadRequestResponse(w, localization.MsgCPSActionWaitPrevious)
			return
		}

		if action.MakerID == makerData.UserID {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		for _, cu := range action.CheckerUsers {
			if cu.RoleID == roleID || cu.CheckerID == makerData.UserID {
				localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
				return
			}
		}

	}
	//---------------------------------------------------

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
		ActionCode:          actionCode,
		CurrentCheckerIndex: *idxDoc.CheckerIndex,
		ActionStatus:        constants.Rejected,
		RejectionReason:     req.RejectionReason,
		CheckerUsers:        append(action.CheckerUsers, CheckerUser),
		RoleCode:            r.Context().Value(constants.ContextKey("role_code")).(string),
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

	actions, err := a.cpsActionApplication.GetCPSActionsByDepartment(ctx, userData.Department, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("cps_action.department", userData.Department))
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, actions)
}

// GetUserCheckedActions retrieves CPS actions checked by the current user
//
//	@Summary		Get user checked actions
//	@Description	Retrieves a paginated list of CPS actions that have been checked (approved/rejected) by the current user
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
//	@Router			/actions/user/checked/actions [get]
func (a *cpsActionAdapter) GetUserCheckedActions(w http.ResponseWriter, r *http.Request) {
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
	res, err := a.cpsActionApplication.GetUserCheckedActions(ctx, userID, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

// GetUserCreatedActions retrieves CPS actions created by the current user
//
//	@Summary		Get user created actions
//	@Description	Retrieves a paginated list of CPS actions created by the current user
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
//	@Router			/actions/user/created/actions [get]
func (a *cpsActionAdapter) GetUserCreatedActions(w http.ResponseWriter, r *http.Request) {
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

	userData := local_util.ExtractUserContext(r)
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserCreatedActions", "handler", "cpsAction")
	defer span.End()

	userID := userData.UserName
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
//	@Param			action_code	path		string	true	"Action Code"
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

// GetUserApproverActions retrieves CPS actions pending approval for the current user's approver role
//
//	@Summary		Get approver checker actions
//	@Description	Retrieves a paginated list of CPS actions that are pending approval for the current user's approver/checker role
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
//	@Router			/actions/approver/checker/actions [get]
func (a *cpsActionAdapter) GetUserApproverActions(w http.ResponseWriter, r *http.Request) {
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
	idxRepo := mid.GetCPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}
	_, _, checkerActions, _, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if checkerActions == nil {
		localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, map[string]interface{}{})
		return
	}

	a.logger.Infof("Checker Actions: %v", checkerActions)

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

	res, err := a.cpsActionApplication.GetCPSActionsForApprover(ctx, userID, reqs, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

// GetUserAuditorActions retrieves CPS actions pending audit for the current user's auditor role
//
//	@Summary		Get auditor checker actions
//	@Description	Retrieves a paginated list of CPS actions that are pending audit for the current user's auditor role
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
//	@Router			/actions/auditor/checker/actions [get]
func (a *cpsActionAdapter) GetUserAuditorActions(w http.ResponseWriter, r *http.Request) {
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
	idxRepo := mid.GetCPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	_, _, _, auditorAllocations, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if auditorAllocations == nil {
		localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, map[string]interface{}{})
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

	// do not force action_status; let API-provided filters decide
	userID := local_util.ExtractUserContext(r).UserID
	res, err := a.cpsActionApplication.GetCPSActionsForAuditor(ctx, userID, reqs, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

func (a *cpsActionAdapter) GetUserApproverApprovedActions(w http.ResponseWriter, r *http.Request) {
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

	idxRepo := mid.GetCPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}
	_, _, checkerActions, _, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
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

	userID := local_util.ExtractUserContext(r).UserID
	res, err := a.cpsActionApplication.GetCPSActionsForApprover(ctx, userID, reqs, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, res)
}

// GetActionCounts retrieves action counts for the current user's role
//
//	@Summary		Get action counts
//	@Description	Retrieves counts of CPS actions by status (pending, approved, rejected, canceled, inprogress, completed) for the current user's role (maker/checker/auditor)
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			role	query		string															true	"Role type: maker, checker, or auditor"
//	@Success		200		{object}	localization.StandardResponse{data=cpsaction.CPSActionCountResponse}	"Action counts retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/counts [get]
func (a *cpsActionAdapter) GetActionCounts(w http.ResponseWriter, r *http.Request) {
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

	idxRepo := mid.GetCPSActionApproveRepo()
	if idxRepo == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	_, makerActions, checkerActions, auditorActions, _, err := idxRepo.PopulateUserApproverAllocations(ctx, rawRoleID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if makerActions == nil && checkerActions == nil && auditorActions == nil {
		resp := &cpsactionDto.CPSActionCountResponse{
			Pending:    0,
			Approved:   0,
			Rejected:   0,
			Canceled:   0,
			Inprogress: 0,
			Completed:  0,
		}
		localization.SendSuccessResponse(w, localization.SuccessCPSActionCount, resp)
		return
	}

	// resolve action_names -> request_actions (same as GetUserApproverActions)
	var reqs []string
	seen := map[string]struct{}{}

	if makerActions != nil && requestedRole == "maker" {
		for _, mod := range makerActions {
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

	}

	if checkerActions != nil && requestedRole == "checker" {
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
	}

	if auditorActions != nil && requestedRole == "auditor" {
		for _, mod := range auditorActions {
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
	}

	// If no mapped request actions, return zero counts
	if len(reqs) == 0 {
		resp := &cpsactionDto.CPSActionCountResponse{
			Pending:    0,
			Approved:   0,
			Rejected:   0,
			Canceled:   0,
			Inprogress: 0,
			Completed:  0,
		}
		localization.SendSuccessResponse(w, localization.SuccessCPSActionCount, resp)
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

	var pendingCount, approvedCount, rejectedCount, canceledCount, inprogressAuditCount, completedAuditCount int

	// Pending
	if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Pending), "")); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	} else if res != nil && res.Meta.TotalDocs > 0 {
		pendingCount = int(res.Meta.TotalDocs)
	}

	// Approved
	if auditorActions != nil && requestedRole == "auditor" {
		if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Approved), string(model.AUDITORNOTCHECKED))); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			approvedCount = int(res.Meta.TotalDocs)
		}
	} else {
		if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Approved), "")); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			approvedCount = int(res.Meta.TotalDocs)
		}

	}

	// Rejected
	if auditorActions != nil && requestedRole == "auditor" {
		if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Rejected), string(model.AUDITORNOTCHECKED))); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			rejectedCount = int(res.Meta.TotalDocs)
		}
	} else {
		if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Rejected), "")); err != nil {
			span.RecordError(err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		} else if res != nil && res.Meta.TotalDocs > 0 {
			rejectedCount = int(res.Meta.TotalDocs)
		}
	}

	// Canceled
	if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter(string(constants.Canceled), "")); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	} else if res != nil && res.Meta.TotalDocs > 0 {
		canceledCount = int(res.Meta.TotalDocs)
	}

	// Audior's Inprogress Count
	if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter("", string(model.AUDITORINPROGRESS))); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	} else if res != nil && res.Meta.TotalDocs > 0 {
		inprogressAuditCount = int(res.Meta.TotalDocs)
	}

	// Audior's Completed Count
	if res, err := a.cpsActionApplication.GetCPSActions(ctx, userID, requestedRole, reqs, buildFilter("", string(model.AUDITORCHECKED))); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	} else if res != nil && res.Meta.TotalDocs > 0 {
		completedAuditCount = int(res.Meta.TotalDocs)
	}

	resp := &cpsactionDto.CPSActionCountResponse{
		Pending:    pendingCount,
		Approved:   approvedCount,
		Rejected:   rejectedCount,
		Canceled:   canceledCount,
		Inprogress: inprogressAuditCount,
		Completed:  completedAuditCount,
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCount, resp)
}

// GetAuthorizerIndex retrieves the authorizer index for a specific request action
//
//	@Summary		Get authorizer index
//	@Description	Retrieves the authorizer index configuration for a specific request action type
//	@Tags			CPS Actions
//	@Accept			json
//	@Produce		json
//	@Param			request_action	path		string									true	"Request Action (e.g., CREATE_BANK, UPDATE_BANK)"
//	@Success		200				{object}	localization.StandardResponse{data=object}	"Authorizer index retrieved successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}		"Bad request - Invalid request action"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/actions/authorizer/index/{request_action} [get]
func (a *cpsActionAdapter) GetAuthorizerIndex(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getUserAuthorizerIndex", "handler", "cpsAction")
	defer span.End()
	requestAction := chi.URLParam(r, "request_action")

	if requestAction == "" {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	authorizerIndex, err := a.cpsActionApplication.GetUserAuthorizerIndex(ctx, constants.RequestAction(requestAction))
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[CPSAction.GetActionCounts] service failed %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, authorizerIndex)

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
	// viewer, maker, checker, auditor, portalCards
	_, _, checkerMods, _, _, err := repo.PopulateUserApproverAllocations(ctx, roleCode)
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
	// viewer, maker, checker, auditor, portalCards
	_, _, _, auditorMods, _, err := repo.PopulateUserApproverAllocations(ctx, roleCode)
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

func (a *cpsActionAdapter) GetAuthorizersLevel(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "rejectCpsAction", "handler", "cpsAction")
	defer span.End()

	requestAction := chi.URLParam(r, "request_action")
	actionVersion := chi.URLParam(r, "action_version")
	parsedVersion, err := strconv.ParseInt(actionVersion, 10, 64)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
		return
	}

	actionName := ""
	var idxDoc *imodel.CPSActionApproveIndex
	if mod, ok := cpsactionsvc.ResolveModuleForRA(cpsactionsvc.RequestAction(requestAction)); ok {
		actionName = mod
	}

	var checkerIdx, auditorIdx int64
	if repo := mid.GetCPSActionApproveRepo(); repo != nil && actionName != "" {
		role, _ := r.Context().Value(constants.ContextKey("role_code")).(string)
		if role == "" {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		uppercasedActionName := strings.ToUpper(actionName)
		idxDoc, err = repo.FindByRoleAndAction(ctx, role, uppercasedActionName, parsedVersion)
		if err != nil {
			span.RecordError(err)
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		if idxDoc == nil || idxDoc.CheckerIndex == nil {
			localization.SendBadRequestResponse(w, localization.ErrorOperationNotAllowed.Message)
			return
		}

		if idxDoc.CheckerIndex != nil {
			checkerIdx = int64(*idxDoc.CheckerIndex)
		}
		if idxDoc.AuditorIndex != nil {
			auditorIdx = int64(*idxDoc.AuditorIndex)
		}
	}

	autorizersLevel := cpsactionDto.AutorizersLevelResponse{
		CheckerIndex: checkerIdx,
		AuditorIndex: auditorIdx,
	}

	localization.SendSuccessResponse(w, localization.AutorizersLevelFetchedSuccessfully, autorizersLevel)
}
