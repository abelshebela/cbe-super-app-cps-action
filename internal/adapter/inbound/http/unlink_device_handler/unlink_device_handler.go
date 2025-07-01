package unlink_device_handler

import (
	"encoding/json"
	"net/http"

	ctx_util "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/unlink"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/unlink"

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
	INVALID_JSON_PAYLOAD    = "INVALID_JSON_PAYLOAD"
	INVALID_INPUT           = "INVALID_INPUT"
	DEVICE_UNLINK_SUCCESS   = "DEVICE_UNLINK_SUCCESS"
	UNLINK_DECISION_SUCCESS = "UNLINK_DECISION_SUCCESS"
)

type UnlinkHandler struct {
	unlinkApp unlink.ApplicationService
	logger        utils.Logger
}

func NewHTTPUnlinkHandler(service unlink.ApplicationService, logger utils.Logger) inbound.UnlinkPortHandler {
	return &UnlinkHandler{
		unlinkApp: service,
		logger:        logger,
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

	var groupMap map[string]errors.ErrorDefinition

	switch group {
	case GENERAL:
		groupMap = errors.DefineError.General
	case AUTH:
		groupMap = errors.DefineError.Auth
	case USER:
		groupMap = errors.DefineError.User
	case TXN:
		groupMap = errors.DefineError.Transaction
	case UNLINK:
		groupMap = errors.DefineError.Unlink
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
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, INVALID_JSON_PAYLOAD))
		return
	}

	// request validate
	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UnlinkDevice] validation failed: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, INVALID_INPUT))
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
	h.writeResponse(w, http.StatusOK, h.getMessageMap(NONE, UNLINK, DEVICE_UNLINK_SUCCESS))
}

func (h *UnlinkHandler) ApproveUnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.ApproveUnlinkDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] failed to decode request: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, INVALID_JSON_PAYLOAD))
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[ApproveUnlinkDevice] validation failed: %v", err)
		h.writeResponse(w, http.StatusBadRequest, h.getMessageMap(NONE, GENERAL, INVALID_INPUT))
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
	h.writeResponse(w, http.StatusOK, h.getMessageMap(NONE, UNLINK, UNLINK_DECISION_SUCCESS))
}
