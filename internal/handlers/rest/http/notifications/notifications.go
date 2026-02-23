package notifications

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/notification"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/notifications/core"
	"cbe-super-app-cps-action/internal/service"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type notificationRequest notification.NotificationRequest
type paginatedNotificationResponse types.PaginatedResponse[[]*notification.NotificationResponse]
type handler struct {
	service service.NotificationService
	logger  utils.Logger
}

func InitNotificationHandler(svc service.NotificationService, logger utils.Logger) *handler {
	return &handler{service: svc, logger: logger}
}

// CreateNotification godoc
//
//	@Summary		Create a new notification
//	@Description	Create a new notification with the provided information
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			request	body		notificationRequest						true	"Notification payload"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Notification creation request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications [post]
func (h *handler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "createNotification", "handler", "notification")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	req, ok := core.ParseAndValidateNotificationRequest(w, r, true)
	if !ok {
		return
	}

	maker, err := common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		span.RecordError(err)
		return
	}

	req.CreatedBy = maker.UserID
	domainReq := core.ToDomainNotificationRequest(req)
	span.SetAttributes(
		attribute.String("notification.created_by", maker.UserID),
	)

	_, err = h.service.CreateNotification(ctx, domainReq)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CreateNotification] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessNotificationCreatedSP, nil)
	} else {
		log.Infof("[CreateNotification] request sent successfully by user: %s", maker.UserID)
		localization.SendSuccessResponse(w, localization.SuccessNotificationCreationRequestSubmitted, nil)
	}
}

// UpdateNotification godoc
//
//	@Summary		Update a notification
//	@Description	Update an existing notification by ID
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Notification ID"
//	@Param			request	body		notificationRequest						true	"Notification payload"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Notification update request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Notification not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications/{id} [patch]
func (h *handler) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "updateNotification", "handler", "notification")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	req, ok := core.ParseAndValidateNotificationRequest(w, r, false)
	if !ok {
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	domainReq := core.ToDomainNotificationRequest(req)
	span.SetAttributes(attribute.String("notification.id", id))
	_, err = h.service.UpdateNotification(ctx, id, domainReq)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateNotification] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		log.Infof("[UpdateNotification] Updated successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessNotificationUpdatedSP, nil)
	} else {
		log.Infof("[UpdateNotification] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessNotificationUpdateRequestSubmitted, nil)
	}
}

// DeleteNotification godoc
//
//	@Summary		Delete a notification
//	@Description	Permanently delete a notification by ID
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Notification ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Notification delete request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Notification not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications/{id} [delete]
func (h *handler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "deleteNotification", "handler", "notification")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	span.SetAttributes(attribute.String("notification.id", id))
	if err := h.service.DeleteNotification(ctx, id); err != nil {
		span.RecordError(err)
		log.Errorf("[DeleteNotification] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessNotificationDeletedSP, nil)
	} else {
		log.Infof("[DeleteNotification] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessNotificationDeleteRequestSubmitted, nil)
	}
}

// EnableNotification godoc
//
//	@Summary		Enable a notification
//	@Description	Enable a notification by ID
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Notification ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Notification enable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Notification not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications/enable/{id} [patch]
func (h *handler) EnableNotification(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "enableNotification", "handler", "notification")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		span.RecordError(err)
		return
	}

	span.SetAttributes(attribute.String("notification.id", id))
	if err := h.service.EnableNotification(ctx, id); err != nil {
		span.RecordError(err)
		log.Errorf("[EnableNotification] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessNotificationEnabledSP, nil)
	} else {
		log.Infof("[EnableNotification] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessNotificationEnableRequestSubmitted, nil)
	}

}

// DisableNotification godoc
//
//	@Summary		Disable a notification
//	@Description	Disable a notification by ID
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Notification ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Notification disable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Notification not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications/disable/{id} [patch]
func (h *handler) DisableNotification(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "disableNotification", "handler", "notification")
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		span.RecordError(err)
		return
	}

	span.SetAttributes(attribute.String("notification.id", id))
	if err := h.service.DisableNotification(ctx, id); err != nil {
		span.RecordError(err)
		log.Errorf("[DisableNotification] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessNotificationDisabledSP, nil)
	} else {
		log.Infof("[DisableNotification] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessNotificationDisableRequestSubmitted, nil)
	}

}

// FetchNotificationByID godoc
//
//	@Summary		Get notification by ID
//	@Description	Retrieve a notification's details by ID
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string																	true	"Notification ID"
//	@Success		200	{object}	localization.StandardResponse{data=notification.NotificationResponse}	"Notification retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}									"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications/{id} [get]
func (h *handler) FetchNotificationByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "fetchNotificationById", "handler", "notification")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		return
	}

	span.SetAttributes(attribute.String("notification.id", id))
	data, err := h.service.FetchNotificationByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[FetchNotificationByID] service error: %v", err)
		localization.SendErrorResponse(w, localization.ErrorNotificationFetchFailed, nil, nil)
		return
	}

	log.Infof("[FetchNotificationByID] notification retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessNotificationRetrieved, data)
}

// FetchNotifications godoc
//
//	@Summary		List notifications
//	@Description	Retrieve notifications with pagination, filtering, and search. Searchable fields: title, notification_code, notification_body, notification_type.
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			page				query		int																	false	"Page number"		default(1)
//	@Param			per_page			query		int																	false	"Items per page"	default(10)
//	@Param			search				query		string																false	"Search term (searches title, notification_code, notification_body, notification_type)"
//	@Param			is_public			query		bool																false	"Filter by public status"
//	@Param			notification_type	query		string																false	"Filter by notification type"
//	@Param			notification_code	query		string																false	"Filter by notification code"
//	@Param			for					query		string																false	"Filter by notification for"
//	@Param			seen				query		bool																false	"Filter by seen status"
//	@Param			enabled				query		bool																false	"Filter by enabled status"
//	@Param			title				query		string																false	"Filter by title"
//	@Success		200					{object}	localization.StandardResponse{data=paginatedNotificationResponse}	"Notifications retrieved successfully"
//	@Failure		500					{object}	localization.StandardResponse{data=nil}								"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications [get]
func (h *handler) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "fetchNotifications", "handler", "notification")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	// Assuming a utility to parse query into types.Filter exists; pass empty for now
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

	data, err := h.service.FetchNotifications(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[FetchNotifications] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("notification.count", len(data.Data)))
	log.Infof("[FetchNotifications] retrieved %d notifications", len(data.Data))
	localization.SendSuccessResponse(w, localization.SuccessNotificationsRetrieved, data)
}
