package notifications

import (
	"net/http"

	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/notifications/core"
	"cbe-super-app-cps-action/internal/service"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type handler struct {
	service service.NotificationService
	logger  utils.Logger
}

func InitNotificationHandler(svc service.NotificationService, logger utils.Logger) *handler {
	return &handler{service: svc, logger: logger}
}

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
		localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationEnableRequestSubmitted, nil)
}

// UpdateNotification handles updating an existing notification
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
		localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationUpdateRequestSubmitted, nil)
}

// DeleteNotification handles deleting a notification
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
		localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationDeleteRequestSubmitted, nil)
}

// EnableNotification handles enabling a notification
func (h *handler) EnableNotification(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	if err := h.service.EnableDisableNotification(r.Context(), id); err != nil {
		localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationEnableRequestSubmitted, nil)
}

// DisableNotification handles disabling a notification
func (h *handler) DisableNotification(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		return
	}

	_, err = common_utils.ExtractUserInfo(r.Context(), h.logger)
	if err != nil {
		return
	}

	if err := h.service.EnableDisableNotification(r.Context(), id); err != nil {
		localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationDisableRequestSubmitted, nil)
}

// FetchNotificationByID fetches a notification by its ID
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

// FetchNotifications fetches notifications with pagination and filtering
func (h *handler) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	// Assuming a utility to parse query into types.Filter exists; pass empty for now
	data, err := h.service.FetchNotifications(r.Context(), &types.Filter{})
	if err != nil {
		localization.SendErrorResponse(w, localization.ErrorNotificationFetchFailed, nil, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessNotificationsRetrieved, data)
}
