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
	_, err := h.Service.CreateUserRequest(r.Context(), r)
	if err != nil {
		h.logger.Errorf("CreateUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	// response := map[string]interface{}{"action_code": dataCPSAction.ActionCode}
	// local_util.BaseResponseMaker(response, w, "User request submitted successfully", 200)
	local_util.WriteSuccessResponse(w, nil, "Create user request action submitted successfully")
}

func (h CPSUserMakerHandler) UpdateUserRequest(w http.ResponseWriter, r *http.Request) {
	_, err := h.Service.UpdateUserRequest(r.Context(), r)
	if err != nil {
		h.logger.Errorf("UpdateUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	// response := map[string]interface{}{"action_code": dataCPSAction.ActionCode}
	// local_util.BaseResponseMaker(response, w, "User update request processed successfully", 200)
	local_util.WriteSuccessResponse(w, nil, "Update user request action submitted successfully")
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
	pendingUserAction, err := h.Service.GetPendingUserActions(r.Context())
	if err != nil {
		h.logger.Errorf("GetPendingUserAction failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_util.BaseResponseMaker(pendingUserAction, w, "Pending users fetched successfully", 200)
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

func (h CPSUserMakerHandler) GetAllCPSUsers(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	users, err := h.Service.GetAllCPSUsers(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("GetAllCPSUser request failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	doc, _ := local_util.StructToMap(users)
	local_util.BaseResponseMaker(doc, w, "CPS users fetched successfully", 200)
}

func (h CPSUserMakerHandler) DeleteUserRequest(w http.ResponseWriter, r *http.Request) {
	action, err := h.Service.DeleteUserRequest(r.Context(), r)
	if err != nil {
		h.logger.Errorf("DeleteUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := map[string]interface{}{"action_code": action.ActionCode}
	local_util.BaseResponseMaker(response, w, "User deletion action submitted successfully", 200)
}

func (h CPSUserMakerHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	err := h.Service.DisableUser(r.Context(), r)
	if err != nil {
		h.logger.Errorf("Disable user request failed: %v\n", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_util.BaseResponseMaker(nil, w, "User disable request submitted successfully", 200)
}


func (h CPSUserMakerHandler) EnableUser(w http.ResponseWriter, r *http.Request) {
	err := h.Service.EnableUser(r.Context(), r)
	if err != nil {
		h.logger.Errorf("Enable user request failed: %v\n", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_util.BaseResponseMaker(nil, w, "User Enable request submitted successfully", 200)
}
