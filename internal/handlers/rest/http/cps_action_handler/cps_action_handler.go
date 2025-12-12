package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	cpsactionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	cpsaction "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	core "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"fmt"
	"net/http"

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

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUserForbidden.Message)
		return
	}
	// First, retrieve the existing CPS action to get ALL the data
	span.SetAttributes(attribute.String("cps_action.code", actionCode))
	action, err := a.cpsActionApplication.GetCPSActionByActionCode(ctx, actionCode, userData.Department)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	approvalAction := core.MapCPSActionToApproval(action, &userData)
	// Then, approve the action
	if err := a.cpsActionApplication.ApproveCPSAction(ctx, approvalAction); err != nil {
		fmt.Printf("Errors : %v\n", err)
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

	if err := a.cpsActionApplication.RejectCPSAction(ctx, actionCode, &model.CPSAction{
		ActionCode:         actionCode,
		ActionStatus:       constants.Rejected,
		RejectionReason:    req.RejectionReason,
		CheckerID:          userData.UserID,
		CheckerName:        userData.FullName,
		CheckerPhoneNumber: userData.PhoneNumber,
		Department:         userData.Department,
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
//	@Param			action_code	path		string								whitespace				true	"Action Code"
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
