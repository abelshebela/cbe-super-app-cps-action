package cpsmakerhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cpsapp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/cps_user_maker"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
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
		err = fmt.Errorf("failed to bind error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	if strings.TrimSpace(req.UserName) == "" ||
		strings.TrimSpace(req.FullName) == "" ||
		strings.TrimSpace(req.PhoneNumber) == "" ||
		strings.TrimSpace(req.UserRole) == "" ||
		strings.TrimSpace(req.Department) == "" {
		err := fmt.Errorf("missing required fields %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing required fields",
		})
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[string]{
		ResponseWriter: w,
		Status:         http.StatusCreated,
		Data:           "User request created and pending approval",
	}
	res.SendJSON()
}

func (h CPSUserMakerHandler) UpdateUserRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserCode string         `json:"user_code"`
		Updated  action.CPSUser `json:"updated"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.UserCode) == "" {
		err := fmt.Errorf("user_code is required %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "user_code is required",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		err := fmt.Errorf("user info missing in context %w", constant.ErrorDefinition{
			Code:    http.StatusUnauthorized,
			Message: "user info missing in context",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	ctx := r.Context()
	if err := h.Service.UpdateUserRequest(ctx, req.Updated, maker); err != nil {
		h.Logger.Errorf("UpdateUserRequest failed: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[string]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           "User updated",
	}
	res.SendJSON()
}

func (h CPSUserMakerHandler) ApproveUserAction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActionID string  `json:"action_id"`
		Approve  bool    `json:"approve"`
		Reason   *string `json:"reason,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		err := fmt.Errorf("action_id is required %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "action_id is required",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		err := fmt.Errorf("user info missing in context %w", constant.ErrorDefinition{
			Code:    http.StatusUnauthorized,
			Message: "user info missing in context",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	ctx := r.Context()
	if err := h.Service.ApproveUserAction(ctx, req.ActionID, req.Approve, req.Reason); err != nil {
		h.Logger.Errorf("ApproveUserAction failed: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[string]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           "Action processed successfully",
	}
	res.SendJSON()
}

func (h CPSUserMakerHandler) GetPendingUserActions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actions, err := h.Service.GetPendingUserActions(ctx)
	if err != nil {
		h.Logger.Errorf("GetPendingUserActions failed: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}
	res := common.Response[[]action.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           actions,
	}
	res.SendJSON()
}
