package eventhandler

import (
	eventInbound "cbe-super-app-cps-action/internal/constants/interfaces/event"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	eventcore "cbe-super-app-cps-action/internal/handlers/rest/http/event/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type eventAdapter struct {
	eventApp service.EventService
	logger   utils.Logger
}

func InitEventAdapter(eventApp service.EventService, logger utils.Logger) eventInbound.EventAdapter {
	return &eventAdapter{
		eventApp: eventApp,
		logger:   logger,
	}
}

// CreateEvent godoc
//
//	@Summary		Create a new event
//	@Description	Create a new event with the provided information (multipart form with ticket arrays)
//	@Tags			Event
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			account_number				formData	string									false	"Account number"
//	@Param			merchant_id					formData	string									true	"Merchant ID"
//	@Param			merchant_name				formData	string									false	"Merchant name"
//	@Param			merchant_phone_number		formData	string									false	"Merchant phone number"
//	@Param			merchant_email				formData	string									false	"Merchant email"
//	@Param			event_name					formData	string									true	"Event name"
//	@Param			event_venue					formData	string									true	"Event venue"
//	@Param			event_city					formData	string									true	"Event city"
//	@Param			start_date					formData	string									true	"Event start date (RFC3339 format)"
//	@Param			due_date					formData	string									true	"Event due date (RFC3339 format)"
//	@Param			cover_image					formData	file									true	"Event cover image"
//	@Param			event_description			formData	string									false	"Event description"
//	@Param			total_ticket_count			formData	integer									true	"Total ticket count"
//	@Param			ticket_name[0]				formData	string									true	"Ticket name (e.g. VIP Ticket)"
//	@Param			ticket_category[0]			formData	string									true	"Ticket category (e.g. VIP)"
//	@Param			ticket_type[0]				formData	string									true	"Ticket type (e.g. Seated)"
//	@Param			ticket_price[0]				formData	number									true	"Ticket price (e.g. 1000)"
//	@Param			ticket_number_of_ticker[0]	formData	integer									true	"Ticket quantity (e.g. 50)"
//	@Param			ticket_name[1]				formData	string									false	"Another Ticket name (e.g. Regular Ticket)"
//	@Param			ticket_category[1]			formData	string									false	"Another Ticket category (e.g. Regular)"
//	@Param			ticket_type[1]				formData	string									false	"Another Ticket type (e.g. Standing)"
//	@Param			ticket_price[1]				formData	number									false	"Another Ticket price (e.g. 500)"
//	@Param			ticket_number_of_ticker[1]	formData	integer									false	"Another Ticket quantity (e.g. 250)"
//	@Success		201							{object}	localization.StandardResponse{data=nil}	"Event created successfully"
//	@Failure		400							{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		401							{object}	localization.StandardResponse{data=nil}	"Unauthorized"
//	@Failure		500							{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/events [post]
func (a *eventAdapter) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createEvent", "handler", "event")
	defer span.End()
	req, err := eventcore.ParseEventRequestFromMultipartForm(r, true)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to parse event request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(true); err != nil {
		span.RecordError(err)
		a.logger.Errorf("event request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("event.merchant_id", req.MerchantID),
		attribute.String("event.event_name", req.EventName),
	)

	if err := a.eventApp.CreateEvent(ctx, req); err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to create event: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("event creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessEventCreationRequestSubmitted, nil)
}

// UpdateEvent godoc
//
//	@Summary		Update an existing event
//	@Description	Updates event details (multipart form with optional tickets & cover image)
//	@Tags			Event
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id							path		string									true	"Event ID"
//	@Param			merchant_id					formData	string									true	"Merchant ID"
//	@Param			event_name					formData	string									true	"Event name"
//	@Param			account_number				formData	string									false	"Account number"
//	@Param			merchant_name				formData	string									false	"Merchant name"
//	@Param			merchant_phone_number		formData	string									false	"Merchant phone number"
//	@Param			merchant_email				formData	string									false	"Merchant email"
//	@Param			event_venue					formData	string									false	"Event venue"
//	@Param			event_city					formData	string									false	"Event city"
//	@Param			start_date					formData	string									false	"Event start date (RFC3339 format)"
//	@Param			due_date					formData	string									false	"Event due date (RFC3339 format)"
//	@Param			cover_image					formData	file									false	"Event cover image"
//	@Param			event_description			formData	string									false	"Event description"
//	@Param			total_ticket_count			formData	integer									false	"Total ticket count"
//	@Param			ticket_name[0]				formData	string									false	"Ticket name (e.g. VIP Ticket)"
//	@Param			ticket_category[0]			formData	string									false	"Ticket category (e.g. VIP)"
//	@Param			ticket_type[0]				formData	string									false	"Ticket type (e.g. Seated)"
//	@Param			ticket_price[0]				formData	number									false	"Ticket price (e.g. 1000)"
//	@Param			ticket_number_of_ticker[0]	formData	integer									false	"Ticket quantity (e.g. 50)"
//	@Param			ticket_name[1]				formData	string									false	"Another Ticket name (e.g. Regular Ticket)"
//	@Param			ticket_category[1]			formData	string									false	"Another Ticket category (e.g. Regular)"
//	@Param			ticket_type[1]				formData	string									false	"Another Ticket type (e.g. Standing)"
//	@Param			ticket_price[1]				formData	number									false	"Another Ticket price (e.g. 500)"
//	@Param			ticket_number_of_ticker[1]	formData	integer									false	"Another Ticket quantity (e.g. 250)"
//	@Success		200							{object}	localization.StandardResponse{data=nil}	"Event updated successfully"
//	@Failure		400							{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		401							{object}	localization.StandardResponse{data=nil}	"Unauthorized"
//	@Failure		404							{object}	localization.StandardResponse{data=nil}	"Event not found"
//	@Failure		500							{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/events/{id} [patch]
func (a *eventAdapter) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateEvent", "handler", "event")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("event ID is required for update")
		localization.SendBadRequestResponse(w, localization.ErrorEventIDRequired.Message)
		return
	}

	req, err := eventcore.ParseEventRequestFromMultipartForm(r, false)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to parse event update request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if req.IsEmpty() {
		a.logger.Warnf("no data provided for event update, event ID: %s", id)
		localization.SendBadRequestResponse(w, localization.ErrorUpdateEventEmptyPayload.Message)
		return
	}

	span.SetAttributes(
		attribute.String("event.id", id),
		attribute.String("event.merchant_id", req.MerchantID),
	)

	if err := a.eventApp.UpdateEvent(ctx, id, req); err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to update event (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("event update request submitted successfully, event ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessEventUpdateRequestSubmitted, nil)
}

// DeleteEvent godoc
//
//	@Summary		Delete an event
//	@Description	Permanently delete an event by ID
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Event ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Event deleted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Event not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/events/{id} [delete]
func (a *eventAdapter) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteEvent", "handler", "event")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("event.id", id))

	if err := a.eventApp.DeleteEvent(ctx, id); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventDeleteRequestSubmitted, nil)
}

// EnableEvent godoc
//
//	@Summary		Enable an event
//	@Description	Enable an event by ID
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Event ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Event enable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Event not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/events/enable/{id} [patch]
func (a *eventAdapter) EnableEvent(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableEvent", "handler", "event")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("event.id", id))
	if err := a.eventApp.EnableDisableEvent(ctx, id, true); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventEnableRequestSubmitted, nil)
}

// DisableEvent godoc
//
//	@Summary		Disable an event
//	@Description	Disable an event by ID
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Event ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Event disable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Event not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/events/disable/{id} [patch]
func (a *eventAdapter) DisableEvent(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableEvent", "handler", "event")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("event.id", id))
	if err := a.eventApp.EnableDisableEvent(ctx, id, false); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventDisableRequestSubmitted, nil)
}

// GetEvent godoc
//
//	@Summary		Get event by ID
//	@Description	Retrieve an event's details by ID
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string											true	"Event ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.Event}	"Event retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}			"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}			"Event not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}			"Internal server error"
//	@Security		BearerAuth
//	@Router			/events/{id} [get]
func (a *eventAdapter) FetchEventByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchEventById", "handler", "event")
	defer span.End()
	var _ model.Event
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("event.id", id))
	event, err := a.eventApp.FetchEventByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventRetrieved, event)
}

// GetEvents godoc
//
//	@Summary		List events
//	@Description	Retrieve events with pagination and optional search
//	@Tags			Event
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																	false	"Page number"		default(1)
//	@Param			per_page	query		int																	false	"Items per page"	default(10)
//	@Param			search		query		string																false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=[]model.Event}	"Events retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"Internal server error"
//	@Security		BearerAuth
//	@Router			/events [get]
func (a *eventAdapter) FetchEvents(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchEvents", "handler", "event")
	defer span.End()
	filter := local_util.ExtractFilterParams(r)
	a.logger.Infof("fetching events with filter: %+v", filter)

	list, err := a.eventApp.FetchEvent(ctx, *filter)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("failed to fetch events: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("event.count", len(list.Data)))
	a.logger.Infof("events fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessEventsRetrieved, list)
}
