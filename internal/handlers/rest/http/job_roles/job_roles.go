package job_roles

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	roles_dto "cbe-super-app-cps-action/internal/constants/dto/job_role"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/job_role"
	"cbe-super-app-cps-action/internal/constants/localization"
	service "cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"

	sharedmodel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"

	"cbe-super-app-cps-action/internal/constants"
	types "cbe-super-app-cps-action/internal/constants/types"
)

type JobRoleHandler struct {
	service service.JobRoleService
	logger  utils.Logger
}

func NewJobRoleHandler(service service.JobRoleService, logger utils.Logger) inbound.RolesInbound {
	return &JobRoleHandler{
		service: service,
		logger:  logger,
	}
}

// GetAll godoc
//
//	@Summary		Get all job roles
//	@Description	Retrieve all job roles with pagination and optional search
//	@Tags			Job Title with Role
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int									false	"Page number"		default(1)
//	@Param			per_page	query		int									false	"Items per page"	default(10)
//	@Param			search		query		string								false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Job roles retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/job_roles [get]
func (j *JobRoleHandler) GetAllWithPagination(w http.ResponseWriter, r *http.Request) {

	filterParams := common_utils.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := common_utils.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := common_utils.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	resp, err := j.service.FindAllWithPagination(r.Context(), *filterParams)
	if err != nil {
		j.logger.Errorf("[Roles][GetAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessJobRolesFetchedSuccessfully, resp)
}

func (j *JobRoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := j.service.FindAll(r.Context())
	if err != nil {
		j.logger.Errorf("[Roles][GetAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	j.logger.Infof("[Roles][GetAll] data: %v", data)
	localization.SendSuccessResponse(w, localization.SuccessJobRolesFetchedSuccessfully, data)
}

// GetByID godoc
//
//	@Summary		Get job role by ID
//	@Description	Retrieve a job role by its ID
//	@Tags			Job Title with Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Job role ID"
//	@Success		200	{object}	localization.StandardResponse{data=object}	"Job role retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/job_roles/{id} [get]
func (j *JobRoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	role, err := j.service.FindById(r.Context(), id)
	if err != nil {
		j.logger.Errorf("[Roles][GetByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessJobRoleFetchedSuccessfully, role)
}

// Create godoc
//
//	@Summary		Create job role
//	@Description	Create a new job role with job title and role
//	@Tags			Job Title with Role
//	@Accept			json
//	@Produce		json
//	@Param			body	body		roles.RequestRolesCreate	true	"Create job role request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}		"Job role creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/job_roles [post]
func (j *JobRoleHandler) Create(w http.ResponseWriter, r *http.Request) {

	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	var body roles_dto.RequestRolesCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := body.Validate(); err != nil {
		j.logger.Errorf("[JobRoleHandler] error: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	role := sharedmodel.Role{
		JobTitle:  strings.TrimSpace(body.JobTitle),
		Role:      strings.TrimSpace(body.Role),
		CreatedAt: time.Now(),
	}

	if err := j.service.Create(ctx, role); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleCreatedRequestSent, nil)
	}
}

// Update godoc
//
//	@Summary		Update job role
//	@Description	Update a job role by its ID
//	@Tags			Job Title with Role
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Job role ID"
//	@Param			body	body		roles.RequestRolesUpdate	true	"Update job role request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}		"Job role update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/job_roles/{id} [patch]
func (j *JobRoleHandler) Update(w http.ResponseWriter, r *http.Request) {

	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var body roles_dto.RequestRolesUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := body.Validate(); err != nil {
		j.logger.Errorf("[Job Role Handler] error: %v", err.Error())
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	updated := sharedmodel.Role{
		UpdatedAt: time.Now(),
	}
	if body.JobTitle != "" {
		updated.JobTitle = strings.TrimSpace(body.JobTitle)
	}
	if body.Role != "" {
		updated.Role = strings.TrimSpace(body.Role)
	}

	if err := j.service.Update(ctx, id, updated); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleUpdatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleUpdatedRequestSent, nil)
	}
}

// Enable godoc
//
//	@Summary		Enable job role
//	@Description	Enable a job role by its ID
//	@Tags			Job Title with Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Job role ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Job role enabled successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/job_roles/{id}/enable [patch]
func (j *JobRoleHandler) Enable(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.EnableOrDisable(ctx, id, true); err != nil {
		j.logger.Errorf("[JobRoleHandler][Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleEnabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleEnabledRequestSent, nil)
	}
}

// Disable godoc
//
//	@Summary		Disable job role
//	@Description	Disable a job role by its ID
//	@Tags			Job Title with Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Job role ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Job role disabled successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/job_roles/{id}/disable [patch]
func (j *JobRoleHandler) Disable(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.EnableOrDisable(ctx, id, false); err != nil {
		j.logger.Errorf("[JobRoleHandler][Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleDisabledRequestSent, nil)
	}
}

// Delete godoc
//
//	@Summary		Delete job role
//	@Description	Soft delete a job role by its ID
//	@Tags			Job Title with Role
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Job role ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Job role deleted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/job_roles/{id} [delete]
func (j *JobRoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := j.service.Delete(ctx, id); err != nil {
		j.logger.Errorf("[JobRoleHandler][Delete] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleDeletedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessJobRoleDeletedRequestSent, nil)
	}
}
