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

func (e *EventMerchantHandler) CreateEventMerchant(w http.ResponseWriter, r *http.Request) {
	var req event_merchant_dto.CreateEventMerchantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e.logger.Errorf("[CreateEventMerchant] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := event_merchant_dto.Validation(req); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	m := core.CreateEventMerchantRequestToModel(req)

	if err := e.service.Create(r.Context(), m); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessEventMerchantCreated, nil)
}

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

func (e *EventMerchantHandler) UpdateEventMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	var req event_merchant_dto.UpdateEventMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e.logger.Errorf("[UpdateEventMerchant] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := event_merchant_dto.Validation(req); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
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
