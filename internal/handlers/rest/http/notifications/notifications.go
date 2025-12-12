package notifications

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants/dto/notification"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/notifications/core"
	"cbe-super-app-cps-action/internal/service"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
	req, ok := core.ParseAndValidateNotificationRequest(w, r, true)
	if !ok {
		return
	}

	maker, err := common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	req.CreatedBy = maker.UserID
	domainReq := core.ToDomainNotificationRequest(req)
	_, err = h.service.CreateNotification(r.Context(), domainReq)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationCreationRequestSubmitted, nil)
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
	_, err = h.service.UpdateNotification(r.Context(), id, domainReq)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationUpdateRequestSubmitted, nil)
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
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	if err := h.service.DeleteNotification(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationDeleteRequestSubmitted, nil)
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
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	if err := h.service.EnableNotification(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationEnableRequestSubmitted, nil)
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
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	if err := h.service.DisableNotification(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationDisableRequestSubmitted, nil)
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
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	data, err := h.service.FetchNotificationByID(r.Context(), id)
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorNotificationFetchFailed, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessNotificationRetrieved, data)
}

// FetchNotifications godoc
//
//	@Summary		List notifications
//	@Description	Retrieve notifications with pagination, filtering, and search. Searchable fields: title, notification_code, notification_body, notification_type.
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			page				query	int		false	"Page number"		default(1)
//	@Param			per_page			query	int		false	"Items per page"	default(10)
//	@Param			search				query	string	false	"Search term (searches title, notification_code, notification_body, notification_type)"
//	@Param			is_public			query	bool	false	"Filter by public status"
//	@Param			notification_type	query	string	false	"Filter by notification type"
//	@Param			notification_code	query	string	false	"Filter by notification code"
//	@Param			for					query	string	false	"Filter by notification for"
//	@Param			seen				query	bool	false	"Filter by seen status"
//	@Param			enabled				query	bool	false	"Filter by enabled status"
//	@Param			title				query	string	false	"Filter by title"
//	@Success		200	{object}	localization.StandardResponse{data=paginatedNotificationResponse}	"Notifications retrieved successfully"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/notifications [get]
func (h *handler) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	// Assuming a utility to parse query into types.Filter exists; pass empty for now
	filterParams := common_utils.ExtractFilterParams(r)
	data, err := h.service.FetchNotifications(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationsRetrieved, data)
}
