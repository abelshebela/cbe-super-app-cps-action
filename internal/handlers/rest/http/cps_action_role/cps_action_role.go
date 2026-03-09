package cps_actionrole_handler

import (
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/cps_action_role"
	cps_actionrole_inbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action_role"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"

	"cbe-super-app-cps-action/internal/constants"
	types "cbe-super-app-cps-action/internal/constants/types"
)

type CPSActionRoleHandler struct {
	service service.CPSActionRoleService
	logger  utils.Logger
}

func NewCPSActionRoleHandler(svc service.CPSActionRoleService, logger utils.Logger) cps_actionrole_inbound.CPSActionRoleHandler {
	return &CPSActionRoleHandler{service: svc, logger: logger}
}

// GetAllActionList godoc
//
//	@Summary		Get all CPS action list
//	@Description	Retrieve all CPS actions with pagination and optional search
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int									false	"Page number"		default(1)
//	@Param			per_page	query		int									false	"Items per page"	default(10)
//	@Param			search		query		string								false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Action list retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-list [get]
func (h *CPSActionRoleHandler) GetAllActionList(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllCpsActionRoles", "handler", "cpsActionRole")
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
		log.Errorf("[CpsRoleH][ListActions] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("cps_action_role.count", len(res.Data)))
	localization.SendSuccessResponse(w, localization.SuccessActionRolesFetched, res)
}

// GetAll godoc
//
//	@Summary		Get all CPS action roles
//	@Description	Retrieve all CPS action roles with pagination and optional search
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int									false	"Page number"		default(1)
//	@Param			per_page	query		int									false	"Items per page"	default(10)
//	@Param			search		query		string								false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Action roles retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles [get]
func (h *CPSActionRoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllCpsActionRoles", "handler", "cpsActionRole")
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

	res, err := h.service.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CpsRoleH][ListRoles] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("cps_action_role.count", len(res.Data)))
	localization.SendSuccessResponse(w, localization.SuccessActionRolesFetched, res)
}

// GetByActionCode godoc
//
//	@Summary		Get CPS action role by code
//	@Description	Retrieve a single CPS action role by its action code
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string													true	"Action Code"
//	@Success		200		{object}	localization.StandardResponse{data=actionrole_dto.GetActionRoleByActionCodeRes}	"Action role retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}					"Action role not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles/{code} [get]
func (h *CPSActionRoleHandler) GetByActionCode(w http.ResponseWriter, r *http.Request) {
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
		log.Errorf("[CpsRoleH][GetByCode] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleFetched, res)
}

// Create godoc
//
//	@Summary		Create CPS action role (maker)
//	@Description	Create a new CPS action role with assigned roles for viewers, makers, checkers, and auditors
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			body	body		actionrole_dto.CreateActionRoleRequest	true	"Create action role request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Action role creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles [post]
func (h *CPSActionRoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createCpsActionRole", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req actionrole_dto.CreateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[CpsRoleH][Create] decode body err: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
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
		log.Errorf("[CpsRoleH][Create] svc err: %v", err)
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
//	@Summary		Update CPS action role (maker)
//	@Description	Update an existing CPS action role by action code. Provide only fields to change.
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string									true	"Action Code"
//	@Param			body	body		actionrole_dto.UpdateActionRoleRequest	true	"Update action role request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Action role update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Action role not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles/{code} [patch]
func (h *CPSActionRoleHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("cps_action_role.code", code))
	err := h.service.Update(ctx, code, req)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CpsRoleH][Update] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleUpdateRequestCreated, nil)
	}
}

// Enable godoc
//
//	@Summary		Enable CPS action role (checker)
//	@Description	Approve enable request for a CPS action role by action code
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string									true	"Action Code"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Action role enable request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Action role not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles/{code}/enable [patch]
func (h *CPSActionRoleHandler) Enable(w http.ResponseWriter, r *http.Request) {
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
		log.Errorf("[CpsRoleH][Enable] svc err: %v", err)
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
//	@Summary		Disable CPS action role (checker)
//	@Description	Approve disable request for a CPS action role by action code
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string									true	"Action Code"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Action role disable request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Action role not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles/{code}/disable [patch]
func (h *CPSActionRoleHandler) Disable(w http.ResponseWriter, r *http.Request) {
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
		log.Errorf("[CpsRoleH][Disable] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessActionRoleDisableRequestCreated, nil)
	}
}

// GetVersions godoc
//
//	@Summary		Get versions for a CPS action role
//	@Description	Return the list of distinct versions stored in cps_action_approver_index for the given action code
//	@Tags			CPS Action Role
//	@Produce		json
//	@Param			code	path		string	true	"Action Code"
//	@Success		200		{object}	localization.StandardResponse{data=[]int64}	"Versions retrieved"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles/{code}/versions [get]
func (h *CPSActionRoleHandler) GetVersions(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getVersions", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	versions, err := h.service.GetVersionsByActionCode(ctx, code)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CpsRoleH][GetVersions] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleFetched, versions)
}

// GetConfiguredRoles godoc
//
//	@Summary		Get configured roles for a CPS action role
//	@Description	Return the configured viewer, maker, checker, and auditor roles for a given action code
//	@Tags			CPS Action Role
//	@Produce		json
//	@Param			code	path		string	true	"Action Code"
//	@Success		200		{object}	localization.StandardResponse{data=cps_actionrole_dto.ConfiguredRolesResponse}	"Configured roles retrieved"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles/{code}/roles [get]
func (h *CPSActionRoleHandler) GetConfiguredRoles(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getConfiguredRoles", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	res, err := h.service.GetConfiguredRoles(ctx, code)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CpsRoleH][GetConfiguredRoles] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleFetched, res)
}

// UpdateVersionRoleCode godoc
//
//	@Summary		Update role_code in approver index for a specific version
//	@Description	Swap a role_code in cps_action_approver_index for a given action code and version
//	@Tags			CPS Action Role
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string										true	"Action Code"
//	@Param			version	path		int											true	"Version"
//	@Param			body	body		actionrole_dto.UpdateIndexRoleCodeRequest	true	"Role code swap request"
//	@Success		200		{object}	localization.StandardResponse{data=object}	"Role code updated"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps-action-roles/{code}/versions/{version}/role [patch]
func (h *CPSActionRoleHandler) UpdateVersionRoleCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateVersionRoleCode", "handler", "cpsActionRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	code := chi.URLParam(r, "code")
	versionStr := chi.URLParam(r, "version")
	if code == "" || versionStr == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}

	var req actionrole_dto.UpdateIndexRoleCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[CpsRoleH][UpdateVersionRoleCode] decode err: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	count, err := h.service.UpdateIndexRoleCode(ctx, code, version, req.OldRoleCode, req.NewRoleCode)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CpsRoleH][UpdateVersionRoleCode] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessActionRoleFetched, map[string]interface{}{
		"updated_count": count,
	})
}
