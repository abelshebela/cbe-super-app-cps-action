package event

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	event_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/event"
	event_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	event_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/event"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	c "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type EventHTTPStore struct {
	Application event_application.ApplicationAbstracts
	logger      c.Logger
}

func NewEventHTTPHandler(app event_application.ApplicationAbstracts, logger c.Logger) event_inbound.EventHandler {
	return &EventHTTPStore{
		Application: app,
		logger:      logger,
	}
}

// extractUserAndMaker extracts user context and creates a maker, sending an error response if incomplete.
func (h *EventHTTPStore) extractUserAndMaker(w http.ResponseWriter, r *http.Request) (cps_entities.User, bool) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
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

// extractID extracts and validates the ID parameter, sending an error response if invalid.
func (h *EventHTTPStore) extractID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return "", false
	}
	return id, true
}
func PrettyPrintJSON(data interface{}) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON))
	return nil
}

// parseAndValidateEventRequest parses and validates the event request from multipart form.
func (h *EventHTTPStore) parseAndValidateEventRequest(w http.ResponseWriter, r *http.Request, isCreate bool) (EventRequest, bool) {
	req, err := ParseEventRequestFromMultipartForm(r, isCreate)
	if err != nil {
		h.logger.Errorf("failed to parse event request: %v", err)
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return EventRequest{}, false
	}

	if err := req.Validate(isCreate); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return EventRequest{}, false
	}
	return req, true
}

// sendResponse sends success or error responses, preserving original logic.
func (h *EventHTTPStore) sendResponse(w http.ResponseWriter, err error, successMessage string, data interface{}) {
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, data, successMessage)
}

// toDomainEventRequest converts EventRequest to the domain request struct.
func (h *EventHTTPStore) toDomainEventRequest(req EventRequest) event_domain.EventRequest {
	return ToDomainEventRequest(struct {
		MerchantID       string
		EventName        string
		EventVenue       string
		EventCity        string
		StartDate        time.Time
		DueDate          time.Time
		CoverImage       *multipart.FileHeader
		EventDescription string
		TotalTicketCount uint
		Tickets          []entities.Ticket
	}{
		MerchantID:       req.MerchantID,
		EventName:        req.EventName,
		EventVenue:       req.EventVenue,
		EventCity:        req.EventCity,
		StartDate:        req.StartDate,
		DueDate:          req.DueDate,
		CoverImage:       req.CoverImage,
		EventDescription: req.EventDescription,
		TotalTicketCount: req.TotalTicketCount,
		Tickets:          req.Tickets,
	})
}

func (h *EventHTTPStore) CreateEvent(w http.ResponseWriter, r *http.Request) {
	req, ok := h.parseAndValidateEventRequest(w, r, true)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	domainReq := h.toDomainEventRequest(req)
	err := h.Application.CreateEvent(r.Context(), domainReq, maker)
	h.sendResponse(w, err, "Event creation request submitted successfully", nil)
}

func (h *EventHTTPStore) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	req, ok := h.parseAndValidateEventRequest(w, r, false)
	if !ok {
		return
	}

	PrettyPrintJSON(req)
	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	if req.IsEmpty() {
		utils.SendErrorResponse(w, common_util.NoDataProvidedForUpdate, http.StatusBadRequest, nil)
		return
	}

	domainReq := h.toDomainEventRequest(req)
	err := h.Application.UpdateEvent(r.Context(), id, domainReq, maker)
	h.sendResponse(w, err, "Event update request submitted successfully", nil)
}

func (h *EventHTTPStore) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.DeleteEvent(r.Context(), id, maker)
	h.sendResponse(w, err, "Event delete request submitted successfully", nil)
}

func (h *EventHTTPStore) EnableEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableDisableEvent(r.Context(), id, maker, true)
	h.sendResponse(w, err, "Event enable request submitted successfully", nil)
}

func (h *EventHTTPStore) DisableEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableDisableEvent(r.Context(), id, maker, false)
	h.sendResponse(w, err, "Event disable request submitted successfully", nil)
}

func (h *EventHTTPStore) FetchEventByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	data, err := h.Application.FetchEventByID(r.Context(), id)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	res := ToEventResponse(data)
	h.sendResponse(w, nil, "Event successfully retrieved", res)
}

func (h *EventHTTPStore) FetchEvents(w http.ResponseWriter, r *http.Request) {
	filterParam := common_util.ExtractFilterParams(r)

	list, err := h.Application.FetchEvent(r.Context(), filterParam)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	var docs []EventResponse
	for _, doc := range list.Data {
		res := ToEventResponse(doc)
		docs = append(docs, res)
	}

	res := common_util.PaginatedResponse[*[]EventResponse]{
		Data: &docs,
		Meta: list.Meta,
	}
	h.sendResponse(w, nil, "Events successfully retrieved", res)
}
