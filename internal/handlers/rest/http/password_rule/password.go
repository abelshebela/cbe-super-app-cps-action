package passwordrule

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/password_rule"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/password_rule/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
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

// Get All Password Rules
//
//	@Summary		Get All Password Rules
//	@Description	Retrieves all password rules with pagination
//	@Tags			PasswordRule
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=passwordrule.PaginatedPasswordRulesResponse}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/password_rule/ [get]
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

// Request Password Rule Update
//
//	@Summary		Request Password Rule Update
//	@Description	Requests an update to a password rule
//	@Tags			PasswordRule
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string							true	"Password Rule ID"
//	@Param			body			body		passwordrule.PasswordRuleUpdate	true	"Password Rule Update DTO"
//	@Success		200				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/password_rule/{id} [put]
func (p *passwordRuleHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
	var req dto.PasswordRuleUpdate

	id, ok := core.ExtractID(w, r, p.logger)
	if !ok {
		return
	}

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

// Check Password Rule
//
//	@Summary		Check Password Rule
//	@Description	Checks if a password meets the rule
//	@Tags			PasswordRule
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body			body		passwordrule.CheckPasswordDTO	true	"Check Password DTO"
//	@Success		200				{object}	object{status=int,message=string,data=object{valid=bool}}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/password_rule/check [post]
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
