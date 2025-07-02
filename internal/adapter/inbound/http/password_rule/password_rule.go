package passwordrule

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	passwordrule "cbe-super-app-cps-action/internal/application/password_rule"
	"cbe-super-app-cps-action/internal/domain/action"
	"cbe-super-app-cps-action/internal/domain/password_rule/services"
	"cbe-super-app-cps-action/pkgs/utils"
)

type PasswordRuleHTTPHandler struct {
	service services.PasswordRuleService
}

func NewPasswordRuleHTTPHandler(service services.PasswordRuleService) *PasswordRuleHTTPHandler {
	return &PasswordRuleHTTPHandler{
		service: service,
	}
}

func (h *PasswordRuleHTTPHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
	var req passwordrule.RequestPasswordRuleUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}
	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userID, _ := r.Context().Value("user_id").(string)
	fullName, _ := r.Context().Value("full_name").(string)
	phoneNumber, _ := r.Context().Value("phone_number").(string)
	department, _ := r.Context().Value("department").(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	if strings.TrimSpace(department) == "" {
		utils.SendErrorResponse(w, "INVALID_INPUT", 0, map[string]interface{}{"message": "department is required in context"})
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
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Action created successfully"})
}

func (h *PasswordRuleHTTPHandler) ApproveOrRejectPasswordRuleAction(w http.ResponseWriter, r *http.Request) {
	var req passwordrule.ApproveOrRejectPasswordRuleActionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}
	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userID, _ := r.Context().Value("user_id").(string)
	fullName, _ := r.Context().Value("full_name").(string)
	phoneNumber, _ := r.Context().Value("phone_number").(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
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
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Action processed successfully"})
}

func (h *PasswordRuleHTTPHandler) GetPasswordRuleUpdateActionByID(w http.ResponseWriter, r *http.Request) {
	actionID := r.URL.Query().Get("action_id")
	if actionID == "" {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"message": "action_id is required"})
		return
	}

	userRole, _ := r.Context().Value("user_role").(string)
	if userRole != "CHECKER" && userRole != "IFB-CHECKER" {
		utils.SendErrorResponse(w, "ACTION_NOT_ALLOWED", http.StatusForbidden, map[string]interface{}{"message": "forbidden: only checker can access this resource"})
		return
	}

	ctx := r.Context()
	result, err := h.service.GetPasswordRuleUpdateActionByID(ctx, actionID)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *PasswordRuleHTTPHandler) GetUpdateAction(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(string)
	userCode, _ := r.Context().Value("user_code").(string)
	fullName, _ := r.Context().Value("full_name").(string)
	phoneNumber, _ := r.Context().Value("phone_number").(string)
	department, _ := r.Context().Value("department").(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(userCode) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" || strings.TrimSpace(department) == "" {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"message": "missing required user or department info"})
		return
	}

	maker := action.User{
		UserID:      userID,
		UserCode:    userCode,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	ctx := r.Context()
	result, err := h.service.GetUpdateAction(ctx, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

type CheckPasswordDTO struct {
	Password string `json:"password"`
}

func (h *PasswordRuleHTTPHandler) CheckPasswordRule(w http.ResponseWriter, r *http.Request) {
	var req CheckPasswordDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Password) == "" {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"message": "invalid request: password required"})
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
