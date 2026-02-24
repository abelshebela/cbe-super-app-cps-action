package passwordrule

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/password_rule"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/password_rule/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	constants "cbe-super-app-cps-action/internal/constants"
	types "cbe-super-app-cps-action/internal/constants/types"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "passwordRule", "passwordRuleHandler", "GetPasswordRule")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, p.logger)
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	passwordRules, err := p.service.GetAllPasswordRules(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[GetPasswordRule] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Password rules retrieved", trace.WithAttributes(attribute.Int("count", len(passwordRules.Data))))
	log.Infof("[GetPasswordRule] retrieved %d password rules", len(passwordRules.Data))
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
//	@Router			/password_rule/{id} [patch]
func (p *passwordRuleHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "passwordRule", "passwordRuleHandler", "RequestPasswordRuleUpdate")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, p.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	var req dto.PasswordRuleUpdate

	id, ok := core.ExtractID(w, r, p.logger)
	if !ok {
		span.AddEvent("Missing id param")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.AddEvent("Failed to decode request", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[RequestPasswordRuleUpdate] failed to decode request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.AddEvent("Validation failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[PwdRuleH][Create] validate err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	err := p.service.RequestPasswordRuleUpdate(ctx, id, req)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[RequestPasswordRuleUpdate] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		span.AddEvent("PasswordRuleUpdate request sent", trace.WithAttributes(attribute.String("id", id)))
		log.Infof("[RequestPasswordRuleUpdate] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessUpdatePasswordRuleSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessUpdatePasswordRule, nil)
	}
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "passwordRule", "passwordRuleHandler", "CheckPasswordRule")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, p.logger)
	var body dto.CheckPasswordDTO

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Password) == "" {
		span.AddEvent("Failed to decode request or empty password", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[CheckPasswordRule] failed to decode request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	valid, msg := p.service.CheckPasswordRule(ctx, body.Password)
	// Note: Password is not logged for security reasons

	statusCode := http.StatusOK
	status := 200

	if !valid {
		span.AddEvent("Password validation failed", trace.WithAttributes(attribute.String("message", msg)))
		status = http.StatusBadRequest
		statusCode = http.StatusBadRequest
	}
	response := map[string]interface{}{
		"status":  status,
		"message": msg,
		"data":    map[string]interface{}{"valid": valid},
	}
	log.Infof("[CheckPasswordRule] password validation completed, valid: %v", valid)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}
