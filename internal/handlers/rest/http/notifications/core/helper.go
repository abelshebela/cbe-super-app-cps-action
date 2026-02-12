package core

import (
	"encoding/json"
	"net/http"

	notify "cbe-super-app-cps-action/internal/constants/dto/notification"
	localization "cbe-super-app-cps-action/internal/constants/localization"
)

// ParseAndValidateNotificationRequest parses and validates the notification request from JSON
func ParseAndValidateNotificationRequest(w http.ResponseWriter, r *http.Request, isCreate bool) (notify.NotificationRequest, bool) {
	var req notify.NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return notify.NotificationRequest{}, false
	}

	if err := req.Validate(isCreate); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return notify.NotificationRequest{}, false
	}
	return req, true
}

// ToDomainNotificationRequest converts NotificationRequest to the domain request struct
func ToDomainNotificationRequest(req notify.NotificationRequest) notify.NotificationRequest {
	return *notify.MapNotificationRequestToDomain(&req)
}
