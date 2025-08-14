package notification

import (
	"encoding/json"
	"fmt"
	"net/http"

	notification_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/notification"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	notification_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	notification_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/notification"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	c "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// Constants for pagination
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// NotificationHTTPStore handles HTTP requests for notifications
type NotificationHTTPStore struct {
	Application notification_application.NotificationApplicationAbstracts
	logger      c.Logger
}

// NewNotificationHTTPHandler creates a new notification HTTP handler
func NewNotificationHTTPHandler(app notification_application.NotificationApplicationAbstracts, logger c.Logger) notification_inbound.NotificationHandler {
	return &NotificationHTTPStore{
		Application: app,
		logger:      logger,
	}
}

// extractUserAndMaker extracts user context and creates a maker, sending an error response if incomplete
func (h *NotificationHTTPStore) extractUserAndMaker(w http.ResponseWriter, r *http.Request) (cps_entities.User, bool) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		h.logger.Errorf("Incomplete user context")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return cps_entities.User{}, false
	}
	maker := cps_entities.User{
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}
	return maker, true
}

// extractID extracts and validates the ID parameter, sending an error response if invalid
func (h *NotificationHTTPStore) extractID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("Missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return "", false
	}
	return id, true
}

// parseAndValidateNotificationRequest parses and validates the notification request from JSON
func (h *NotificationHTTPStore) parseAndValidateNotificationRequest(w http.ResponseWriter, r *http.Request, isCreate bool) (NotificationRequest, bool) {
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Failed to parse notification request: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return NotificationRequest{}, false
	}

	if err := req.Validate(isCreate); err != nil {
		h.logger.Errorf("Validation failed: %v", err)
		common_util.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return NotificationRequest{}, false
	}
	return req, true
}

// sendResponse sends success or error responses, preserving original logic
func (h *NotificationHTTPStore) sendResponse(w http.ResponseWriter, err error, successMessage string, data interface{}) {
	if err != nil {
		h.logger.Errorf("Request failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	common_util.WriteSuccessResponse(w, data, successMessage)
}

// toDomainNotificationRequest converts NotificationRequest to the domain request struct
func (h *NotificationHTTPStore) toDomainNotificationRequest(req NotificationRequest) notification_domain.NotificationRequest {
	return *MapNotificationRequestToDomain(&req)
}

// CreateNotification handles the creation of a new notification
func (h *NotificationHTTPStore) CreateNotification(w http.ResponseWriter, r *http.Request) {
	req, ok := h.parseAndValidateNotificationRequest(w, r, true)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	req.CreatedBy = maker.UserCode
	domainReq := h.toDomainNotificationRequest(req)
	err := h.Application.CreateNotification(r.Context(), domainReq, maker)
	h.sendResponse(w, err, "Notification creation request submitted successfully", nil)
}

// UpdateNotification handles updating an existing notification
func (h *NotificationHTTPStore) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	req, ok := h.parseAndValidateNotificationRequest(w, r, false)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	if req.IsEmpty() {
		h.logger.Errorf("No data provided for update")
		common_util.SendErrorResponse(w, common_util.NoDataProvidedForUpdate, http.StatusBadRequest, nil)
		return
	}

	domainReq := h.toDomainNotificationRequest(req)
	err := h.Application.UpdateNotification(r.Context(), id, domainReq, maker)
	h.sendResponse(w, err, "Notification update request submitted successfully", nil)
}

// DeleteNotification handles deleting a notification
func (h *NotificationHTTPStore) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.DeleteNotification(r.Context(), id, maker)
	h.sendResponse(w, err, "Notification delete request submitted successfully", nil)
}

// EnableNotification handles enabling a notification
func (h *NotificationHTTPStore) EnableNotification(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableDisableNotification(r.Context(), id, maker, true)
	h.sendResponse(w, err, "Notification enable request submitted successfully", nil)
}

// DisableNotification handles disabling a notification
func (h *NotificationHTTPStore) DisableNotification(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableDisableNotification(r.Context(), id, maker, false)
	h.sendResponse(w, err, "Notification disable request submitted successfully", nil)
}

// FetchNotificationByID fetches a notification by its ID
func (h *NotificationHTTPStore) FetchNotificationByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	data, err := h.Application.FetchNotificationByID(r.Context(), id)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	res := MapNotificationToResponse(data)
	h.sendResponse(w, nil, "Notification successfully retrieved", res)
}

// FetchNotifications fetches notifications with pagination and filtering
func (h *NotificationHTTPStore) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	filterParam := common_util.ExtractFilterParams(r)

	fmt.Println(filterParam)

	list, err := h.Application.FetchNotifications(r.Context(), filterParam)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	docs := []NotificationResponse{}
	for _, doc := range list.Data {
		res := MapNotificationToResponse(doc)
		docs = append(docs, *res)
	}

	res := common_util.PaginatedResponse[*[]NotificationResponse]{
		Data: &docs,
		Meta: list.Meta,
	}
	h.sendResponse(w, nil, "Notifications successfully retrieved", res)
}
