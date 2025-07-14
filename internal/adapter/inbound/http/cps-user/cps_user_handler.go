package cpsmakerhandler

import (
	"net/http"

	cpsapp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
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
		return
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
		return
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
		return
	}
	local_util.BaseResponseMaker(data, w, "User update request approved successfully", 200)
}

func (h CPSUserMakerHandler) GetPendingUserActions(w http.ResponseWriter, r *http.Request) {
	userActions, err := h.Service.GetPendingUserActions(r.Context())
	if err != nil {
		h.logger.Errorf("GetPendingUserAction failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := local_util.StructToMap(userActions)
	if err != nil {
		local_util.SendErrorResponse(w, "UNHANDLER_SERVER_ERROR", 500, nil)
		return
	}
	local_util.BaseResponseMaker(data, w, "Pending users fetched successfully", 200)
}

func (h CPSUserMakerHandler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	user, err := h.Service.FetchUserByUserCode(r.Context(), r)
	if err != nil {
		h.logger.Errorf("FetchUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := local_util.StructToMap(user)
	if err != nil {
		h.logger.Errorf("failed to convert data to map: %v", err)
		local_util.SendErrorResponse(w, "Failed to convert data to map", http.StatusInternalServerError, nil)
		return
	}
	local_util.BaseResponseMaker(data, w, "CPS user fetched successfully", 200)
}
