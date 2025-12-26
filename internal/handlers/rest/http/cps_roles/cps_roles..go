package cpsroles

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/cps_roles"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type cpsRolesHandler struct {
	svc    service.CPSRolesService
	logger utils.Logger
}

func NewCPSRolesHandler(svc service.CPSRolesService, logger utils.Logger) *cpsRolesHandler {
	return &cpsRolesHandler{svc: svc, logger: logger}
}

func (c *cpsRolesHandler) CreateCPSRole(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCPSRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Errorf("[CreateCPSRole] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		c.logger.Errorf("[CreateCPSRole] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.Create(r.Context(), req); err != nil {
		c.logger.Errorf("[CreateCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleCreated, nil)
}

func (c *cpsRolesHandler) UpdateCPSRole(w http.ResponseWriter, r *http.Request) {
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[UpdateCPSRole] extractID: %v", err)
		return
	}

	var req dto.UpdateCPSRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Errorf("[UpdateCPSRole] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		c.logger.Errorf("[UpdateCPSRole] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.Update(r.Context(), id, req); err != nil {
		c.logger.Errorf("[UpdateCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleUpdated, nil)
}

func (c *cpsRolesHandler) GetAllCPSRoles(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	roles, err := c.svc.FindAllWithPagination(r.Context(), filterParams)
	if err != nil {
		c.logger.Errorf("[GetAllCPSRoles] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRolesFetched, roles)
}

func (c *cpsRolesHandler) GetCPSRole(w http.ResponseWriter, r *http.Request) {
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[GetCPSRole] extractID: %v", err)
		return
	}

	role, err := c.svc.FindById(r.Context(), id)
	if err != nil {
		c.logger.Errorf("[GetCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleFetched, role)
}

func (c *cpsRolesHandler) EnableCPSRole(w http.ResponseWriter, r *http.Request) {
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[EnableCPSRole] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(r.Context(), id, true); err != nil {
		c.logger.Errorf("[EnableCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleEnabled, nil)
}

func (c *cpsRolesHandler) DisableCPSRole(w http.ResponseWriter, r *http.Request) {
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[DisableCPSRole] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(r.Context(), id, false); err != nil {
		c.logger.Errorf("[DisableCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleDisabled, nil)
}
