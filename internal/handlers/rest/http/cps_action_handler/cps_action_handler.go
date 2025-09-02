package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	cpsactionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	cpsaction "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	core "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

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

func (a *cpsActionAdapter) ApproveCPSAction(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, string(constants.ActionCode))

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUserForbidden.Message)
		return
	}

	// First, retrieve the existing CPS action to get ALL the data
	existingAction, err := a.cpsActionApplication.GetCPSActionByActionCode(r.Context(), actionCode, userData.Department)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	// Map to approval action using helper
	approvalAction := core.MapCPSActionToApproval(existingAction, &userData)

	if err := a.cpsActionApplication.ApproveCPSAction(r.Context(), approvalAction); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionAuthorized, nil)
}

func (a *cpsActionAdapter) RejectCPSAction(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, string(constants.ActionCode))
	var req cpsactionDto.ActionRequest
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorCPSActionRejectionPayloadDecodeFailed.Message)
		return
	}

	if req.Validate() != nil {
		localization.SendBadRequestResponse(w, localization.MsgCPSActionRejectionReason)
		return
	}

	if err := a.cpsActionApplication.RejectCPSAction(r.Context(), actionCode, &model.CPSAction{
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
func (a *cpsActionAdapter) GetCPSActionsByDepartment(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	actions, err := a.cpsActionApplication.GetCPSActionsByDepartment(r.Context(), userData.Department, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionsRetrieved, actions)
}
func (a *cpsActionAdapter) GetCPSActionByID(w http.ResponseWriter, r *http.Request) {
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	actionID := chi.URLParam(r, string(constants.ActionID))

	action, err := a.cpsActionApplication.GetCPSActionByID(r.Context(), actionID, userData.Department)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionFetched, action)
}
func (a *cpsActionAdapter) GetCPSActionByActionCode(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, string(constants.ActionCode))
	userData, err := local_util.ParseUserContext(r)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorUserForbidden, nil, nil)
		return
	}

	action, err := a.cpsActionApplication.GetCPSActionByActionCode(r.Context(), actionCode, userData.Department)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorInternalServerError, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionFetched, action)
}
