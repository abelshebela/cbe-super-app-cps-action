package cpsroles

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/cps_roles"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	constants "cbe-super-app-cps-action/internal/constants"
	types "cbe-super-app-cps-action/internal/constants/types"
)

type cpsRolesHandler struct {
	svc    service.CPSRolesService
	logger utils.Logger
}

func NewCPSRolesHandler(svc service.CPSRolesService, logger utils.Logger) *cpsRolesHandler {
	return &cpsRolesHandler{svc: svc, logger: logger}
}

func (c *cpsRolesHandler) CreateCPSRole(w http.ResponseWriter, r *http.Request) {

	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

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

	if err := c.svc.Create(ctx, req); err != nil {
		c.logger.Errorf("[CreateCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleCreated, nil)
	}
}

func (c *cpsRolesHandler) UpdateCPSRole(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

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

	if err := c.svc.Update(ctx, id, req); err != nil {
		c.logger.Errorf("[UpdateCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleUpdatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleUpdated, nil)
	}
}

func (c *cpsRolesHandler) GetAllCPSRoles(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

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
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[EnableCPSRole] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, true); err != nil {
		c.logger.Errorf("[EnableCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleEnabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleEnabled, nil)
	}
}

func (c *cpsRolesHandler) DisableCPSRole(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		c.logger.Errorf("[DisableCPSRole] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, false); err != nil {
		c.logger.Errorf("[DisableCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleDisabled, nil)
	}
}
