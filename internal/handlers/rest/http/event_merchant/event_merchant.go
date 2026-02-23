package event_merchant_handler

import (
	local_util "cbe-super-app-cps-action/pkgs/utils"

	event_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/event_merchant"
	event_merchant_port "cbe-super-app-cps-action/internal/constants/interfaces/event_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/event_merchant/core"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type EventMerchantHandler struct {
	service service.EventMerchantService
	logger  utils.Logger
}

// CreateEventMerchant godoc
//
//	@Summary		Create event merchant
//	@Description	Create a new event merchant with the provided information
//	@Tags			Event Merchant
//	@Accept			json
//	@Produce		json
//	@Param			body	body		event_merchant_dto.CreateEventMerchantRequest	true	"Create event merchant request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Event merchant created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/event_merchants [post]
func (e *EventMerchantHandler) CreateEventMerchant(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), e.logger)
	var req event_merchant_dto.CreateEventMerchantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[CreateEventMerchant] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := event_merchant_dto.Validation(req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	m := core.CreateEventMerchantRequestToModel(req)

	if err := e.service.Create(r.Context(), m); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantCreated, nil)
}

// DeleteEventMerchant godoc
//
//	@Summary		Delete event merchant
//	@Description	Delete an event merchant by ID
//	@Tags			Event Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Event merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Event merchant deleted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Event merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/event_merchant/{id} [delete]
func (e *EventMerchantHandler) DeleteEventMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := e.service.Delete(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantDeleted, nil)
}

// DisableEventMerchant godoc
//
//	@Summary		Disable event merchant
//	@Description	Disable an event merchant by ID
//	@Tags			Event Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Event merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Event merchant disabled successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Event merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/event_merchants/disable/{id} [patch]
func (e *EventMerchantHandler) DisableEventMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := e.service.EnableOrDisable(r.Context(), id, false); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantDisabled, nil)
}

// EnableEventMerchant godoc
//
//	@Summary		Enable event merchant
//	@Description	Enable an event merchant by ID
//	@Tags			Event Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Event merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Event merchant enabled successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Event merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/event_merchants/enable/{id} [patch]
func (e *EventMerchantHandler) EnableEventMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := e.service.EnableOrDisable(r.Context(), id, true); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantEnabled, nil)
}

// GetEventMerchantByID godoc
//
//	@Summary		Get event merchant by ID
//	@Description	Retrieve a single event merchant by its identifier
//	@Tags			Event Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string													true	"Event merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=event_merchant_dto.EventMerchantResponse}	"Event merchant retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}					"Event merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/event_merchants/{id} [get]
func (e *EventMerchantHandler) GetEventMerchantByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	result, err := e.service.FindByID(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantFetched, result)
}

// GetEventMerchants godoc
//
//	@Summary		Get all event merchants
//	@Description	Retrieve all event merchants with pagination and optional search
//	@Tags			Event Merchant
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int									false	"Page number"		default(1)
//	@Param			per_page	query		int									false	"Items per page"	default(10)
//	@Param			search		query		string								false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Event merchants retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/event_merchants [get]
func (e *EventMerchantHandler) GetEventMerchants(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

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

	result, err := e.service.FindAllWithPagination(r.Context(), *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantFetched, result)
}

// UpdateEventMerchant godoc
//
//	@Summary		Update event merchant
//	@Description	Update an existing event merchant by ID. Provide only fields to change.
//	@Tags			Event Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string										true	"Event merchant ID"
//	@Param			body	body		event_merchant_dto.UpdateEventMerchantRequest	true	"Update event merchant request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}		"Event merchant updated successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}		"Event merchant not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/event_merchants/{id} [patch]
func (e *EventMerchantHandler) UpdateEventMerchant(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), e.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	var req event_merchant_dto.UpdateEventMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[UpdateEventMerchant] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := event_merchant_dto.Validation(req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	m := core.UpdateEventMerchantRequestToModel(req)
	if err := e.service.Update(r.Context(), id, m); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantUpdated, nil)
}

func NewEventMerchantHandler(service service.EventMerchantService, logger utils.Logger) event_merchant_port.EventMerchantInboundAdaptor {
	return &EventMerchantHandler{
		service: service,
		logger:  logger,
	}
}
