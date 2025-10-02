package passwordrule

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/password_rule"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/password_rule/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type passwordRuleHandler struct {
	service service.PasswordRuleService
	logger  utils.Logger
}

func InitPasswordRuleHandler(service service.PasswordRuleService, logger utils.Logger) *passwordRuleHandler {
	return &passwordRuleHandler{
		service: service,
		logger:  logger,
	}
}

func (p *passwordRuleHandler) GetPasswordRule(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	passwordRules, err := p.service.GetAllPasswordRules(r.Context(), *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data, _ := core.StructToMap(passwordRules)
	localization.SendSuccessResponse(w, localization.SuccessFetchAllPasswordRules, data)
}

func (p *passwordRuleHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
	var req dto.PasswordRuleUpdate

	id, ok := core.ExtractID(w, r, p.logger)
	if !ok {
		return
	}
	fmt.Println("Body:-----------------", r.Body)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		p.logger.Errorf("Failed to validate incoming password rules | Error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	body := core.PasswordRuleDtoToModel(req)

	err := p.service.RequestPasswordRuleUpdate(r.Context(), id, body)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUpdatePasswordRule, nil)
}

func (p *passwordRuleHandler) CheckPasswordRule(w http.ResponseWriter, r *http.Request) {
	var body dto.CheckPasswordDTO

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Password) == "" {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	valid, msg := p.service.CheckPasswordRule(r.Context(), body.Password)

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
