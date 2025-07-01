package passwordrule

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cbe-super-app-cps-action/internal/application/middleware"
	passwordrule "cbe-super-app-cps-action/internal/application/password_rule"
	"cbe-super-app-cps-action/internal/domain/action"
	"cbe-super-app-cps-action/internal/domain/password_rule/services"
	constant "cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PasswordRuleHTTPHandler struct {
	service services.PasswordRuleService
	logger  utils.Logger
}

func NewPasswordRuleHTTPHandler(service services.PasswordRuleService, logger utils.Logger) *PasswordRuleHTTPHandler {
	return &PasswordRuleHTTPHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PasswordRuleHTTPHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
	var req passwordrule.RequestPasswordRuleUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind password rule update data: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           map[string]string{"message": "invalid request"},
		}
		resp.SendJSON()
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusUnauthorized,
			Data:           map[string]string{"message": "user info missing in context"},
		}
		resp.SendJSON()
		return
	}

	if strings.TrimSpace(department) == "" {
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusUnauthorized,
			Data:           map[string]string{"message": "department is required in context"},
		}
		resp.SendJSON()
		return
	}

	ctx := context.WithValue(r.Context(), "department", department)

	maker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	_, err := h.service.RequestPasswordRuleUpdate(ctx, &req.Rule, maker)
	if err != nil {
		h.logger.Errorf("RequestPasswordRuleUpdate failed: %v", err)
		status := http.StatusInternalServerError
		msg := err.Error()
		if strings.Contains(msg, "pending action already exists") {
			status = http.StatusConflict
			msg = "A pending action already exists for this maker"
		}
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         status,
			Data:           map[string]string{"message": msg},
		}
		resp.SendJSON()
		return
	}

	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusCreated,
		Data:           map[string]string{"message": "Action created successfully"},
	}
	resp.SendJSON()
}
func (h *PasswordRuleHTTPHandler) ApproveOrRejectPasswordRuleAction(w http.ResponseWriter, r *http.Request) {
	var req passwordrule.ApproveOrRejectPasswordRuleActionDTO
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

	checker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	ctx := r.Context()
	err := h.service.ApproveOrRejectPasswordRuleAction(ctx, req.ActionID, req.Approve, checker, req.RejectionReason)
	if err != nil {
		h.logger.Errorf("ApproveOrRejectPasswordRuleAction failed: %v", err)
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

func (h *PasswordRuleHTTPHandler) GetPasswordRuleUpdateActionByID(w http.ResponseWriter, r *http.Request) {
	actionID := r.URL.Query().Get("action_id")
	if actionID == "" {
		err := fmt.Errorf("action_id is required %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "action_id is required",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	userRole, _ := r.Context().Value(constant.ContextKey("user_role")).(string)
	if userRole != "CHECKER" && userRole != "IFB-CHECKER" {
		err := fmt.Errorf("forbidden: only checker can access this resource %w", constant.ErrorDefinition{
			Code:    http.StatusForbidden,
			Message: "forbidden: only checker can access this resource",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	ctx := r.Context()
	result, err := h.service.GetPasswordRuleUpdateActionByID(ctx, actionID)
	if err != nil {
		h.logger.Errorf("GetPasswordRuleUpdateActionByID failed: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}
	res := common.Response[*action.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           result,
	}
	res.SendJSON()
}
func (h *PasswordRuleHTTPHandler) GetUpdateAction(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	userCode, _ := r.Context().Value(constant.ContextKey("user_code")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(userCode) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" || strings.TrimSpace(department) == "" {
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           map[string]string{"message": "missing required user or department info"},
		}
		resp.SendJSON()
		return
	}

	maker := action.User{
		UserID:      userID,
		UserCode:    userCode,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	ctx := context.WithValue(r.Context(), "department", department)
	result, err := h.service.GetUpdateAction(ctx, maker)
	if err != nil {
		h.logger.Errorf("GetUpdateAction failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusInternalServerError,
			Data:           map[string]string{"message": err.Error()},
		}
		resp.SendJSON()
		return
	}

	res := common.Response[*action.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           result,
	}
	res.SendJSON()
}

type CheckPasswordDTO struct {
	Password string `json:"password"`
}

func (h *PasswordRuleHTTPHandler) CheckPasswordRule(w http.ResponseWriter, r *http.Request) {
	var req CheckPasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Password) == "" {
		resp := map[string]interface{}{"valid": false, "message": "invalid request: password required"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	valid, msg := h.service.CheckPasswordRule(r.Context(), req.Password)
	resp := map[string]interface{}{"valid": valid, "message": msg}
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}
