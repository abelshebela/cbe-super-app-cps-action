// Package unlink_device_handler handles the unlink functionality
package unlink_device_handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/unlink"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"github.com/go-chi/chi/v5"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, 0, nil)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UnlinkDevice] validation failed: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidInput, 0, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	cpsAction := entities.CPSAction{
		MakerID:          userContext.UserID,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		RequestAction:    entities.UnlinkDevice,
		MakerActionTime:  time.Now(),
	}

	action_code, err := h.unlinkService.UnlinkDevice(request.UserCode, cpsAction)
	if err != nil {
		h.logger.Errorf("[UnlinkDevice] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[UnlinkDevice] success for userCode: %s", request.UserCode)
	common_util.BaseResponseMaker(map[string]interface{}{"action_code": action_code}, w, "Device unlink request processed successfully", 200)
}

func (h *UnlinkHandler) ApproveUnlinkDevice(w http.ResponseWriter, r *http.Request) {
	var request unlink.ApproveUnlinkDeviceRequest
	action_code := chi.URLParam(r, "action_code")

	request.UserCode = action_code
	if r.Method == "GET" {
		request.Decision = "AUTHORIZED"
	} else {
		request.Decision = "DENIED"
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Errorf("[ApproveUnlinkDevice] failed to decode request: %v", err)
			common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, 0, nil)
			return
		}
		if err := request.Validate(); err != nil {
			h.logger.Warnf("[ApproveUnlinkDevice] validation failed: %v", err)
			common_util.SendErrorResponse(w, common_util.InvalidInput, 0, nil)
			return
		}
	}

	userContext := ctx_util.ExtractUserContext(r)
	cpsAction := entities.CPSAction{
		CheckerID:          userContext.UserID,
		CheckerName:        userContext.FullName,
		CheckerPhoneNumber: userContext.PhoneNumber,
		Department:         userContext.Department,
	}

	if err := h.unlinkService.ApproveOrDecline(request.UserCode, request.Decision, request.Reason, cpsAction); err != nil {
		h.logger.Errorf("[ApproveUnlinkDevice] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[ApproveUnlinkDevice] decision: %s for userCode: %s", request.Decision, request.UserCode)
	common_util.BaseResponseMaker(nil, w, "Unlink decision processed successfully", 200)
}
