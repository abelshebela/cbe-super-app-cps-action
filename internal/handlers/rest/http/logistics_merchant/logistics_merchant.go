package logistics_merchant_handler

import (
	local_util "cbe-super-app-cps-action/pkgs/utils"

	logistics_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/logistics_merchant"
	logistics_merchant_adaptor "cbe-super-app-cps-action/internal/constants/interfaces/logistics_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/logistics_merchant/core"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type LogisticsMerchantHandler struct {
	service service.LogisticsMerchantService
	logger  utils.Logger
}

// CreateLogisticMerchant godoc
//
//	@Summary		Create logistics merchant
//	@Description	Create a new logistics merchant with the provided information
//	@Tags			Logistics Merchant
//	@Accept			json
//	@Produce		json
//	@Param			body	body		logistics_merchant_dto.CreateLogisticsMerchantRequest	true	"Create logistics merchant request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Logistics merchant created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/logistics_merchants [post]
func (e *LogisticsMerchantHandler) CreateLogisticMerchant(w http.ResponseWriter, r *http.Request) {
	var req logistics_merchant_dto.CreateLogisticsMerchantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e.logger.Errorf("[CreateLogisticsMerchant] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := logistics_merchant_dto.Validation(req); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	m := core.CreateLogisticsMerchantRequestToModel(req)

	if err := e.service.Create(r.Context(), m); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessLogisticsMerchantCreated, nil)
}

// DeleteLogisticMerchant godoc
//
//	@Summary		Delete logistics merchant
//	@Description	Delete a logistics merchant by ID
//	@Tags			Logistics Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Logistics Merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Logistics merchant deleted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Logistics merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/logistics_merchants/{id} [delete]
func (e *LogisticsMerchantHandler) DeleteLogisticMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := e.service.Delete(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessLogisticsMerchantDeleted, nil)
}

// DisableLogisticMerchant godoc
//
//	@Summary		Disable logistics merchant
//	@Description	Disable a logistics merchant by ID
//	@Tags			Logistics Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Logistics Merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Logistics merchant disabled successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Logistics merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/logistics_merchants/disable/{id} [patch]
func (e *LogisticsMerchantHandler) DisableLogisticMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := e.service.EnableOrDisable(r.Context(), id, false); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessLogisticsMerchantDisabled, nil)
}

func (e *LogisticsMerchantHandler) EnableLogisticMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := e.service.EnableOrDisable(r.Context(), id, true); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessLogisticsMerchantEnabled, nil)
}

// GetLogisticMerchantByID godoc
//
//	@Summary		Get logistics merchant by ID
//	@Description	Retrieve a single logistics merchant by its identifier
//	@Tags			Logistics Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string													true	"Logistics Merchant ID"
//	@Success		200	{object}	localization.StandardResponse{data=logistics_merchant_dto.LogisticsMerchantResponse}	"Logistics merchant retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}													"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}													"Logistics merchant not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}													"Internal server error"
//	@Security		BearerAuth
//	@Router			/logistics_merchants/{id} [get]
func (e *LogisticsMerchantHandler) GetLogisticMerchantByID(w http.ResponseWriter, r *http.Request) {
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
	localization.SendSuccessResponse(w, localization.SuccessLogisticsMerchantFetched, result)
}

func (e *LogisticsMerchantHandler) GetLogisticMerchants(w http.ResponseWriter, r *http.Request) {
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
	localization.SendSuccessResponse(w, localization.SuccessLogisticsMerchantFetched, result)
}

// UpdateLogisticMerchant godoc
//
//	@Summary		Update logistics merchant
//	@Description	Update an existing logistics merchant by ID
//	@Tags			Logistics Merchant
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string													true	"Logistics Merchant ID"
//	@Param			body	body		logistics_merchant_dto.UpdateLogisticsMerchantRequest	true	"Update logistics merchant request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}				"Logistics merchant updated successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}				"Logistics merchant not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/logistics_merchants/{id} [patch]
func (e *LogisticsMerchantHandler) UpdateLogisticMerchant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	var req logistics_merchant_dto.UpdateLogisticsMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e.logger.Errorf("[UpdateLogisticsMerchant] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := logistics_merchant_dto.Validation(req); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	m := core.UpdateLogisticsMerchantRequestToModel(req)
	if err := e.service.Update(r.Context(), id, m); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessLogisticsMerchantUpdated, nil)
}

func NewLogisticsMerchantHandler(service service.LogisticsMerchantService, logger utils.Logger) logistics_merchant_adaptor.LogisticMerchantInboundAdaptor {
	return &LogisticsMerchantHandler{
		service: service,
		logger:  logger,
	}
}
