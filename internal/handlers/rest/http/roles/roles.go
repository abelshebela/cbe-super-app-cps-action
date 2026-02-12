package roles

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	roles_dto "cbe-super-app-cps-action/internal/constants/dto/roles"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/roles"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	service "cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"

	constants "cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
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

// FindAll godoc
//
//	@Summary		Get all roles
//	@Description	Retrieve all roles with pagination and optional search
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int									false	"Page number"		default(1)
//	@Param			per_page	query		int									false	"Items per page"	default(10)
//	@Param			search		query		string								false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Roles retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/roles [get]
func (j *RoleHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	filterParams := common_utils.ExtractFilterParams(r)

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
		j.logger.Errorf("[Roles][GetAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessGetAllBanks, resp)
}

// FindById godoc
//
//	@Summary		Get role by ID
//	@Description	Retrieve a single role by its identifier
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string													true	"Role ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.JobRole}	"Role retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}				"Role not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/roles/{id} [get]
func (j *RoleHandler) FindById(w http.ResponseWriter, r *http.Request) {
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
	localization.SendSuccessResponse(w, localization.SuccessGetOneBank, role)
}

// Create godoc
//
//	@Summary		Create role (maker)
//	@Description	Create a new role with name and optional portal cards
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			body	body		roles_dto.CreateJobRoleRequest	true	"Create role request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Role creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/roles [post]
func (j *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {

	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	var body roles_dto.CreateJobRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := body.Validate(); err != nil {
		j.logger.Errorf("[JobRoleHandler] error: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	role := imodel.JobRole{
		Name:        strings.TrimSpace(body.Name),
		Code:        "ROLE_" + local_util.UniqueIdGenerator(),
		PortalCards: body.PortalCards,
		CreatedAt:   time.Now(),
	}

	if err := j.service.Create(ctx, role); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessRoleCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessRoleCreatedRequestSent, nil)
	}
}

// Update godoc
//
//	@Summary		Update role (maker)
//	@Description	Update an existing role by ID. Provide only fields to change. At least one field must be provided.
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Role ID"
//	@Param			body	body		roles_dto.UpdateJobRoleRequest	true	"Update role request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Role update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Role not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/roles/{id} [patch]
func (j *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

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
		j.logger.Errorf("[Job Role Handler] error: %v", err.Error())
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	updated := imodel.JobRole{
		UpdatedAt: time.Now(),
	}
	if body.Name != "" {
		updated.Name = strings.TrimSpace(body.Name)
	}
	if body.Code != "" {
		updated.Code = strings.TrimSpace(body.Code)
	}

	if len(body.PortalCards) > 0 {
		updated.PortalCards = body.PortalCards
	}

	if err := j.service.Update(ctx, id, updated); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessRoleUpdatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessRoleUpdatedRequestSent, nil)

	}
}
