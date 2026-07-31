package roles

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	roles_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/roles"
	inbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/roles"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	service "github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	common_utils "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"

	constants "github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
)

type RoleHandler struct {
	service service.RoleService
	logger  utils.Logger
}

func NewRoleHandler(service service.RoleService, logger utils.Logger) inbound.RolesInbound {
	return &RoleHandler{
		service: service,
		logger:  logger,
	}
}

func (j *RoleHandler) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	log := common_utils.LoggerFromCtx(r.Context(), j.logger)

	filterParams, err := common_utils.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

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

	resp, err := j.service.FindAllWithPagination(r.Context(), *filterParams)
	if err != nil {
		log.Errorf("[Roles][GetAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSRolesFetched, resp)
}

func (j *RoleHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	log := common_utils.LoggerFromCtx(r.Context(), j.logger)
	data, err := j.service.FindAll(r.Context())
	if err != nil {
		log.Errorf("[Roles][GetAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[Roles][GetAll] fetched %d roles", len(*data))
	localization.SendSuccessResponse(w, localization.SuccessCPSRoleFetched, data)
}

func (j *RoleHandler) FindById(w http.ResponseWriter, r *http.Request) {
	log := common_utils.LoggerFromCtx(r.Context(), j.logger)
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	role, err := j.service.FindById(r.Context(), id)
	if err != nil {
		log.Errorf("[Roles][GetByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCPSRolesFetched, role)
}

func (j *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	var body roles_dto.CreateJobRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := body.Validate(); err != nil {
		log.Errorf("[JobRoleHandler] error: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	role := imodel.Role{
		Name:        strings.TrimSpace(body.Name),
		Code:        body.Code,
		Type:        strings.ToUpper(strings.TrimSpace(body.Type)),
		Description: strings.TrimSpace(body.Description),
		CreatedAt:   time.Now(),
	}

	if err := j.service.Create(ctx, role); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleCreatedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleCreatedRequestSent, nil)
	}
}

func (j *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var body roles_dto.UpdateJobRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := body.Validate(); err != nil {
		log.Errorf("[Job Role Handler] error: %v", err.Error())
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	updated := imodel.Role{
		UpdatedAt: time.Now(),
	}
	if body.Name != "" {
		updated.Name = strings.TrimSpace(body.Name)
	}
	if body.Code != "" {
		updated.Code = strings.TrimSpace(body.Code)
	}
	if body.Type != "" {
		updated.Type = strings.ToUpper(strings.TrimSpace(body.Type))
	}
	if body.Description != "" {
		updated.Description = strings.TrimSpace(body.Description)
	}

	if err := j.service.Update(ctx, id, updated); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[Job Role Handler] error: %v", err.Error())
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleUpdatedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleUpdatedRequestSent, nil)

	}
}

func (j *RoleHandler) Enable(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.EnableOrDisable(ctx, id, true); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[RoleHandler][Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleEnabledSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleEnabledRequestSent, nil)
	}
}

func (j *RoleHandler) Disable(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.EnableOrDisable(ctx, id, false); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[RoleHandler][Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDisabledSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDisabledRequestSent, nil)
	}
}

func (j *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := common_utils.LoggerFromCtx(ctx, j.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.Delete(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[RoleHandler][Delete] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDeletedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessRoleDeletedRequestSent, nil)
	}
}
