package bps_action_role_handler

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
//	@Summary	Get all action list
//	@Description	Retrieve all action list with pagination and optional search
//	@Tags		BPS Action Role
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int									false	"Page number"		default(1)
//	@Param		per_page	query		int									false	"Items per page"	default(10)
//	@Param		search		query		string								false	"Search term"
//	@Success	200			{object}	localization.StandardResponse{data=object}	"Action list retrieved successfully"
//	@Failure	400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure	500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security	BearerAuth
//	@Router		/bps-action-list [get]
func (h *BPSActionRoleHandler) GetAllActionList(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllBpsActionRoles", "handler", "bpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
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
		log.Errorf("[BpsRoleH][ListActions] svc err: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, h.logger)
	filter := *local_util.ExtractFilterParams(r)
	res, err := h.service.FindAllWithPagination(ctx, filter)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[BpsRoleH][ListRoles] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("cps_action_role.count", len(res.Data)))
	localization.SendSuccessResponse(w, localization.SuccessActionRolesFetched, res)
}

// GetByActionCode godoc
//
//	@Summary	Get BPS action role by code
//	@Description	Retrieve a BPS action role by action code
//	@Tags		BPS Action Role
//	@Accept		json
//	@Produce	json
//	@Param		code	path		string									true	"Action Code"
//	@Success	200		{object}	localization.StandardResponse{data=actionrole_dto.GetActionRoleByActionCodeRes}	"BPS action role retrieved successfully"
//	@Failure	400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure	404		{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure	500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security	BearerAuth
//	@Router		/bps-action-roles/{code} [get]
func (h *BPSActionRoleHandler) GetByActionCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionRoleByCode", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("cps_action_role.code", code))
	res, err := h.service.GetByActionCode(ctx, code)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[BpsRoleH][GetByCode] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleFetched, res)
}

// GetByActionName godoc
//
//	@Summary	Get action role by name
//	@Tags		ActionRole
//	@Produce	json
//	@Param		name	path		string	true	"Action Name"
//	@Success	200		{object}	localization.StandardResponse
//	@Security	BearerAuth
//	@Router		/action-roles/name/{name} [get]
func (h *BPSActionRoleHandler) GetByActionNameCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getCpsActionRoleByCode", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	name := chi.URLParam(r, "name")
	if name == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	span.SetAttributes(attribute.String("cps_action_role.name", name))
	res, err := h.service.GetByActionNameCode(ctx, name)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[BpsRoleH][GetByCode] svc err: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	var req actionrole_dto.CreateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[BPSActionRoleHandler] invalid input payload: %v", err)
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
		log.Errorf("[BpsRoleH][Create] svc err: %v", err)
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
//	@Summary	Update BPS action role
//	@Description	Update a BPS action role by action code
//	@Tags		BPS Action Role
//	@Accept		json
//	@Produce	json
//	@Param		code	path		string									true	"Action Code"
//	@Param		body	body		actionrole_dto.UpdateActionRoleRequest	true	"Update BPS action role request"
//	@Success	200		{object}	localization.StandardResponse{data=nil}		"BPS action role update request submitted successfully"
//	@Failure	400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure	404		{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure	500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security	BearerAuth
//	@Router		/bps-action-roles/{code} [patch]
func (h *BPSActionRoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateCpsActionRole", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

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
		log.Errorf("[BpsRoleH][Update] svc err: %v", err)
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
//	@Summary	Enable BPS action role
//	@Description	Enable a BPS action role by action code
//	@Tags		BPS Action Role
//	@Accept		json
//	@Produce	json
//	@Param		code	path		string									true	"Action Code"
//	@Success	200		{object}	localization.StandardResponse{data=nil}		"BPS action role enable request submitted successfully"
//	@Failure	400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure	404		{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure	500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security	BearerAuth
//	@Router		/bps-action-roles/{code}/enable [patch]
func (h *BPSActionRoleHandler) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableCpsActionRole", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

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
		log.Errorf("[BpsRoleH][Enable] svc err: %v", err)
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
//	@Summary	Disable BPS action role
//	@Description	Disable a BPS action role by action code
//	@Tags		BPS Action Role
//	@Accept		json
//	@Produce	json
//	@Param		code	path		string									true	"Action Code"
//	@Success	200		{object}	localization.StandardResponse{data=nil}		"BPS action role disable request submitted successfully"
//	@Failure	400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure	404		{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure	500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security	BearerAuth
//	@Router		/bps-action-roles/{code}/disable [patch]
func (h *BPSActionRoleHandler) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableCpsActionRole", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

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
		log.Errorf("[BpsRoleH][Disable] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleDisableRequestCreated, nil)
	}
}
