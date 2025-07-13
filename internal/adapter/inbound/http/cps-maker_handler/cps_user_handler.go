package cpsmakerhandler

import (
	"encoding/json"

	"net/http"
	"strings"

	cpsapp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_user_maker"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	util_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CPSUserMakerHandler struct {
	Service cpsapp.ApplicationService
	logger  utils.Logger
}

func InitCPSUserMakerHandler(service cpsapp.ApplicationService, logger utils.Logger) inbound.CPSUserMakerHandler {
	return CPSUserMakerHandler{
		Service: service,
		logger:  logger,
	}
}

func (h CPSUserMakerHandler) CreateUserRequest(w http.ResponseWriter, r *http.Request) {
	dataCPSAction, err := h.Service.CreateUserRequest(r.Context(), r)
	if err != nil {
		h.logger.Errorf("CreateUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := local_util.StructToMap(dataCPSAction)
	if err != nil {
		h.logger.Errorf("failed to convert data to map: %v", err)
		local_util.SendErrorResponse(w, "Failed to convert data to map", http.StatusInternalServerError, nil)
	}
	local_util.BaseResponseMaker(data, w, "User request submitted successfully", 200)
}

func (h CPSUserMakerHandler) UpdateUserRequest(w http.ResponseWriter, r *http.Request) {
	dataCPSAction, err := h.Service.UpdateUserRequest(r.Context(), r)
	if err != nil {
		h.logger.Errorf("UpdateUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := local_util.StructToMap(dataCPSAction)
	if err != nil {
		h.logger.Errorf("failed to convert data to map: %v", err)
		local_util.SendErrorResponse(w, "Failed to convert data to map", http.StatusInternalServerError, nil)
	}
	local_util.BaseResponseMaker(data, w, "User update request processed successfully", 200)
}

func (h CPSUserMakerHandler) ApproveUserAction(w http.ResponseWriter, r *http.Request) {
	dataCPSAction, err := h.Service.ApproveUserAction(r.Context(), r)
	if err != nil {
		h.logger.Errorf("UpdateUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	if dataCPSAction == nil {
		local_util.BaseResponseMaker(nil, w, "User update request rejected successfully", 200)
		return
	}

	data, err := local_util.StructToMap(dataCPSAction)
	if err != nil {
		h.logger.Errorf("failed to convert data to map: %v", err)
		local_util.SendErrorResponse(w, "Failed to convert data to map", http.StatusInternalServerError, nil)
	}
	local_util.BaseResponseMaker(data, w, "User update request approved successfully", 200)
}

func (h CPSUserMakerHandler) GetPendingUserActions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actionCode := r.URL.Query().Get("action_code")
	actions, err := h.Service.GetPendingUserActions(ctx, actionCode)
	if err != nil {
		h.logger.Errorf("GetPendingUserActions failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := local_util.StructToMap(actions)
	local_util.BaseResponseMaker(data, w, def.Message, http.StatusOK)
}

func (h CPSUserMakerHandler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserCode string `json:"user_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.UserCode) == "" {
		h.logger.Errorf("user_code is required or failed to decode: %v", err)
		local_util.SendErrorResponse(w, "user_code is required", http.StatusBadRequest, nil)
		return
	}

	ctx := r.Context()
	user, err := h.Service.FetchUserByUserCode(ctx, req.UserCode)
	if err != nil {
		h.logger.Errorf("FetchUserByUserCode failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := local_util.StructToMap(user)
	local_util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}
