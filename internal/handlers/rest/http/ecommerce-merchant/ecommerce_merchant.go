package ecommercemerchant

import (
	"encoding/json"
	"errors"
	"net/http"

	miniappmerchant "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"
	miniappmerchat "cbe-super-app-cps-action/internal/constants/interfaces/ecommerce_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type ecommerceMerchantAdapter struct {
	srv    service.EcommerceMerchantService
	logger shared_utils.Logger
}

func NewEcommerceMerchantdapter(srv service.EcommerceMerchantService, logger shared_utils.Logger) miniappmerchat.EcommerceMerchant {
	return &ecommerceMerchantAdapter{srv: srv, logger: logger}
}

// Create Mini App Merchant
//
//	@Summary		Create Mini App Merchant
//	@Description	Creates a new mini app merchant
//	@Tags			EcommerceMerchant
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body			body		miniappmerchant.EcommerceMerchant	true	"Mini App Merchant DTO"
//	@Success		201				{object}	localization.StandardResponse{data=miniappmerchant.EcommerceMerchant}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant [post]
func (h *ecommerceMerchantAdapter) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "createMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	var reqDTO miniappmerchant.EcommerceMerchant

	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		span.RecordError(err)
		h.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	if err := reqDTO.Validate(true); err != nil {
		span.RecordError(err)
		h.logger.Errorf("Validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	// formattedPhone := local_util.FormatPhoneNumber(reqDTO.PhoneNumber)

	// reqDTO.PhoneNumber = formattedPhone
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		h.logger.Warnf("Incomplete user context: %+v", userContext)
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	createdMerchant, err := h.srv.Create(ctx, &reqDTO)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("Failed to create merchant: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Debugf("Created merchant: %+v", createdMerchant)

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantCreateRequestCreated, nil)
}

// Update Mini App Merchant
//
//	@Summary		Update Mini App Merchant
//	@Description	Updates an existing mini app merchant
//	@Tags			EcommerceMerchant
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string								true	"Merchant ID"
//	@Param			body				body		miniappmerchant.EcommerceMerchant	true	"Mini App Merchant DTO"
//	@Success		200					{object}	localization.StandardResponse
//	@Failure		400,401,404,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/{id} [patch]
func (h *ecommerceMerchantAdapter) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "updateMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var reqDTO miniappmerchant.EcommerceMerchant
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		span.RecordError(err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	if reqDTO.IsEmpty() {
		span.RecordError(errors.New("no data provided for update"))
		localization.SendErrorResponse(w, localization.ErrorNoDataProvidedForUpdate, nil, nil)
		return
	}

	if err := reqDTO.Validate(false); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("ecommerce_merchant.id", id))
	_, _, err := h.srv.Update(ctx, id, &reqDTO)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantUpdateRequestCreated, nil)
}

// Delete Mini App Merchant
//
//	@Summary		Delete Mini App Merchant
//	@Description	Deletes a mini app merchant by ID
//	@Tags			EcommerceMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/{id} [delete]
func (h *ecommerceMerchantAdapter) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "deleteMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("missing or invalid parameter 'id'"))
		h.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameters, nil, nil)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("ecommerce_merchant.id", id))
	err := h.srv.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("failed to delete merchant %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantDeleteRequestCreated, nil)
}

// Enable Mini App Merchant
//
//	@Summary		Enable Mini App Merchant
//	@Description	Enables a mini app merchant by ID
//	@Tags			EcommerceMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/enable/{id} [patch]
func (h *ecommerceMerchantAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "enableMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("missing or invalid parameter 'id'"))
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("ecommerce_merchant.id", id))
	if err := h.srv.EnableOrDisable(ctx, id, true); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantEnableRequestCreated, nil)
}

// Disable Mini App Merchant
//
//	@Summary		Disable Mini App Merchant
//	@Description	Disables a mini app merchant by ID
//	@Tags			EcommerceMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/disable/{id} [patch]
func (h *ecommerceMerchantAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "disableMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("missing or invalid parameter 'id'"))
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("ecommerce_merchant.id", id))
	if err := h.srv.EnableOrDisable(ctx, id, false); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantDisableRequestCreated, nil)
}

// {
//   "ok": false,
//   "status": 401,
//   "timestamp": "2025-10-01T16:32:52.1069555+03:00",
//   "message": "User is not authorized"
// }

// Get Mini App Merchant by ID
//
//	@Summary		Get Mini App Merchant by ID
//	@Description	Retrieves a mini app merchant by ID
//	@Tags			EcommerceMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse{data=miniappmerchant.MiniAppMerchantResponseDTO}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/{id} [get]
func (h *ecommerceMerchantAdapter) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "findMiniAppMerchantById", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	span.SetAttributes(attribute.String("ecommerce_merchant.id", id))
	result, err := h.srv.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, result)
}

// List Mini App Merchants with Pagination
//
//	@Summary		List Mini App Merchants
//	@Description	Retrieves a paginated list of mini app merchants
//	@Tags			EcommerceMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=miniappmerchant.PaginatedMiniAppResponseResponse}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant [get]
func (h *ecommerceMerchantAdapter) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "findAllMiniAppMerchants", "handler", "miniAppMerchant")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	phoneNumber, ok := filterParams.Filters["phone_number"].(string)
	if ok {
		formattedPhone := local_util.FormatPhoneNumber(phoneNumber)
		filterParams.Filters["phone_number"] = formattedPhone
	}

	miniAppMerchant, err := h.srv.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, miniAppMerchant)
}

// Merchant Lookup
// @Summary Merchant Lookup
// @Description Retrieves a merchant by ID
// @Tags EcommerceMerchant
// @Security BearerAuth
// @Produce json
// @Param merchant_id path string true "Merchant ID"
// @Success 200 {object} localization.StandardResponse{data=miniappmerchant.MerchantLookUpResponse}
// @Failure 400,401,404,500 {object} localization.StandardResponse{data=nil}
// @Router /ecommerce-merchant/merchant-lookup/{merchant_id} [get]
func (h *ecommerceMerchantAdapter) MerchantLookup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "miniAppMerchantLookup", "handler", "miniAppMerchant")
	defer span.End()
	merchantID := chi.URLParam(r, "merchant_id")
	if merchantID == "" {
		span.RecordError(errors.New("missing or invalid parameter 'merchant_id'"))
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}
	span.SetAttributes(attribute.String("ecommerce_merchant.id", merchantID))
	result, err := h.srv.MerchantLookup(ctx, merchantID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, result)
}
