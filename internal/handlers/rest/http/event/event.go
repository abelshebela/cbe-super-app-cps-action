package eventhandler

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	eventcore "cbe-super-app-cps-action/internal/handlers/rest/http/event/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type eventAdapter struct {
	eventApp service.EventService
	logger   utils.Logger
}

func InitEventAdapter(eventApp service.EventService, logger utils.Logger) *eventAdapter {
	return &eventAdapter{
		eventApp: eventApp,
		logger:   logger,
	}
}

func (a *eventAdapter) CreateEvent(w http.ResponseWriter, r *http.Request) {
	req, err := eventcore.ParseEventRequestFromMultipartForm(r, true)
	if err != nil {
		a.logger.Errorf("failed to parse event request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(true); err != nil {
		a.logger.Errorf("event request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.eventApp.CreateEvent(r.Context(), req); err != nil {
		a.logger.Errorf("failed to create event: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("event creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessEventCreationRequestSubmitted, nil)
}

func (a *eventAdapter) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("event ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	req, err := eventcore.ParseEventRequestFromMultipartForm(r, false)
	if err != nil {
		a.logger.Errorf("failed to parse event update request: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	if req.IsEmpty() {
		a.logger.Warnf("no data provided for event update, event ID: %s", id)
		localization.SendErrorResponse(w, localization.ErrorNoDataProvidedForUpdate, nil, nil)
		return
	}

	if err := a.eventApp.UpdateEvent(r.Context(), id, req); err != nil {
		a.logger.Errorf("failed to update event (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("event update request submitted successfully, event ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessEventUpdateRequestSubmitted, nil)
}

func (a *eventAdapter) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	if err := a.eventApp.DeleteEvent(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventDeleteRequestSubmitted, nil)
}

func (a *eventAdapter) EnableEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	if err := a.eventApp.EnableDisableEvent(r.Context(), id, true); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventEnableRequestSubmitted, nil)
}

func (a *eventAdapter) DisableEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	if err := a.eventApp.EnableDisableEvent(r.Context(), id, false); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventDisableRequestSubmitted, nil)
}

func (a *eventAdapter) FetchEventByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorEventIDRequired, nil, nil)
		return
	}

	event, err := a.eventApp.FetchEventByID(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEventRetrieved, event)
}
func (a *eventAdapter) FetchEvents(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	a.logger.Infof("fetching events with filter: %+v", filter)

	list, err := a.eventApp.FetchEvent(r.Context(), *filter)
	if err != nil {
		a.logger.Errorf("failed to fetch events: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("events fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessEventsRetrieved, list)
}
