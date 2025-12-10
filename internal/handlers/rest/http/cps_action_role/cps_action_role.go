package cps_actionrole_handler

import (
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	cps_actionrole_inbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action_role"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CPSActionRoleHandler struct {
	service service.CPSActionRoleService
	logger  utils.Logger
}

func NewCPSActionRoleHandler(svc service.CPSActionRoleService, logger utils.Logger) cps_actionrole_inbound.CPSActionRoleHandler {
	return &CPSActionRoleHandler{service: svc, logger: logger}
}

// GetAll godoc
// @Summary      List action roles
// @Tags         ActionRole
// @Accept       json
// @Produce      json
// @Param        page     query  int  false  "Page"
// @Param        per_page query  int  false  "Per Page"
// @Param        search   query  string false "Search"
// @Success      200 {object} localization.StandardResponse
// @Security     BearerAuth
// @Router       /action-roles [get]
func (h *CPSActionRoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	filter := *local_util.ExtractFilterParams(r)
	res, err := h.service.FindAllWithPagination(r.Context(), filter)
	if err != nil {
		h.logger.Errorf("list action roles error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRolesFetched, res)
}

// GetByActionCode godoc
// @Summary      Get action role by code
// @Tags         ActionRole
// @Produce      json
// @Param        code  path string true "Action Code"
// @Success      200 {object} localization.StandardResponse
// @Security     BearerAuth
// @Router       /action-roles/{code} [get]
func (h *CPSActionRoleHandler) GetByActionCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	res, err := h.service.GetByActionCode(r.Context(), code)
	if err != nil {
		h.logger.Errorf("get action role error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleFetched, res)
}

// Create godoc
// @Summary      Create action role (maker)
// @Tags         ActionRole
// @Accept       json
// @Produce      json
// @Param        body body actionrole_dto.CreateActionRoleRequest true "Create"
// @Success      201 {object} localization.StandardResponse
// @Security     BearerAuth
// @Router       /action-roles [post]
func (h *CPSActionRoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req actionrole_dto.CreateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}
	if req.ActionName == "" {
		localization.SendBadRequestResponse(w, localization.ErrorActionNameIsRequired.Message)
		return
	}
	err := h.service.Create(r.Context(), actionrole_dto.CreateActionRoleRequest{ActionCode: req.ActionCode, ActionName: req.ActionName, AssignedMakersRoles: req.AssignedMakersRoles, AssignedCheckerRoles: req.AssignedCheckerRoles})
	if err != nil {
		h.logger.Errorf("create action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleCreateRequestCreated, nil)
}

// Update godoc
// @Summary      Update action role (maker)
// @Tags         ActionRole
// @Accept       json
// @Produce      json
// @Param        code path string true "Action Code"
// @Param        body body actionrole_dto.UpdateActionRoleRequest true "Update"
// @Success      201 {object} localization.StandardResponse
// @Security     BearerAuth
// @Router       /action-roles/{code} [patch]
func (h *CPSActionRoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	var req actionrole_dto.UpdateActionRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}
	err := h.service.Update(r.Context(), code, req)
	if err != nil {
		h.logger.Errorf("update action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleUpdateRequestCreated, nil)
}

// Enable godoc
// @Summary      Enable action role (maker)
// @Tags         ActionRole
// @Produce      json
// @Param        code path string true "Action Code"
// @Success      201 {object} localization.StandardResponse
// @Security     BearerAuth
// @Router       /action-roles/{code}/enable [patch]
func (h *CPSActionRoleHandler) Enable(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	if err := h.service.Enable(r.Context(), code); err != nil {
		h.logger.Errorf("enable action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleEnableRequestCreated, nil)
}

// Disable godoc
// @Summary      Disable action role (maker)
// @Tags         ActionRole
// @Produce      json
// @Param        code path string true "Action Code"
// @Success      201 {object} localization.StandardResponse
// @Security     BearerAuth
// @Router       /action-roles/{code}/disable [patch]
func (h *CPSActionRoleHandler) Disable(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameter.Message)
		return
	}
	if err := h.service.Disable(r.Context(), code); err != nil {
		h.logger.Errorf("disable action role failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessActionRoleDisableRequestCreated, nil)
}
