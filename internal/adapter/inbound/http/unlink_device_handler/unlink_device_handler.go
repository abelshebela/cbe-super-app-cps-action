package unlink_device_handler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/unlink"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/unlink"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
    constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type UnlinkHandler struct {
	unlinkService unlink.ApplicationService
	logger        utils.Logger
}

func NewHTTPUnlinkHandler(service unlink.ApplicationService, logger utils.Logger) inbound.UnlinkPortHandler {
	return &UnlinkHandler{
		unlinkService: service,
		logger:        logger,
	}
}

func (h *UnlinkHandler) UnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.UnlinkDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UnlinkDevice] failed to decode request: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message":  "Invalid JSON payload"},
		}
		resp.SendJSON()
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UnlinkDevice] validation failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message": "Invalid input provided"},
		}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	branchCode, _ := r.Context().Value(constant.ContextKey("branch_code")).([]string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if err := h.unlinkService.UnlinkDevice(request.UserCode, userID, branchCode, department); err != nil {
		h.logger.Errorf("[UnlinkDevice] service error: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{
			"message": err.Error(),
		},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[UnlinkDevice] success for userCode: %s", request.UserCode)
	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data: map[string]string{
			"message": "Device unlink request processed successfully",
		},
	}
	resp.SendJSON()
}

func (h *UnlinkHandler) ApproveUnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.ApproveUnlinkDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] failed to decode request: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message": "Invalid JSON payload"},
		}
		resp.SendJSON()
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[ApproveUnlinkDevice] validation failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message": "Invalid input provided"},
		}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)

	if err := h.unlinkService.ApproveOrDecline(request.UserCode, request.Decision, request.Reason, userID); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] service error: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{
			"message": err.Error(),
		},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[ApproveUnlinkDevice] decision: %s for userCode: %s", request.Decision, request.UserCode)
	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data: map[string]string{
			"message": "Unlink decision processed successfully",
		},
	}
	resp.SendJSON()
}
