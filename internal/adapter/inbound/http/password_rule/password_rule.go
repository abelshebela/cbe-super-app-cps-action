package passwordrule

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/password_rule"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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
	filterParams := utils.ExtractFilterParams(r)
	ctx := r.Context()

	passwordRules, err := h.service.GetAllPasswordRules(ctx, filterParams)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, _ := utils.StructToMap(passwordRules)
	utils.BaseResponseMaker(data, w, "Password Rules Successfully Fetched", 200)
}
func (h *PasswordRuleHTTPHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
	var req passwordrule.RequestPasswordRuleUpdateDTO
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}

	req.ID = id

	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	ctx := context.WithValue(r.Context(), constant.ContextKey("department"), userContext.Department)

	maker := action.User{
		UserID:      userContext.UserID,
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}

	req.Rule.ID = req.ID
	actionID, err := h.service.RequestPasswordRuleUpdate(ctx, &req.Rule, maker, userContext.Department)

	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, map[string]interface{}{"action_id": actionID}, "Update request submitted for approval")
}

func (h *PasswordRuleHTTPHandler) GetPasswordRuleUpdateActionByID(w http.ResponseWriter, r *http.Request) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	ctx := r.Context()
	result, err := h.service.GetPasswordRuleUpdateActionByID(ctx, actionCode)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, result, "Fetched password rule update action successfully")

}

func (h *PasswordRuleHTTPHandler) GetUpdateAction(w http.ResponseWriter, r *http.Request) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	maker := action.User{
		UserID:      userContext.UserID,
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}

	ctx := r.Context()
	result, err := h.service.GetUpdateAction(ctx, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, result, "Fetched update action successfully")

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
		status = http.StatusBadRequest
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
