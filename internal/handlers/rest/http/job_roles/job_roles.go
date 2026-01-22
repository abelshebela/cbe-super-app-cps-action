package job_roles

import (
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

func (j *JobRoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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
	localization.SendSuccessResponse(w, localization.SuccessGetAllBanks, resp)
}

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
	localization.SendSuccessResponse(w, localization.SuccessGetOneBank, role)
}

func (j *JobRoleHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	if err := j.service.Create(r.Context(), role); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessJobRoleCreatedRequestSent, nil)
}

func (j *JobRoleHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	if err := j.service.Update(r.Context(), id, updated); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessJobRoleUpdatedRequestSent, nil)
}
