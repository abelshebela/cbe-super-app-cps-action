package actionrole_handler

import (
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	actionrole_inbound "cbe-super-app-cps-action/internal/constants/interfaces/action_role"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type BPSActionRoleHandler struct {
	service service.BPSActionRoleService
	logger  utils.Logger
}

func NewBPSActionRoleHandler(svc service.BPSActionRoleService, logger utils.Logger) actionrole_inbound.BPSActionRoleHandler {
	return &BPSActionRoleHandler{service: svc, logger: logger}
}

// GetAll godoc
//
//	@Summary	List action roles
//	@Tags		ActionRole
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int		false	"Page"
//	@Param		per_page	query		int		false	"Per Page"
//	@Param		search		query		string	false	"Search"
//	@Success	200			{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles [get]
func (h *BPSActionRoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllActionRoles", "handler", "actionRole")
	defer span.End()
	filter := *local_util.ExtractFilterParams(r)
	res, err := h.service.FindAllWithPagination(ctx, filter)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("action_role.count", len(res.Data)))
	h.logger.Infof("[GetAll] retrieved %d action roles", len(res.Data))
	localization.SendSuccessResponse(w, localization.SuccessActionRolesFetched, res)
}

// GetByActionCode godoc
//
//	@Summary	Get action role by code
//	@Tags		ActionRole
//	@Produce	json
//	@Param		code	path		string	true	"Action Code"
//	@Success	200		{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles/{code} [get]
func (h *BPSActionRoleHandler) GetByActionCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getActionRoleByCode", "handler", "actionRole")
	defer span.End()
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("action_role.code", code))
	res, err := h.service.GetByActionCode(ctx, code)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetByActionCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[GetByActionCode] action role retrieved successfully for code: %s", code)
	localization.SendSuccessResponse(w, localization.SuccessActionRoleFetched, res)
}

// Create godoc
//
//	@Summary	Create action role (maker)
//	@Tags		ActionRole
//	@Accept		json
//	@Produce	json
//	@Param		body	body		actionrole_dto.CreateActionRoleRequest	true	"Create"
//	@Success	201		{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles [post]
func (h *BPSActionRoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createActionRole", "handler", "actionRole")
	defer span.End()
	var req actionrole_dto.CreateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Create] failed to decode request: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}
	if req.ActionName == "" {
		localization.SendBadRequestResponse(w, localization.ErrorActionNameIsRequired.Message)
		return
	}
	span.SetAttributes(
		attribute.String("action_role.code", req.ActionCode),
		attribute.String("action_role.name", req.ActionName),
	)
	err := h.service.Create(ctx, req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Create] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[Create] request sent successfully for action_code: %s", req.ActionCode)
	localization.SendSuccessResponse(w, localization.SuccessActionRoleCreateRequestCreated, nil)
}

// Update godoc
//
//	@Summary	Update action role (maker)
//	@Tags		ActionRole
//	@Accept		json
//	@Produce	json
//	@Param		code	path		string									true	"Action Code"
//	@Param		body	body		actionrole_dto.UpdateActionRoleRequest	true	"Update"
//	@Success	201		{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles/{code} [patch]
func (h *BPSActionRoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateActionRole", "handler", "actionRole")
	defer span.End()
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	var req actionrole_dto.UpdateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Update] failed to decode request: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}
	span.SetAttributes(attribute.String("action_role.code", code))
	err := h.service.Update(ctx, code, req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Update] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[Update] request sent successfully for code: %s", code)
	localization.SendSuccessResponse(w, localization.SuccessActionRoleUpdateRequestCreated, nil)
}

// Enable godoc
//
//	@Summary	Enable action role (maker)
//	@Tags		ActionRole
//	@Produce	json
//	@Param		code	path		string	true	"Action Code"
//	@Success	201		{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles/{code}/enable [patch]
func (h *BPSActionRoleHandler) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableActionRole", "handler", "actionRole")
	defer span.End()
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("action_role.code", code))
	if err := h.service.Enable(ctx, code); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[Enable] request sent successfully for code: %s", code)
	localization.SendSuccessResponse(w, localization.SuccessActionRoleEnableRequestCreated, nil)
}

// Disable godoc
//
//	@Summary	Disable action role (maker)
//	@Tags		ActionRole
//	@Produce	json
//	@Param		code	path		string	true	"Action Code"
//	@Success	201		{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles/{code}/disable [patch]
func (h *BPSActionRoleHandler) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableActionRole", "handler", "actionRole")
	defer span.End()
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("action_role.code", code))
	if err := h.service.Disable(ctx, code); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("[Disable] request sent successfully for code: %s", code)
	localization.SendSuccessResponse(w, localization.SuccessActionRoleDisableRequestCreated, nil)
}
