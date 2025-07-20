package passwordrule

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/password_rule"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type PasswordRuleHTTPHandler struct {
	service services.PasswordRuleService
}

func NewPasswordRuleHTTPHandler(service services.PasswordRuleService) *PasswordRuleHTTPHandler {
	return &PasswordRuleHTTPHandler{
		service: service,
	}
}

func (h *PasswordRuleHTTPHandler) GetPasswordRule(w http.ResponseWriter, r *http.Request) {

	passwordRules, err := h.service.GetAllPasswordRules(r.Context())

	fmt.Printf("Password rule update request took \n")

	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data := map[string]interface{}{"password_rules": passwordRules}
	utils.BaseResponseMaker(data, w, "Password Rules Successfuly Fetched", 200)

}
func (h *PasswordRuleHTTPHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
	var req passwordrule.RequestPasswordRuleUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}
	fmt.Println(req, "required")
	if err := req.Validate(); err != nil {
		fmt.Println("HOLAAA", err)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	userCode, _ := r.Context().Value(constant.ContextKey("user_code")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(userCode) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" || strings.TrimSpace(department) == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	ctx := context.WithValue(r.Context(), constant.ContextKey("department"), department)

	maker := action.User{
		UserID:      userID,
		UserCode:    userCode,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	req.Rule.ID = req.ID

	start := time.Now()
	actionID, err := h.service.RequestPasswordRuleUpdate(ctx, &req.Rule, maker, department)
	elapsed := time.Since(start)
	fmt.Printf("Password rule update request took %s\n", elapsed)

	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := map[string]interface{}{
		"status":  200,
		"message": "Update request submitted for approval",
		"data":    map[string]interface{}{"action_id": actionID},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// func (h *PasswordRuleHTTPHandler) ApproveOrRejectPasswordRuleAction(w http.ResponseWriter, r *http.Request) {
// 	var req passwordrule.ApproveOrRejectPasswordRuleActionDTO
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
// 		return
// 	}
// 	fmt.Println(req, "required")
// 	if err := req.Validate(); err != nil {
// 		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
// 		return
// 	}

// 	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
// 	userCode, _ := r.Context().Value(constant.ContextKey("user_code")).(string)
// 	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
// 	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
// 	department, _ := r.Context().Value(constant.ContextKey("department")).(string)
// 	fmt.Println(userID, userCode, fullName, phoneNumber, department, "nodjghdjkfgheidgjfdk")
// 	if strings.TrimSpace(userID) == "" || strings.TrimSpace(userCode) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" || strings.TrimSpace(department) == "" {
// 		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
// 		return
// 	}

// 	checker := action.User{
// 		UserID:      userID,
// 		UserCode:    userCode,
// 		FullName:    fullName,
// 		PhoneNumber: phoneNumber,
// 	}

// 	ctx := r.Context()
// 	err := h.service.ApproveOrRejectPasswordRuleAction(ctx, req.ActionID, req.Decision, checker, req.RejectionReason, department)
// 	if err != nil {
// 		utils.SendErrorResponse(w, err.Error(), 0, nil)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"status":  200,
// 		"message": "Action processed successfully",
// 		"data":    map[string]interface{}{"action_id": req.ActionID},
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	_ = json.NewEncoder(w).Encode(response)
// }

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
	response := map[string]interface{}{
		"status":  200,
		"message": "Fetched password rule update action successfully",
		"data":    result,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
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
	response := map[string]interface{}{
		"status":  200,
		"message": "Fetched update action successfully",
		"data":    result,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
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

	statusCode := http.StatusOK
	status := 200
	if !valid {
		status = 200
		statusCode = http.StatusBadRequest
	}
	response := map[string]interface{}{
		"status":  status,
		"message": msg,
		"data":    map[string]interface{}{"valid": valid},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *PasswordRuleHTTPHandler) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
