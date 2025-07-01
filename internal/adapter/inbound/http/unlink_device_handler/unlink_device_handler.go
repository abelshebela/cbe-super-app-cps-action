package unlink_device_handler

import (
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/application/unlink"
	inbound "cbe-super-app-cps-action/internal/port/inbound/unlink"
	ctx_util "cbe-super-app-cps-action/pkgs/context"
	err_def "cbe-super-app-cps-action/pkgs/common"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	GENERAL = "General"
	AUTH    = "Auth"
	USER    = "User"
	TXN     = "Transaction"
	UNLINK  = "Unlink"
	NONE    = ""
)

const (
	InvalidJSONPayload    = "INVALID_JSON_PAYLOAD"
	InvalidInput          = "INVALID_INPUT"
	DeviceUnlinkSuccess   = "DEVICE_UNLINK_SUCCESS"
	UnlinkDecisionSuccess = "UNLINK_DECISION_SUCCESS"
)

type UnlinkHandler struct {
	unlinkApp unlink.ApplicationService
	logger    utils.Logger
}

func NewHTTPUnlinkHandler(service unlink.ApplicationService, logger utils.Logger) inbound.UnlinkPortHandler {
	return &UnlinkHandler{
		unlinkApp: service,
		logger:    logger,
	}
}

func (h *UnlinkHandler) writeResponse(w http.ResponseWriter, status int, data any) {
	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         status,
		Data:           data,
	}
	resp.SendJSON()
}

func (h *UnlinkHandler) getMessageMap(errorMsg string, group string, key ...string) map[string]string {
	if len(key) == 0 {
		return map[string]string{"message": errorMsg}
	}

	var groupMap map[string]err_def.ErrorDefinition

	switch group {
	case GENERAL:
		groupMap = err_def.DefineError.General
	case AUTH:
		groupMap = err_def.DefineError.Auth
	case USER:
		groupMap = err_def.DefineError.User
	case TXN:
		groupMap = err_def.DefineError.Transaction
	case UNLINK:
		groupMap = err_def.DefineError.Unlink
	default:
		// NONE group
		return map[string]string{"message": errorMsg}
	}

	if errDef, ok := groupMap[key[0]]; ok {
		return map[string]string{"message": errDef.Message}
	}

	// Fallback to custom error message
	return map[string]string{"message": errorMsg}
}

func (h *UnlinkHandler) UnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.UnlinkDeviceRequest
	// decode it
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UnlinkDevice] failed to decode request: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, InvalidJSONPayload))
		return
	}

	// request validate
	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UnlinkDevice] validation failed: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, InvalidInput))
		return
	}

	ctx := ctx_util.ExtractUserContext(r)

	userID := ctx.UserID
	branchCode := ctx.BranchCode
	department := ctx.Department

	// calling the application service
	if err := h.unlinkApp.UnlinkDevice(request.UserCode, userID, branchCode, department); err != nil {
		h.logger.Errorf("[UnlinkDevice] service error: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(err.Error(), NONE))
		return
	}

	h.logger.Infof("[UnlinkDevice] success for userCode: %s", request.UserCode)
	h.writeResponse(w, http.StatusOK, h.getMessageMap(NONE, UNLINK, DeviceUnlinkSuccess))
}

func (h *UnlinkHandler) ApproveUnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.ApproveUnlinkDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] failed to decode request: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, InvalidJSONPayload))
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[ApproveUnlinkDevice] validation failed: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, InvalidInput))
		return
	}
	ctx := ctx_util.ExtractUserContext(r)
	userID := ctx.UserID

	if err := h.unlinkApp.ApproveOrDecline(request.UserCode, request.Decision, request.Reason, userID); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] service error: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(err.Error(), NONE))
		return
	}

	h.logger.Infof("[ApproveUnlinkDevice] decision: %s for userCode: %s", request.Decision, request.UserCode)
	h.writeResponse(w, http.StatusOK, h.getMessageMap(NONE, UNLINK, UnlinkDecisionSuccess))
}
