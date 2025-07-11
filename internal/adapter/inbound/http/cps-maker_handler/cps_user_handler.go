package cpsmakerhandler

import (
	"encoding/json"
	"time"

	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cpsapp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_user_maker"
	cpsUserDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	util_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

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
	local_util.BaseResponseMaker(data, w, "User request submitted for approval", 200)
}

func (h CPSUserMakerHandler) UpdateUserRequest(w http.ResponseWriter, r *http.Request) {

	var req cpsapp.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind user data: %v", err)
		local_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if strings.TrimSpace(req.UserCode) == "" {
		local_util.SendErrorResponse(w, "user_code is required", http.StatusBadRequest, nil)
		return
	}

	userID := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber := r.Context().Value(constant.ContextKey("phone_number")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		local_util.SendErrorResponse(w, "user info missing in context", http.StatusUnauthorized, nil)
		return
	}

	updated := action.CPSUser{
		UserCode:           req.UserCode,
		UserName:           req.UserName,
		FullName:           req.FullName,
		PhoneNumber:        req.PhoneNumber,
		Role:               req.UserRole,
		Department:         req.Department,
		PermissionCategory: req.PermissionCategory,
		PermissionGroup:    req.PermissionGroups,
	}

	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	ctx := r.Context()
	dataCPSAction, err := h.Service.UpdateUserRequest(ctx, updated, maker)
	if err != nil {
		h.logger.Errorf("UpdateUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := local_util.StructToMap(dataCPSAction)
	local_util.BaseResponseMaker(data, w, def.Message, http.StatusOK)
}

func (h CPSUserMakerHandler) ApproveUserAction(w http.ResponseWriter, r *http.Request) {
	actionID := strings.TrimSpace(chi.URLParam(r, "action_id"))
	if actionID == "" {
		h.logger.Errorf("action_id is required")
		local_util.SendErrorResponse(w, "action_id is required", http.StatusBadRequest, nil)
		return
	}

	var req cpsUserDTO.ApproveUserActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		local_util.SendErrorResponse(w, local_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("validattion failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	userPayload := ctx_util.ExtractUserContext(r)
	cpsAction := model.CPSAction{
		ActionCode:         actionID,
		CheckerID:          userPayload.UserID,
		CheckerName:        userPayload.FullName,
		CheckerPhoneNumber: userPayload.PhoneNumber,
		Department:         userPayload.Department,
		RejectionReason:    *req.Reason,
		CheckerActionTime:  time.Now(),
	}

	if err := h.Service.ApproveUserAction(r.Context(), cpsAction); err != nil {
		h.logger.Errorf("ApproveUserAction failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_util.BaseResponseMaker(nil, w, "Action approved successfully", 200)
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
