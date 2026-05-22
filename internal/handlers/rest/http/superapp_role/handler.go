package superapprole

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	superapproledto "cbe-super-app-cps-action/internal/constants/dto/superapp_role"
	sar_iface "cbe-super-app-cps-action/internal/constants/interfaces/superapp_role"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SuperAppRoleAdapter struct {
	svc    service.SuperAppRoleService
	logger utils.Logger
}

func NewSuperAppRoleHandler(svc service.SuperAppRoleService, logger utils.Logger) sar_iface.SuperAppRole {
	return &SuperAppRoleAdapter{svc: svc, logger: logger}
}

// GetAllSuperAppRoles lists all superapp roles grouped by role name.
//
//	@Summary		List SuperApp Roles
//	@Description	Retrieves a paginated list of superapp roles with their segments
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200	{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles [get]
func (h *SuperAppRoleAdapter) GetAllSuperAppRoles(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, err := h.svc.FindAllWithPagination(r.Context(), *filterParams)
	if err != nil {
		log.Errorf("[GetAllSuperAppRoles] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendPaginatedSuccessResponse(w, localization.SuperAppRoleFetchedSuccessfully, result.Data, result.Meta)
}

// GetTransferLimitByRole fetches service-level transfer limits for a superapp role from CoreIO.
//
//	@Summary		Get Transfer Limits by SuperApp Role
//	@Description	Fetches service-level transfer limits from the core banking system for the given superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Produce		json
//	@Param			role		path		string	true	"SuperApp Role code"
//	@Param			page		query		int		false	"Page number"
//	@Param			per_page	query		int		false	"Items per page"
//	@Success		200	{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role}/limits [get]
func (h *SuperAppRoleAdapter) GetTransferLimitByRole(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)

	role := chi.URLParam(r, "role")
	if role == "" {
		log.Errorf("[GetTransferLimitByRole] missing role param")
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	filterParams := local_util.ExtractFilterParams(r)

	result, err := h.svc.GetTransferLimitByRole(r.Context(), role, *filterParams)
	if err != nil {
		log.Errorf("[GetTransferLimitByRole] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendPaginatedSuccessResponse(w, localization.SuperAppRoleFetchedSuccessfully, result.Data, result.Meta)
}

// GetGlobalLimitByRole fetches the overall channel-level limits for a superapp role from CoreIO.
//
//	@Summary		Get Global Limits by SuperApp Role
//	@Description	Fetches global (channel-level) limits from the core banking system for the given superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Produce		json
//	@Param			role	path		string	true	"SuperApp Role code"
//	@Success		200	{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role}/global-limits [get]
func (h *SuperAppRoleAdapter) GetGlobalLimitByRole(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)

	role := chi.URLParam(r, "role")
	if role == "" {
		log.Errorf("[GetGlobalLimitByRole] missing role param")
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	result, err := h.svc.GetGlobalLimitByRole(r.Context(), role)
	if err != nil {
		log.Errorf("[GetGlobalLimitByRole] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuperAppRoleGlobalLimitFetchedSuccessfully, result)
}

// EnableByRole enables all segments for a given superapp role.
//
//	@Summary		Enable SuperApp Role
//	@Description	Enables all segments belonging to the given superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Produce		json
//	@Param			role	path		string	true	"SuperApp Role code"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role}/enable [patch]
func (h *SuperAppRoleAdapter) EnableByRole(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "EnableByRole", "handler", "superAppRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	role := chi.URLParam(r, "role")
	if role == "" {
		log.Errorf("[EnableByRole] missing role param")
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	if err := h.svc.EnableByRole(ctx, role); err != nil {
		span.RecordError(err)
		log.Errorf("[EnableByRole] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuperAppRoleEnabledSuccessfully, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuperAppRoleEnableSubmittedSuccessfully, nil)
}

// DisableByRole disables all segments for a given superapp role.
//
//	@Summary		Disable SuperApp Role
//	@Description	Disables all segments belonging to the given superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Produce		json
//	@Param			role	path		string	true	"SuperApp Role code"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role}/disable [patch]
func (h *SuperAppRoleAdapter) DisableByRole(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "DisableByRole", "handler", "superAppRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	role := chi.URLParam(r, "role")
	if role == "" {
		log.Errorf("[DisableByRole] missing role param")
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	if err := h.svc.DisableByRole(ctx, role); err != nil {
		span.RecordError(err)
		log.Errorf("[DisableByRole] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuperAppRoleDisabledSuccessfully, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuperAppRoleDisableSubmittedSuccessfully, nil)
}

// DeleteByRole deletes all segments for a given superapp role.
//
//	@Summary		Delete SuperApp Role
//	@Description	Deletes all segments belonging to the given superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Produce		json
//	@Param			role	path		string	true	"SuperApp Role code"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role} [delete]
func (h *SuperAppRoleAdapter) DeleteByRole(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "DeleteByRole", "handler", "superAppRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	role := chi.URLParam(r, "role")
	if role == "" {
		log.Errorf("[DeleteByRole] missing role param")
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	if err := h.svc.DeleteByRole(ctx, role); err != nil {
		span.RecordError(err)
		log.Errorf("[DeleteByRole] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuperAppRoleDeletedSuccessfully, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuperAppRoleDeleteSubmittedSuccessfully, nil)
}

// GetAccessListsByRole returns enabled and disabled access lists for a superapp role.
//
//	@Summary		Get Access Lists by SuperApp Role
//	@Description	Returns enabled and disabled access lists for the given superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Produce		json
//	@Param			role	path		string	true	"SuperApp Role code"
//	@Success		200	{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role}/access-lists [get]
func (h *SuperAppRoleAdapter) GetAccessListsByRole(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)

	role := chi.URLParam(r, "role")
	if role == "" {
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	enabled, disabled, err := h.svc.GetAccessListsByRole(r.Context(), role)
	if err != nil {
		log.Errorf("[GetAccessListsByRole] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuperAppRoleAccessListsFetchedSuccessfully, map[string]interface{}{
		"enabled_services":  enabled,
		"disabled_services": disabled,
	})
}

// BulkDisableAccessLists disables a list of access lists for a superapp role.
//
//	@Summary		Bulk Disable Access Lists
//	@Description	Adds the given access list IDs to the block list for the superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			role	path		string									true	"SuperApp Role code"
//	@Param			body	body		superapproledto.BulkAccessListByRoleRequest	true	"Access list IDs"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role}/access-lists/disable [patch]
func (h *SuperAppRoleAdapter) BulkDisableAccessLists(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "BulkDisableAccessLists", "handler", "superAppRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	role := chi.URLParam(r, "role")
	if role == "" {
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	var req superapproledto.BulkAccessListByRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidRequestBody.Code).Error())
		return
	}

	if req.AccessListIDs == nil || len(req.AccessListIDs) == 0 {
		log.Errorf("[BulkDisableAccessLists] empty access list IDs in request body")
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidRequestBody.Code).Error())
		return
	}

	log.Infof("[BulkDisableAccessLists] received request to bulk disable access lists for role %s with access list IDs: %v", role, req.AccessListIDs)

	if err := h.svc.BulkDisableAccessLists(ctx, role, req); err != nil {
		span.RecordError(err)
		log.Errorf("[BulkDisableAccessLists] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuperAppRoleAccessListsBulkDisabledSuccessfully, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuperAppRoleAccessListsBulkDisableSubmittedSuccessfully, nil)
}

// BulkEnableAccessLists enables (removes from block list) access lists for a superapp role.
//
//	@Summary		Bulk Enable Access Lists
//	@Description	Removes the given access list IDs from the block list for the superapp role
//	@Tags			SuperAppRole
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			role	path		string									true	"SuperApp Role code"
//	@Param			body	body		superapproledto.BulkAccessListByRoleRequest	true	"Access list IDs"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/superapp-roles/{role}/access-lists/enable [patch]
func (h *SuperAppRoleAdapter) BulkEnableAccessLists(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "BulkEnableAccessLists", "handler", "superAppRole")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	role := chi.URLParam(r, "role")
	if role == "" {
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidID.Code).Error())
		return
	}

	var req superapproledto.BulkAccessListByRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidRequestBody.Code).Error())
		return
	}

	if req.AccessListIDs == nil || len(req.AccessListIDs) == 0 {
		log.Errorf("[BulkDisableAccessLists] empty access list IDs in request body")
		localization.SendErrorByCodeResponse(w, errors.New(localization.ErrorInvalidRequestBody.Code).Error())
		return
	}

	if err := h.svc.BulkEnableAccessLists(ctx, role, req); err != nil {
		span.RecordError(err)
		log.Errorf("[BulkEnableAccessLists] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuperAppRoleAccessListsBulkEnabledSuccessfully, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuperAppRoleAccessListsBulkEnableSubmittedSuccessfully, nil)
}
