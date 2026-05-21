package superapprole

import (
	"context"
	"errors"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
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
