// Package unlink_device_handler handles the unlink functionality
package unlink_device_handler

import (
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/unlink"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func (h *UnlinkHandler) UnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.UnlinkDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UnlinkDevice] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UnlinkDevice] validation failed: %v", err)
		common_util.SendErrorResponse(w, InvalidInput, http.StatusBadRequest, nil)

		return
	}

	ctx := ctx_util.ExtractUserContext(r)

	userID := ctx.UserID
	branchCode := ctx.BranchCode
	department := ctx.Department

	if err := h.unlinkApp.UnlinkDevice(request.UserCode, userID, branchCode, department); err != nil {
		h.logger.Errorf("[UnlinkDevice] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.logger.Infof("[UnlinkDevice] success for userCode: %s", request.UserCode)
	common_util.WriteSuccessResponse(w, nil, "Device unlink request processed successfully")
}

func (h *UnlinkHandler) ApproveUnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.ApproveUnlinkDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, InvalidJSONPayload, http.StatusBadRequest, nil)

		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[ApproveUnlinkDevice] validation failed: %v", err)
		common_util.SendErrorResponse(w, InvalidInput, http.StatusBadRequest, nil)

		return
	}
	ctx := ctx_util.ExtractUserContext(r)
	userID := ctx.UserID

	if err := h.unlinkApp.ApproveOrDecline(request.UserCode, request.Decision, request.Reason, userID); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.logger.Infof("[ApproveUnlinkDevice] decision: %s for userCode: %s", request.Decision, request.UserCode)
	common_util.WriteSuccessResponse(w, nil, "Unlink decision processed successfully")
}
