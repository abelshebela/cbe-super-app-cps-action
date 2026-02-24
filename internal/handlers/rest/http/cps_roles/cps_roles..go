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
	log := local_util.LoggerFromCtx(ctx, c.logger)

	var req dto.CreateCPSRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[CreateCPSRole] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[CreateCPSRole] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.Create(ctx, req); err != nil {
		log.Errorf("[CreateCPSRole] service error: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, c.logger)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[UpdateCPSRole] extractID: %v", err)
		return
	}

	var req dto.UpdateCPSRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[UpdateCPSRole] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[UpdateCPSRole] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.Update(ctx, id, req); err != nil {
		log.Errorf("[UpdateCPSRole] service error: %v", err)
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
	log := local_util.LoggerFromCtx(r.Context(), c.logger)
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
		log.Errorf("[GetAllCPSRoles] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRolesFetched, roles)
}

func (c *cpsRolesHandler) GetCPSRole(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), c.logger)
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[GetCPSRole] extractID: %v", err)
		return
	}

	role, err := c.svc.FindById(r.Context(), id)
	if err != nil {
		log.Errorf("[GetCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleFetched, role)
}

func (c *cpsRolesHandler) EnableCPSRole(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	log := local_util.LoggerFromCtx(ctx, c.logger)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[EnableCPSRole] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, true); err != nil {
		log.Errorf("[EnableCPSRole] service error: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, c.logger)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[DisableCPSRole] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, false); err != nil {
		log.Errorf("[DisableCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleDisabled, nil)
	}
}

func (c *cpsRolesHandler) EnableServiceAccess(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), c.logger)
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[EnableServiceAccess] extractID: %v", err)
		return
	}

	var req dto.ToggleServiceAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[EnableServiceAccess] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[EnableServiceAccess] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.EnableServiceAccess(r.Context(), id, req); err != nil {
		log.Errorf("[EnableServiceAccess] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleServiceEnabled, nil)
}

func (c *cpsRolesHandler) DisableServiceAccess(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), c.logger)
	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[DisableServiceAccess] extractID: %v", err)
		return
	}

	var req dto.ToggleServiceAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[DisableServiceAccess] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[DisableServiceAccess] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := c.svc.DisableServiceAccess(r.Context(), id, req); err != nil {
		log.Errorf("[DisableServiceAccess] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSRoleServiceDisabled, nil)
}

func (c *cpsRolesHandler) DeleteCPSRole(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	log := local_util.LoggerFromCtx(ctx, c.logger)

	id, err := local_util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[DeleteCPSRole] extractID: %v", err)
		return
	}

	if err := c.svc.Delete(ctx, id); err != nil {
		log.Errorf("[DeleteCPSRole] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleDeletedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessCPSRoleDeleted, nil)
	}
}
