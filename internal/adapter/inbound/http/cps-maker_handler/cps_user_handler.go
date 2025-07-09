package cpsmakerhandler

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"strings"

	cpsapp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/cps_user_maker"
	// "cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	util_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CPSUserMakerHandler struct {
	Service cpsapp.ApplicationService
	Logger  utils.Logger
}

func InitCPSUserMakerHandler(service cpsapp.ApplicationService, logger utils.Logger) inbound.CPSUserMakerHandler {
	return CPSUserMakerHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h CPSUserMakerHandler) CreateUserRequest(w http.ResponseWriter, r *http.Request) {
	var req cpsapp.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Errorf("failed to bind user data: %v", err)

		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if strings.TrimSpace(req.UserName) == "" ||
		strings.TrimSpace(req.FullName) == "" ||
		strings.TrimSpace(req.PhoneNumber) == "" ||
		strings.TrimSpace(req.UserRole) == "" ||
		strings.TrimSpace(req.Department) == "" {
		util.SendErrorResponse(w, "missing required field", http.StatusBadRequest, nil)
		return
	}

	userID := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber := r.Context().Value(constant.ContextKey("phone_number")).(string)

	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	cpsuser := action.CPSUser{
		UserName:           req.UserName,
		FullName:           req.FullName,
		PhoneNumber:        req.PhoneNumber,
		Role:               req.UserRole,
		Department:         req.Department,
		PermissionCategory: req.PermissionCategory,
		PermissionGroup:    req.PermissionGroups,
	}

	ctx := r.Context()
	if err := h.Service.CreateUserRequest(ctx, cpsuser, maker); err != nil {
		h.Logger.Errorf("CreateUserRequest failed: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(cpsuser)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}
func (h CPSUserMakerHandler) UpdateUserRequest(w http.ResponseWriter, r *http.Request) {

	var req cpsapp.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Errorf("failed to bind user data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if strings.TrimSpace(req.UserCode) == "" {
		util.SendErrorResponse(w, "user_code is required", http.StatusBadRequest, nil)
		return
	}

	userID := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber := r.Context().Value(constant.ContextKey("phone_number")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		util.SendErrorResponse(w, "user info missing in context", http.StatusUnauthorized, nil)
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
	if err := h.Service.UpdateUserRequest(ctx, updated, maker); err != nil {
		h.Logger.Errorf("UpdateUserRequest failed: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(updated)
	util.BaseResponseMaker(data, w, def.Message, http.StatusOK)
}
func (h CPSUserMakerHandler) ApproveUserAction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		h.Logger.Errorf("action_id is required or failed to decode: %v", err)
		util.SendErrorResponse(w, "action_id is required", http.StatusBadRequest, nil)
		return
	}

	userID := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber := r.Context().Value(constant.ContextKey("phone_number")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		util.SendErrorResponse(w, "user info missing in context", http.StatusUnauthorized, nil)
		return
	}

	ctx := r.Context()
	if err := h.Service.ApproveUserAction(ctx, req.ActionID, req.Approve, req.Reason); err != nil {
		h.Logger.Errorf("ApproveUserAction failed: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	util.BaseResponseMaker(nil, w, def.Message, http.StatusAccepted)
}

func (h CPSUserMakerHandler) GetPendingUserActions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actionCode := r.URL.Query().Get("action_code")
	actions, err := h.Service.GetPendingUserActions(ctx, actionCode)
	if err != nil {
		h.Logger.Errorf("GetPendingUserActions failed: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(actions)
	util.BaseResponseMaker(data, w, def.Message, http.StatusOK)
}

func (h CPSUserMakerHandler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserCode string `json:"user_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.UserCode) == "" {
		h.Logger.Errorf("user_code is required or failed to decode: %v", err)
		util.SendErrorResponse(w, "user_code is required", http.StatusBadRequest, nil)
		return
	}

	ctx := r.Context()
	user, err := h.Service.FetchUserByUserCode(ctx, req.UserCode)
	if err != nil {
		h.Logger.Errorf("FetchUserByUserCode failed: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := util_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(user)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}
