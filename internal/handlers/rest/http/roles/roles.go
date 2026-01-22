package roles

import (
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

func (j *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	if err := j.service.Create(r.Context(), role); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessRoleCreatedRequestSent, nil)
}

func (j *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	if err := j.service.Update(r.Context(), id, updated); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessRoleUpdatedRequestSent, nil)
}
