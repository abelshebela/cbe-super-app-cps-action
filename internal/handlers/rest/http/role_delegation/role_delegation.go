package role_delegation_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	role_delegation_dto "cbe-super-app-cps-action/internal/constants/dto/role_delegation"
	role_delegation_handler "cbe-super-app-cps-action/internal/constants/interfaces/role_delegation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	role_delegation_core "cbe-super-app-cps-action/internal/handlers/rest/http/role_delegation/core"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type RoleDelegationHandler struct {
	service service.RoleDelegationService
	logger  utils.Logger
}

// Create implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "CreateRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	var body role_delegation_dto.RoleDelegationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Errorf("[RoleDelegationHandler][Create] invalid request body: %v", err)
		localization.SendBadRequestResponse(w, "invalid request body")
		return
	}

	delegation, err := role_delegation_core.BuildRoleDelegationRequest(body)
	if err != nil {
		log.Errorf("[RoleDelegationHandler][Create] validation error: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := j.service.Create(ctx, delegation); err != nil {
		log.Errorf("[RoleDelegationHandler][Create] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDelegationCreated, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessRoleDelegationCreateRequestSubmitted, nil)
}

// Update implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "UpdateRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		log.Errorf("[RoleDelegationHandler][Update] missing id path parameter")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var body role_delegation_dto.RoleDelegationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Errorf("[RoleDelegationHandler][Update] invalid request body: %v", err)
		localization.SendBadRequestResponse(w, "invalid request body")
		return
	}

	delegation, err := role_delegation_core.BuildRoleDelegationRequest(body)
	if err != nil {
		log.Errorf("[RoleDelegationHandler][Update] validation error: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := j.service.Update(ctx, id, delegation); err != nil {
		log.Errorf("[RoleDelegationHandler][Update] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDelegationUpdated, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessRoleDelegationUpdateRequestSubmitted, nil)
}

// Delete implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "DeleteRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		log.Errorf("[RoleDelegationHandler][Delete] missing id path parameter")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.Delete(ctx, id); err != nil {
		log.Errorf("[RoleDelegationHandler][Delete] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDelegationDelete, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessRoleDelegationDeleteRequestSubmitted, nil)
}

// Disable implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "DisableRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		log.Errorf("[RoleDelegationHandler][Disable] missing id path parameter")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.EnableOrDisable(ctx, id, false); err != nil {
		log.Errorf("[RoleDelegationHandler][Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDelegationDisabled, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessRoleDelegationDisableRequestSubmitted, nil)
}

// Enable implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "EnableRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		log.Errorf("[RoleDelegationHandler][Enable] missing id path parameter")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.EnableOrDisable(ctx, id, true); err != nil {
		log.Errorf("[RoleDelegationHandler][Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDelegationEnabled, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessRoleDelegationEnableRequestSubmitted, nil)
}

// FindById implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) FindById(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "FindByIdRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		log.Errorf("[RoleDelegationHandler][FindById] missing id path parameter")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	data, err := j.service.FindById(ctx, id)
	if err != nil {
		log.Errorf("[RoleDelegationHandler][FindById] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, data)
}

// FindAll implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "FindAllRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	data, err := j.service.FindAll(ctx)
	if err != nil {
		log.Errorf("[RoleDelegationHandler][FindAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, data)
}

// FindAllWithPagination implements [role_delegation_outbound.RoleDelegation].
func (j *RoleDelegationHandler) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "FindAllWithPaginationRoleDelegation", "handler", "role_delegation")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		log.Errorf("[RoleDelegationHandler][FindAllWithPagination] invalid search filter: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		log.Errorf("[RoleDelegationHandler][FindAllWithPagination] invalid filter value: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	resp, err := j.service.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		log.Errorf("[RoleDelegationHandler][FindAllWithPagination] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, resp)
}

func NewRoleDelegationHandler(service service.RoleDelegationService, logger utils.Logger) role_delegation_handler.RoleDelegation {
	return &RoleDelegationHandler{
		service: service,
		logger:  logger,
	}
}
