package event

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	event_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/event"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	event_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/event"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
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

func (h *EventHTTPStore) MakerCreateEvent(w http.ResponseWriter, r *http.Request) {
	// Use reusable helper to parse the event request from multipart form
	req, err := ParseEventRequestFromMultipartForm(r)
	if err != nil {
		h.logger.Errorf("failed to parse event request: %v", err)
		utils.SendErrorResponse(w, utils.InvalidInput, http.StatusBadRequest, nil)
		return
	}

	// Validate
	if err := req.Validate(true); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}

	maker := entities.Maker{
		ID:          userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}

	// Use mappers to convert to domain DTO
	domainReq := ToDomainEventRequest(struct {
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

	 err = h.Application.CreateEvent(r.Context(), domainReq, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utils.WriteSuccessResponse(w, nil, "Event creation request submitted successfully")
}

func (h *EventHTTPStore) CheckerEvent(w http.ResponseWriter, r *http.Request) {
	// var req dto.EventCheckerRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	h.logger.Errorf("failed to decode JSON payload", "error", err)
	// 	utils.SendErrorResponse(w, utils.InvalidJSONPayload, 0, nil)
	// 	return
	// }
	// defer r.Body.Close()

	// if err := req.Validate(); err != nil {
	// 	h.logger.Errorf("invalid input in event checker request", "request_id", req.RequestID, "action", req.Action, "error", err)
	// 	utils.SendErrorResponse(w, utils.InvalidInput, http.StatusBadRequest, nil)
	// 	return
	// }

	// userContext := ctx_util.ExtractUserContext(r)
	// if userContext.IsIncomplete() {
	// 	h.logger.Errorf("incomplete user context", "request_id", req.RequestID, "action", req.Action)
	// 	utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
	// 	return
	// }

	// err := h.Application.CheckerCreateEvent(r.Context(), req.RequestID, req.Action, userContext.UserID, userContext.FullName, userContext.PhoneNumber)
	// if err != nil {
	// 	h.logger.Errorf("failed to create checker event", "request_id", req.RequestID, "action", req.Action, "user_id", userContext.UserID, "error", err)
	// 	utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
	// 	return
	// }
	// h.logger.Infof("event reviewed successfully", "request_id", req.RequestID, "action", req.Action, "user_id", userContext.UserID)
	// utils.BaseResponseMaker(map[string]any{"status": "success"}, w, "Event request reviewed successfully", http.StatusOK)
}

func (h *EventHTTPStore) FetchEventByID(w http.ResponseWriter, r *http.Request) {
	// var req dto.EventFetchRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	utils.SendErrorResponse(w, utils.InvalidJSONPayload, 0, nil)
	// 	return
	// }
	// if err := req.Validate(); err != nil {
	// 	utils.SendErrorResponse(w, utils.InvalidInput, http.StatusBadRequest, nil)
	// 	return
	// }
	// resp, err := h.Application.FetchEvent(r.Context(), req.RequestID)
	// if err != nil {
	// 	utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
	// 	return
	// }
	// data, err := utils.StructToMap(resp)
	// if err != nil {
	// 	utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
	// 	return
	// }
	// utils.BaseResponseMaker(data, w, "Event successfully retrieved", http.StatusOK)
}

func (h *EventHTTPStore) FetchEvents(w http.ResponseWriter, r *http.Request) {
	// q := r.URL.Query()
	// page, _ := strconv.Atoi(q.Get("page"))
	// if page < 1 {
	// 	page = DefaultPage
	// }
	// pageSize, _ := strconv.Atoi(q.Get("page_size"))
	// if pageSize < 1 {
	// 	pageSize = DefaultPageSize
	// }
	// if pageSize > MaxPageSize {
	// 	pageSize = MaxPageSize
	// }

	// offset := (page - 1) * pageSize
	// limit := pageSize
	// resp, err := h.Application.FetchEvent(r.Context(), limit, offset)
	// if err != nil {
	// 	utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
	// 	return
	// }
	// data, err := utils.StructToMap(resp)
	// if err != nil {
	// 	utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
	// 	return
	// }
	// utils.BaseResponseMaker(data, w, "Events successfully retrieved", http.StatusOK)
}
