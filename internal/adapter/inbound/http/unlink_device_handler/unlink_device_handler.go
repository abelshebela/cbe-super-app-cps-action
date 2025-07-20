// Package unlink_device_handler handles the unlink functionality
package unlink_device_handler

import (
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/unlink"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

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
	UserCode, ok := common_util.GetParam(r, "user_code")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'user_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
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

	action_code, err := h.unlinkService.UnlinkDevice(r.Context(), UserCode, cpsAction)
	if err != nil {
		h.logger.Errorf("[UnlinkDevice] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[UnlinkDevice] success for userCode: %s", UserCode)
	common_util.BaseResponseMaker(map[string]interface{}{"action_code": action_code}, w, "Device unlink request processed successfully", 200)
}
