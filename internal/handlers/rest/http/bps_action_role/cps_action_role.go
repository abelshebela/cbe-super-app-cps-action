package cps_actionrole_handler

import (
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	actionrole_inbound "cbe-super-app-cps-action/internal/constants/interfaces/action_role"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"

	constants "cbe-super-app-cps-action/internal/constants"
	types "cbe-super-app-cps-action/internal/constants/types"
)

type BPSActionRoleHandler struct {
	service service.BPSActionRoleService
	logger  utils.Logger
}

func NewBPSActionRoleHandler(svc service.BPSActionRoleService, logger utils.Logger) actionrole_inbound.BPSActionRoleHandler {
	return &BPSActionRoleHandler{service: svc, logger: logger}
}

// GetAllActionList godoc
//
//	@Summary	List of Action for Cps
//	@Tags		ActionList
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int		false	"Page"
//	@Param		per_page	query		int		false	"Per Page"
//	@Param		search		query		string	false	"Search"
//	@Success	200			{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles/action-list [get]
func (h *BPSActionRoleHandler) GetAllActionList(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllBpsActionRoles", "handler", "bpsActionRole")
	defer span.End()
	filterParams := *local_util.ExtractFilterParams(r)

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

	res, err := h.service.FindAllActionListWithPagination(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("list action list error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("cps_action_role.count", len(res.Data)))
	localization.SendSuccessResponse(w, localization.SuccessActionRolesFetched, res)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllBpsActionRoles", "handler", "bpsActionRole")
	defer span.End()
	filter := *local_util.ExtractFilterParams(r)
	res, err := h.service.FindAllWithPagination(ctx, filter)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("list action roles error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("cps_action_role.count", len(res.Data)))
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionRoleByCode", "handler", "cpsActionRole")
	defer span.End()
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("cps_action_role.code", code))
	res, err := h.service.GetByActionCode(ctx, code)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("get action role error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createCpsActionRole", "handler", "cpsActionRole")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	var req actionrole_dto.CreateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[BPSActionRoleHandler] invalid input payload: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if req.ActionName == "" {
		localization.SendBadRequestResponse(w, localization.ErrorActionNameIsRequired.Message)
		return
	}
	span.SetAttributes(attribute.String("cps_action_role.code", req.ActionCode))
	action_role := actionrole_dto.CreateActionRoleRequest{
		ActionName:           strings.ToUpper(strings.TrimSpace(req.ActionName)),
		ActionCode:           strings.ToUpper(strings.TrimSpace(req.ActionCode)),
		PortalCardName:       req.PortalCardName,
		IsMakerOnly:          req.IsMakerOnly,
		AssignedViewersRoles: req.AssignedViewersRoles,
		AssignedMakersRoles:  req.AssignedMakersRoles,
		AssignedCheckerRoles: req.AssignedCheckerRoles,
		AssignedAuditorRoles: req.AssignedAuditorRoles,
	}

	// actionrole_dto.CreateActionRoleRequest{ActionCode: req.ActionCode, ActionName: req.ActionName, AssignedCheckerRoles: req.AssignedCheckerRoles, AssignedAuditorRoles: req.AssignedAuditorRoles}
	err := h.service.Create(ctx, action_role)

	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("create action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleCreateRequestCreated, nil)
	}
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateCpsActionRole", "handler", "cpsActionRole")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	var req actionrole_dto.UpdateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}

	span.SetAttributes(attribute.String("cps_action_role.code", code))
	err := h.service.Update(ctx, code, req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("update action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleUpdatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleUpdateRequestCreated, nil)
	}
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableCpsActionRole", "handler", "cpsActionRole")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("cps_action_role.code", code))
	if err := h.service.Enable(ctx, code); err != nil {
		span.RecordError(err)
		h.logger.Errorf("enable action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleEnabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleEnableRequestCreated, nil)
	}
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableCpsActionRole", "handler", "cpsActionRole")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("cps_action_role.code", code))
	if err := h.service.Disable(ctx, code); err != nil {
		span.RecordError(err)
		h.logger.Errorf("disable action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleDisableRequestCreated, nil)
	}
}
