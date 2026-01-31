package avatar

import (
	avatarInbound "cbe-super-app-cps-action/internal/constants/interfaces/avatar"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"mime/multipart"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/constants"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type paginatedAvatarResp types.PaginatedResponse[[]*model.Avatar]
type avatarAdapter struct {
	avatarApplication service.AvatarService
	logger            utils.Logger
}

func InitAvatarAdapter(avatarApplication service.AvatarService, logger utils.Logger) avatarInbound.AvatarInbound {
	return &avatarAdapter{
		logger:            logger,
		avatarApplication: avatarApplication,
	}
}

// CreateAvatar godoc
//
//	@Summary		Create avatar (maker)
//	@Description	Upload an avatar image with a label.
//	@Tags			Avatars
//	@Accept			mpfd
//	@Produce		json
//	@Param			label	formData	string									true	"Label"	example("Gold")
//	@Param			avatar	formData	file									true	"Avatar image (<=2MB; jpeg/png/gif)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Avatar create request sent"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/avatar [post]
func (a *avatarAdapter) CreateAvatar(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createAvatar", "handler", "avatar")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	req, err := ReqFileParse(r)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[CreateAvatar] failed to parse file: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		a.logger.Errorf("[CreateAvatar] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("avatar.label", req.Label))
	if err := a.avatarApplication.CreateAvatar(ctx, &model.Avatar{Label: req.Label}, req.Avatar); err != nil {
		span.RecordError(err)
		a.logger.Errorf("[CreateAvatar] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessAvatarCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessAvatarCreated, nil)
		a.logger.Infof("[CreateAvatar] request sent successfully for label: %s", req.Label)
	}
}

// DeleteAvatar godoc
//
//	@Summary		Delete avatar (maker)
//	@Description	Submit delete request for an avatar by ID.
//	@Tags			Avatars
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Avatar ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Avatar delete request sent"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Avatar not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/avatar/{id} [delete]
func (a *avatarAdapter) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteAvatar", "handler", "avatar")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	span.SetAttributes(attribute.String("avatar.id", id))

	if err := a.avatarApplication.DeleteAvatar(ctx, id); err != nil {
		span.RecordError(err)
		a.logger.Errorf("[DeleteAvatar] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessAvatarDeletedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessAvatarDeleted, nil)
		a.logger.Infof("[DeleteAvatar] request sent successfully for id: %s", id)
	}
}

// Enable godoc
//
//	@Summary		Enable avatar (checker)
//	@Description	Approve enable request for an avatar.
//	@Tags			Avatars
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Avatar ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Avatar enabled"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Avatar not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/avatar/enable/{id} [patch]
func (a *avatarAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableAvatar", "handler", "avatar")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	span.SetAttributes(attribute.String("avatar.id", id))

	if err := a.avatarApplication.EnableDisable(ctx, id, true); err != nil {
		span.RecordError(err)
		a.logger.Errorf("[Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessAvatarEnabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessAvatarEnabled, nil)
		a.logger.Infof("[Enable] avatar enabled successfully for id: %s", id)
	}
}

// Disable godoc
//
//	@Summary		Disable avatar (checker)
//	@Description	Approve disable request for an avatar.
//	@Tags			Avatars
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Avatar ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Avatar disabled"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Avatar not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/avatar/disable/{id} [patch]
func (a *avatarAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableAvatar", "handler", "avatar")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	span.SetAttributes(attribute.String("avatar.id", id))

	if err := a.avatarApplication.EnableDisable(ctx, id, false); err != nil {
		span.RecordError(err)
		a.logger.Errorf("[Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessAvatarDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessAvatarDisabled, nil)

		a.logger.Infof("[Disable] avatar disabled successfully for id: %s", id)
	}
}

// FetchAvatar godoc
//
//	@Summary		Get avatar by ID
//	@Description	Retrieve a single avatar by ID.
//	@Tags			Avatars
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string												true	"Avatar ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.Avatar}	"Avatar retrieved"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}				"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}				"Avatar not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}				"Server error"
//	@Security		BearerAuth
//	@Router			/avatar/{id} [get]
func (a *avatarAdapter) FetchAvatar(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchAvatar", "handler", "avatar")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	span.SetAttributes(attribute.String("avatar.id", id))

	res, err := a.avatarApplication.FetchAvatarById(ctx, id)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[FetchAvatar] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("[FetchAvatar] avatar retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessAvatarEnabled, res)
}

// FetchAvatars godoc
//
//	@Summary		List avatars
//	@Description	Retrieve avatars with pagination.
//	@Tags			Avatars
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int														false	"Page number"		default(1)	minimum(1)	example(1)
//	@Param			per_page	query		int														false	"Items per page"	default(10)	minimum(1)	maximum(100)	example(10)
//	@Param			search		query		string													false	"Search term"		example("Gold")
//	@Success		200			{object}	localization.StandardResponse{data=paginatedAvatarResp}	"Avatars retrieved"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}					"Server error"
//	@Security		BearerAuth
//	@Router			/avatar [get]
func (a *avatarAdapter) FetchAvatars(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchAvatars", "handler", "avatar")
	defer span.End()
	filterParam := local_util.ExtractFilterParams(r)

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

	avatars, err := a.avatarApplication.FetchAllAvatar(ctx, *filterParam)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[FetchAvatars] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("avatar.count", len(avatars.Data)))
	a.logger.Infof("[FetchAvatars] retrieved %d avatars", len(avatars.Data))
	localization.SendSuccessResponse(w, localization.SuccessAvatarRetrieved, avatars)
}

// UpdateAvatar godoc
//
//	@Summary		Update avatar (maker)
//	@Description	Update avatar label and/or image. Provide only fields to change.
//	@Tags			Avatars
//	@Accept			mpfd
//	@Produce		json
//	@Param			id		path		string									true	"Avatar ID"
//	@Param			label	formData	string									false	"Label"	example("Silver")
//	@Param			avatar	formData	file									false	"Avatar image (<=2MB; jpeg/png/gif)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Avatar update request sent"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid input"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Avatar not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/avatar/{id} [patch]
func (a *avatarAdapter) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateAvatar", "handler", "avatar")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	var inputData *multipart.FileHeader
	var label string
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Message)
		return
	}

	req, err := ReqFileParse(r)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[UpdateAvatar] failed to parse file: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if req.Avatar != nil {
		inputData = req.Avatar
		label = req.Label
	}
	span.SetAttributes(
		attribute.String("avatar.id", id),
		attribute.String("avatar.label", label),
	)
	if err := a.avatarApplication.UpdateAvatar(ctx, id, &model.Avatar{Label: label}, inputData, false); err != nil {
		span.RecordError(err)
		a.logger.Errorf("[UpdateAvatar] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessAvatarUpdatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessAvatarUpdated, nil)

		a.logger.Infof("[UpdateAvatar] request sent successfully for id: %s", id)
	}
}
